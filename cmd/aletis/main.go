package main

import (
	"log/slog"
	"os"
	"strings"
)

func main() {
	slog.Info(trimGetEnv("PORT"))
	slog.Info(trimGetEnv("PUBLIC"))
	slog.Info(trimGetEnv("AI_ENABLED"))
	slog.Info(trimGetEnv("OPENAI_URL"))
	slog.Info(trimGetEnv("OPENAI_API_KEY"))
	slog.Info(trimGetEnv("SEARXNG_HOST"))
}
func trimGetEnv(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}
