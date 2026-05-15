package url

import (
	"context"
	"smply/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	pgsql *pgxpool.Pool
}

func NewPostgresRepo(pgsql *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{pgsql}
}

func (r *PostgresRepo) Store(ctx context.Context, url string, alias string) (Url, error) {
	tx, err := r.pgsql.Begin(ctx)
	if err != nil {
		return Url{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	var id int64
	err = tx.QueryRow(ctx,
		`INSERT INTO urls (original, alias) VALUES ($1, $2) RETURNING id`,
		url, alias,
	).Scan(&id)
	if err != nil {
		return Url{}, err
	}

	if alias == "" {
		alias = utils.EncodeWithPadding(id, 2)
	}

	_, err = tx.Exec(ctx,
		`UPDATE urls SET alias = $1 WHERE id = $2`,
		alias, id,
	)
	if err != nil {
		return Url{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return Url{}, err
	}

	return Url{
		Id:       id,
		Original: url,
		Alias:    alias,
	}, nil
}

func (r *PostgresRepo) GetExact(ctx context.Context, originalUrl, alias string) (*Url, error) {
	var url Url

	err := r.pgsql.QueryRow(
		ctx,
		`SELECT id, original, alias FROM urls 
			WHERE original = $1 AND alias = $2`,
		originalUrl, alias).Scan(
		&url.Id,
		&url.Original,
		&url.Alias,
	)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (r *PostgresRepo) GetByAlias(ctx context.Context, alias string) (Url, error) {
	var url Url

	err := r.pgsql.QueryRow(
		ctx,
		`SELECT id, original, alias FROM urls WHERE alias = $1`,
		alias).Scan(
		&url.Id,
		&url.Original,
		&url.Alias,
	)
	if err != nil {
		return Url{}, err
	}

	return url, nil
}
