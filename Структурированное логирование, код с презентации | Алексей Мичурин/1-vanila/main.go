package main

import (
	"log/slog"
)

func main() {
	slog.Info("OK (function)", "userID", 111)

	logger := slog.Default()
	logger = logger.With("data_center", "SomeDC")
	logger.Info("OK (method)", "userID", 222)
}
