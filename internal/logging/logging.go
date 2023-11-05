package logging

import (
	"log/slog"
	"os"
	"time"

	"gitlab.com/greyxor/slogor"
)

func InitLogger(level slog.Level, show bool) {
	handler := slogor.NewHandler(os.Stdout, &slogor.Options{
		TimeFormat: time.DateTime,
		Level:      level,
		ShowSource: show,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
