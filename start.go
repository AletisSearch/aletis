package aletis

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	sqlcdb "github.com/AletisSearch/aletis/db"
	"github.com/AletisSearch/aletis/internal/config"
	"github.com/AletisSearch/aletis/internal/db"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func init() {
	location, _ := time.LoadLocation("UTC")
	time.Local = location
}

func StartWebServer(options ...config.Option) error {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf, err := config.NewConfig(options...)
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}
	if err = conf.Validate(config.ValidDefault); err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	// Database
	if err = applyMigrations(conf); err != nil {
		return err
	}

	database, err := pgxpool.New(ctx, conf.DBconnStr())
	if err != nil {
		return err
	}
	err = database.Ping(ctx)
	if err != nil {
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
				err := queries.DeleteOld(ctxLimit, time.Now())
				if err != nil {
					slog.Error("err running DB Cleanup", "ERR", err)
				}
			case <-ctx.Done():
				slog.Info("Closing DB Cleanup")
				return
			}
		}
	})

	router, err := NewApp(ctx, &wg, conf, queries)
	if err != nil {
		return fmt.Errorf("failed to create app: %w", err)
	}

	serverErrChan := make(chan error, 1)

	addr := ":" + conf.Port
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		slog.Info("server listening", "addr", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrChan <- err
		}
	}()

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrChan:
		slog.Error("unexpected server error", "ERR", err)
	case sig := <-quitChan:
		slog.Info("received signal", "signal", sig)
	}

	ctxTmt, cancelTmt := context.WithTimeout(ctx, time.Second*15)
	defer cancelTmt()

	slog.Info("shutting server down...")
	if err = server.Shutdown(ctxTmt); err != nil {
		slog.Error("shutting down server", "ERR", err)
	}

	cancel()

	wg.Wait()

	return nil
}
func applyMigrations(conf *config.Config) error {
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
