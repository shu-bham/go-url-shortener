package logger

import (
	"github.com/shu-bham/go-url-shortener/internal/config"
	"github.com/sirupsen/logrus"
	"os"
)

func NewLogger(cfg config.Config) *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)

	switch cfg.Logger.Format {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "caller",
			},
		})
	default:
		log.SetFormatter(&logrus.TextFormatter{})
	}

	level, err := logrus.ParseLevel(cfg.Logger.Level)
	if err != nil {
		log.WithError(err).Warnf("Invalid logger level '%s'. Defaulting to 'info'", cfg.Logger.Level)
		level = logrus.InfoLevel
	}
	log.SetLevel(level)
	log.SetReportCaller(false)
	return log
}
