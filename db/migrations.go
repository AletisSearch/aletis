package sqlcdb

import (
	"embed"
	"log/slog"
	"net/url"

	"github.com/AletisSearch/aletis/internal/config"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
)

//go:embed migrations/*.sql
var FS embed.FS

func ApplyMigrations(conf *config.Config) error {
	u, _ := url.Parse(conf.DBconnStr())
	db := dbmate.New(u)
	db.FS = FS
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
