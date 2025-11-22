package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/AletisSearch/aletis"
)

func main() {
	location, _ := time.LoadLocation("UTC")
	time.Local = location

	if err := aletis.Start(aletis.EnvConfigOptions()...); err != nil {
		slog.Error("application exited with an error", "ERR", err)
		os.Exit(1)
	}
}
