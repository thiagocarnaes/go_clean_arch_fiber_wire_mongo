package logger

import (
	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999Z", // Formato ISO 8601 para Datadog
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})
	logger.SetLevel(logrus.InfoLevel)
	// Adicionar campos padrão para Datadog
	//logger.AddHook(&DatadogHook{})
	return logger
}

// DatadogHook adiciona campos específicos do Datadog
//type DatadogHook struct{}
//
//func (h *DatadogHook) Levels() []logrus.Level {
//	return logrus.AllLevels
//}
//
//func (h *DatadogHook) Fire(entry *logrus.Entry) error {
//	entry.Data["ddsource"] = "go"
//	entry.Data["service"] = "user-management"
//	entry.Data["ddtags"] = "env:dev,app:fiber"
//	return nil
//}
