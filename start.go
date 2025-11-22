package aletis

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"sync"
	"time"

	sqlcdb "github.com/AletisSearch/aletis/db"
	"github.com/AletisSearch/aletis/internal/db"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Start(options ...Option) error {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf, err := NewConfig(options...)
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	// Database
	database, err := pgxpool.New(ctx, conf.DBconnStr())
	if err != nil {
		return err
	}
	err = database.Ping(ctx)
	if err != nil {
		return err
	}
	if err = applyMigrations(conf); err != nil {
		return err
	}

	queries := db.New(database)

	wg.Go(func() {
		t := time.Tick(5 * time.Minute)
		for {
			select {
			case <-t:
				ctxLimit, cancel := context.WithTimeout(ctx, time.Second*30)
				defer cancel()
				err := queries.DeleteOld(ctxLimit, pgtype.Timestamptz{Time: time.Now()})
				if err != nil {
					slog.Error("err running DB Cleanup", "Error", err)
				}
			case <-ctx.Done():
				slog.Info("Closing DB Cleanup")
				return
			}
		}
	})
	wg.Done()
	return nil
}
func applyMigrations(conf *Config) error {
	u, _ := url.Parse(conf.DBconnStr())
	db := dbmate.New(u)
	db.FS = sqlcdb.FS
	db.MigrationsDir = []string{"./migrations"}

	migrations, err := db.FindMigrations()
	if err != nil {
		return err
	}
	for _, m := range migrations {
		slog.Info("Migration:", "Version", m.Version, "Path", m.FilePath, "Applied", m.Applied)
	}

	slog.Info("Applying migrations...")
	err = db.CreateAndMigrate()
	if err != nil {
		return err
	}
	return nil
}
