package utils

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

func InitLogger() {
	level := slog.LevelInfo
	if viper.GetBool("verbose") {
		level = slog.LevelDebug
	}
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(h))
}
