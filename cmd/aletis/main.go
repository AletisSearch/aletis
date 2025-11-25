package main

import (
	"log/slog"
	"os"

	"github.com/AletisSearch/aletis"
)

func main() {
	if err := aletis.StartWebServer(aletis.EnvConfigOptions()...); err != nil {
		slog.Error("application exited with an error", "ERR", err)
		os.Exit(1)
	}
}
