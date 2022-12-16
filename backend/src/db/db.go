package db

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	migrate "github.com/rubenv/sql-migrate"
	_ "modernc.org/sqlite"
)

//go:embed migrations
var migrations embed.FS

type DB struct {
	d   *sqlx.DB
	log zerolog.Logger
}

func NewDB(fileName string, log zerolog.Logger) (*DB, error) {
	d, err := sqlx.Connect("sqlite", fileName)
	if err != nil {
		return nil, err
	}

	d.SetMaxOpenConns(1)

	fsys, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return nil, err
	}

	mig := &migrate.HttpFileSystemMigrationSource{
		FileSystem: http.FS(fsys),
	}

	if _, err = migrate.Exec(d.DB, "sqlite3", mig, migrate.Up); err != nil {
		return nil, err
	}

	log.Info().Msg("database prepared")

	return &DB{
		d: d,
	}, nil
}

type TX interface {
	squirrel.QueryerContext
	squirrel.ExecerContext
}

func (d *DB) NoTX() TX {
	return d.d
}

func (d *DB) NewTX(ctx context.Context) (*sql.Tx, error) {
	tx, err := d.d.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (d *DB) Close() error {
	d.log.Info().Msg("closing database")
	return d.d.Close()
}
