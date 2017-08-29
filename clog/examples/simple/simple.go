package main

import (
	"github.com/empirefox/cement/clog"
	"go.uber.org/zap"
)

// github.com/go-validator/validator v2
func main() {
	log, err := clog.NewLogger(clog.Config{
		Dev: true,
		DSN: "https://8c095ece96b5482bb4bdb49d109e1f73:7328e84e673f4c5daf59325a9fb1f7a5@sentry.io/156592",
	})
	if err != nil {
		panic(err)
	}

	richLog := log.Module("main-test")

	// Use logger in your service
	richLog.Info("Message describing logging reason", zap.String("key", "value"))
	richLog.Error("Even richer", zap.String("#subsystem", "example"))

	log.Sync()
}
