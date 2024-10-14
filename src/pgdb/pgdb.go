package pgdb

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var db *bun.DB

type Entry struct {
	bun.BaseModel `bun:"table:anomalies"`

	ID        int64 `bun:"id,pk,autoincrement"`
	SessionID string
	Frequency float64
	Timestamp time.Time
}

// conf PgConf
func DbConnect(dsn string) error {
	if dsn == "" {
		dsn = "postgres://postgres:123@localhost:5432/postgres?sslmode=disable"
	}
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db = bun.NewDB(sqldb, pgdialect.New())
	err := db.Ping()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	return nil
}

func AddNewTable(logger *slog.Logger, drop bool) error {
	err := db.Ping()
	if err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	logger.Info("Connected to database")

	if drop {
		err = db.ResetModel(context.Background(), (*Entry)(nil))
		if err != nil {

			return fmt.Errorf("failed to create table: %w", err)
		}
		logger.Info("Table droped")

		return nil
	}
	_, err = db.NewCreateTable().
		Model((*Entry)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	logger.Info("Table ready to use")
	return nil
}

func AddEntry(sessionId string, freq float64, logger *slog.Logger, time time.Time) error {
	if db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	entry := &Entry{
		SessionID: sessionId,
		Frequency: freq,
		Timestamp: time,
	}
	res, err := db.NewInsert().Model(entry).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("failed to insert entry: %w", err)
	}
	logger.Info("Entry inserted", slog.Any("result", res))
	return nil
}
