package service

import (
	"context"
	"smply/internal/storage"
)

type TopLink struct {
	Alias  string
	Clicks int64
}

type DailyStat struct {
	Date   string
	Clicks int64
}

type TopReferer struct {
	Referer string
	Clicks  int64
}

type TopDevice struct {
	Device string
	Clicks int64
}

type AdminStats struct {
	// URLs
	TotalUrls int64
	UrlsToday int64

	// Redirects (sum of urls.visited)
	TotalRedirects int64

	// Click events
	TotalClicks        int64
	ClicksToday        int64
	ClicksThisWeek     int64
	UniqueLinksClicked int64

	// Peak activity
	PeakDay        string
	PeakHour       int
	PeakHourClicks int64

	// Click trend — last 7 days
	DailyTrend []DailyStat

	// Top links
	TopLinks []TopLink

	// Top referers
	TopReferers []TopReferer

	// Device breakdown
	TopDevices []TopDevice

	// API
	TotalApiKeys  int64
	ActiveApiKeys int64

	// Magic tokens
	TotalTokens int64
	TokensToday int64
}

func GetAdminStats(ctx context.Context) (AdminStats, error) {
	var s AdminStats

	// ── URLs ──
	if err := storage.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM urls`,
	).Scan(&s.TotalUrls); err != nil {
		return s, err
	}

	if err := storage.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM urls WHERE DATE(created) = CURRENT_DATE`,
	).Scan(&s.UrlsToday); err != nil {
		return s, err
	}

	// ── Total redirects (from urls.visited) ──
	if err := storage.DB.QueryRow(ctx,
		`SELECT COALESCE(SUM(visited), 0) FROM urls`,
	).Scan(&s.TotalRedirects); err != nil {
		return s, err
	}

	// ── Click events: totals ──
	if err := storage.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM click_events`,
	).Scan(&s.TotalClicks); err != nil {
		return s, err
	}

	if err := storage.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM click_events WHERE DATE(timestamp) = CURRENT_DATE`,
	).Scan(&s.ClicksToday); err != nil {
		return s, err
	}

	if err := storage.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM click_events WHERE timestamp >= NOW() - INTERVAL '7 days'`,
	).Scan(&s.ClicksThisWeek); err != nil {
		return s, err
	}

	// ── Unique links ever clicked ──
	if err := storage.DB.QueryRow(ctx,
		`SELECT COUNT(DISTINCT url_id) FROM click_events`,
	).Scan(&s.UniqueLinksClicked); err != nil {
		return s, err
	}

	// ── Peak day (soft error — no data yet is fine) ──
	_ = storage.DB.QueryRow(ctx, `
		SELECT TO_CHAR(DATE(timestamp), 'Mon DD, YYYY'), COUNT(*)
		FROM click_events
		GROUP BY DATE(timestamp)
		ORDER BY COUNT(*) DESC
		LIMIT 1
	`).Scan(&s.PeakDay, new(int64))

	// ── Peak hour 0–23 ──
	_ = storage.DB.QueryRow(ctx, `
		SELECT EXTRACT(HOUR FROM timestamp)::int, COUNT(*) AS cnt
		FROM click_events
		GROUP BY EXTRACT(HOUR FROM timestamp)
		ORDER BY cnt DESC
		LIMIT 1
	`).Scan(&s.PeakHour, &s.PeakHourClicks)

	// ── Daily trend — last 7 days (always 7 rows via generate_series) ──
	trendRows, err := storage.DB.Query(ctx, `
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
	topRows, err := storage.DB.Query(ctx, `
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
	refRows, err := storage.DB.Query(ctx, `
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
	devRows, err := storage.DB.Query(ctx, `
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
	if err := storage.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM api_keys`,
	).Scan(&s.TotalApiKeys); err != nil {
		return s, err
	}

	if err := storage.DB.QueryRow(ctx, `
		SELECT COUNT(*) FROM api_keys
		WHERE expires_at > NOW() AND last_used_at IS NULL
	`).Scan(&s.ActiveApiKeys); err != nil {
		return s, err
	}

	// ── Magic tokens ──
	if err := storage.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM magic_tokens`,
	).Scan(&s.TotalTokens); err != nil {
		return s, err
	}

	if err := storage.DB.QueryRow(ctx,
		`SELECT COUNT(*) FROM magic_tokens WHERE DATE(created_at) = CURRENT_DATE`,
	).Scan(&s.TokensToday); err != nil {
		return s, err
	}

	return s, nil
}
