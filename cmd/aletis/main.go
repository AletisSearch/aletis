package main

import (
	"log/slog"
	"os"

	"github.com/AletisSearch/aletis"
	"github.com/AletisSearch/aletis/internal/config"
)

func main() {
	if err := aletis.StartWebServer(config.EnvConfigOptions()...); err != nil {
		slog.Error("application exited with an error", "ERR", err)
		os.Exit(1)
	}
}
