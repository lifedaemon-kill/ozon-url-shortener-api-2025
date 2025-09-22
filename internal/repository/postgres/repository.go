package postgres

import (
	"context"
	"database/sql"
	"errors"
	"main/internal/config"
	"main/internal/errs"
	"main/internal/repository/entities"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sqlx.DB
}

func (p *Postgres) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func New(conf config.DBConfig) (*Postgres, error) {
	db, err := sqlx.Connect("postgres", conf.DSN())
	if err != nil {
		return nil, err
	}
	return &Postgres{
		db: db,
	}, nil
}

func (p *Postgres) AddUrlRelation(ctx context.Context, req entities.AddUrlRelationRequest) (entities.AddUrlRelationResponse, error) {
	query := `INSERT INTO relations (origin_url, shortened_url)
              VALUES ($1, $2)`

	_, err := p.db.ExecContext(ctx, query, req.OriginUrl, req.ShortenedUrl)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" { // unique_violation
			switch pqErr.Constraint {
			case "relations_origin_url_key":
				return entities.AddUrlRelationResponse{""}, errs.ErrorAlreadyExist
			case "relations_shortened_url_key":
				return entities.AddUrlRelationResponse{""}, errs.ErrorRepositoryDuplicate
			}
		}
		return entities.AddUrlRelationResponse{""}, err
	}

	return entities.AddUrlRelationResponse{req.ShortenedUrl}, nil
}

func (p *Postgres) GetOriginURLFromShortened(ctx context.Context, req entities.GetOriginURLFromShortenedUrlRequest) (entities.GetOriginURLFromShortenedUrlResponse, error) {
	query := `SELECT origin_url FROM relations WHERE shortened_url=$1`

	row := p.db.QueryRowContext(ctx, query, req.ShortenedUrl)

	var originUrl string
	err := row.Scan(&originUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.GetOriginURLFromShortenedUrlResponse{""}, errs.ErrorRepositoryUrlEmpty
		}
		return entities.GetOriginURLFromShortenedUrlResponse{""}, err
	}
	return entities.GetOriginURLFromShortenedUrlResponse{originUrl}, nil
}
