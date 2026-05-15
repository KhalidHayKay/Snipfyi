package stat

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	pgsql *pgxpool.Pool
}

func NewPostresRepo(pgsql *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{pgsql}
}

func (r *PostgresRepo) Run(ctx context.Context, alias, referer, userAgent string, timestamp time.Time) error {
	tx, err := r.pgsql.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var id int
	err = tx.QueryRow(ctx, `
		UPDATE urls
		SET visited = visited + 1, last_visited = $1
		WHERE alias = $2
		RETURNING id
	`, timestamp, alias).Scan(&id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO click_events (url_id, referer, user_agent, timestamp)
		VALUES ($1, $2, $3, $4)
	`, id, referer, userAgent, timestamp)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepo) Get(ctx context.Context, alias string) (Stats, error) {
	var stats Stats

	err := r.pgsql.QueryRow(ctx,
		`SELECT original, alias, visited, created, last_visited FROM urls WHERE alias = $1`,
		alias).Scan(
		&stats.Original,
		&stats.Alias,
		&stats.Visited,
		&stats.Created,
		&stats.LastVisited,
	)
	if err != nil {
		return Stats{}, err
	}

	return stats, nil
}

func (r *PostgresRepo) GetAdmin(ctx context.Context) (AdminStats, error) {
	var s AdminStats

	// ── URLs ──
	if err := r.pgsql.QueryRow(ctx,
		`SELECT COUNT(*) FROM urls`,
	).Scan(&s.TotalUrls); err != nil {
		return s, err
	}

	if err := r.pgsql.QueryRow(ctx,
		`SELECT COUNT(*) FROM urls WHERE DATE(created) = CURRENT_DATE`,
	).Scan(&s.UrlsToday); err != nil {
		return s, err
	}

	// ── Total redirects (from urls.visited) ──
	if err := r.pgsql.QueryRow(ctx,
		`SELECT COALESCE(SUM(visited), 0) FROM urls`,
	).Scan(&s.TotalRedirects); err != nil {
		return s, err
	}

	// ── Click events: totals ──
	if err := r.pgsql.QueryRow(ctx,
		`SELECT COUNT(*) FROM click_events`,
	).Scan(&s.TotalClicks); err != nil {
		return s, err
	}

	if err := r.pgsql.QueryRow(ctx,
		`SELECT COUNT(*) FROM click_events WHERE DATE(timestamp) = CURRENT_DATE`,
	).Scan(&s.ClicksToday); err != nil {
		return s, err
	}

	if err := r.pgsql.QueryRow(ctx,
		`SELECT COUNT(*) FROM click_events WHERE timestamp >= NOW() - INTERVAL '7 days'`,
	).Scan(&s.ClicksThisWeek); err != nil {
		return s, err
	}

	// ── Unique links ever clicked ──
	if err := r.pgsql.QueryRow(ctx,
		`SELECT COUNT(DISTINCT url_id) FROM click_events`,
	).Scan(&s.UniqueLinksClicked); err != nil {
		return s, err
	}

	// ── Peak day (soft error — no data yet is fine) ──
	_ = r.pgsql.QueryRow(ctx, `
		SELECT TO_CHAR(DATE(timestamp), 'Mon DD, YYYY'), COUNT(*)
		FROM click_events
		GROUP BY DATE(timestamp)
		ORDER BY COUNT(*) DESC
		LIMIT 1
	`).Scan(&s.PeakDay, new(int64))

	// ── Peak hour 0–23 ──
	_ = r.pgsql.QueryRow(ctx, `
		SELECT EXTRACT(HOUR FROM timestamp)::int, COUNT(*) AS cnt
		FROM click_events
		GROUP BY EXTRACT(HOUR FROM timestamp)
		ORDER BY cnt DESC
		LIMIT 1
	`).Scan(&s.PeakHour, &s.PeakHourClicks)

	// ── Daily trend — last 7 days (always 7 rows via generate_series) ──
	trendRows, err := r.pgsql.Query(ctx, `
		SELECT TO_CHAR(gs::date, 'Mon DD') AS day, COALESCE(c.clicks, 0)
		FROM generate_series(
			(CURRENT_DATE - INTERVAL '6 days')::timestamp,
			CURRENT_DATE::timestamp,
			'1 day'
		) gs
		LEFT JOIN (
			SELECT DATE(timestamp) AS d, COUNT(*) AS clicks
			FROM click_events
			GROUP BY DATE(timestamp)
		) c ON c.d = gs::date
		ORDER BY gs
	`)
	if err != nil {
		return s, err
	}
	defer trendRows.Close()
	for trendRows.Next() {
		var d DailyStat
		if err := trendRows.Scan(&d.Date, &d.Clicks); err != nil {
			return s, err
		}
		s.DailyTrend = append(s.DailyTrend, d)
	}

	// ── Top 5 links by clicks ──
	topRows, err := r.pgsql.Query(ctx, `
		SELECT u.alias, COUNT(c.id) AS clicks
		FROM urls u
		JOIN click_events c ON c.url_id = u.id
		GROUP BY u.id, u.alias
		ORDER BY clicks DESC
		LIMIT 5
	`)
	if err != nil {
		return s, err
	}
	defer topRows.Close()
	for topRows.Next() {
		var l TopLink
		if err := topRows.Scan(&l.Alias, &l.Clicks); err != nil {
			return s, err
		}
		s.TopLinks = append(s.TopLinks, l)
	}

	// ── Top 5 referers ──
	refRows, err := r.pgsql.Query(ctx, `
		SELECT
			CASE WHEN referer = '' OR referer IS NULL THEN 'Direct' ELSE referer END,
			COUNT(*) AS clicks
		FROM click_events
		GROUP BY 1
		ORDER BY clicks DESC
		LIMIT 5
	`)
	if err != nil {
		return s, err
	}
	defer refRows.Close()
	for refRows.Next() {
		var r TopReferer
		if err := refRows.Scan(&r.Referer, &r.Clicks); err != nil {
			return s, err
		}
		s.TopReferers = append(s.TopReferers, r)
	}

	// ── Device breakdown (SQL-level UA bucketing) ──
	devRows, err := r.pgsql.Query(ctx, `
		SELECT
			CASE
				WHEN user_agent ILIKE '%mobile%'
				  OR user_agent ILIKE '%android%'
				  OR user_agent ILIKE '%iphone%' THEN 'Mobile'
				WHEN user_agent ILIKE '%tablet%'
				  OR user_agent ILIKE '%ipad%'   THEN 'Tablet'
				WHEN user_agent = '' OR user_agent IS NULL THEN 'Unknown'
				ELSE 'Desktop'
			END AS device,
			COUNT(*) AS clicks
		FROM click_events
		GROUP BY 1
		ORDER BY clicks DESC
	`)
	if err != nil {
		return s, err
	}
	defer devRows.Close()
	for devRows.Next() {
		var d TopDevice
		if err := devRows.Scan(&d.Device, &d.Clicks); err != nil {
			return s, err
		}
		s.TopDevices = append(s.TopDevices, d)
	}

	// ── API keys ──
	if err := r.pgsql.QueryRow(ctx,
		`SELECT COUNT(*) FROM api_keys`,
	).Scan(&s.TotalApiKeys); err != nil {
		return s, err
	}

	if err := r.pgsql.QueryRow(ctx, `
		SELECT COUNT(*) FROM api_keys
		WHERE expires_at > NOW() AND last_used_at IS NULL
	`).Scan(&s.ActiveApiKeys); err != nil {
		return s, err
	}

	// ── Magic tokens ──
	if err := r.pgsql.QueryRow(ctx,
		`SELECT COUNT(*) FROM magic_tokens`,
	).Scan(&s.TotalTokens); err != nil {
		return s, err
	}

	if err := r.pgsql.QueryRow(ctx,
		`SELECT COUNT(*) FROM magic_tokens WHERE DATE(created_at) = CURRENT_DATE`,
	).Scan(&s.TokensToday); err != nil {
		return s, err
	}

	return s, nil
}
