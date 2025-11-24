//go:build dev

package web

import (
	"io/fs"
	"log/slog"
	"os"
)

const IsDev = true

func init() {
	rootfs, err := os.OpenRoot("./web/dist")
	dynamicFS = rootfs.FS()
	if err != nil {
		slog.Error("failed to open web/dist root filesystem", "error", err)
		os.Exit(1)
	}
	if AssetsFs, err = fs.Sub(dynamicFS, "assets"); err != nil {
		slog.Error("failed to create assets sub-filesystem (dev)", "error", err)
		os.Exit(1)
	}
}
