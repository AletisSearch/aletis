//go:build !dev

package web

import (
	"encoding/json"
	"io/fs"
	"log/slog"
	"os"
)

const IsDev = false

func init() {
	var err error
	if AssetsFs, err = fs.Sub(staticFs, "dist/assets"); err != nil {
		slog.Error("failed to create assets sub-filesystem (prod)", "error", err)
		os.Exit(1)
	}
	f, err := staticFs.ReadFile("dist/.vite/manifest.json")
	if err != nil {
		slog.Error("failed to read manifest.json", "error", err)
		os.Exit(1)
	}
	if err = json.Unmarshal(f, &pManifest); err != nil {
		slog.Error("failed to unmarshal manifest.json", "error", err)
		os.Exit(1)
	}
}
