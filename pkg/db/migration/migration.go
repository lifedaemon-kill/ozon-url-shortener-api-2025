package migration

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Run(migrationsPath, dsn string) error {
	goose.SetDialect("postgres")

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}
	defer db.Close()

	if err = goose.Up(db, migrationsPath); err != nil {
		panic(err)
	}
	return nil
}
