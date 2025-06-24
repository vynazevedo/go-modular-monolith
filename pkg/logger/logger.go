package logger

import (
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Log *StructuredLogger

type Config struct {
	Level       string `json:"level" mapstructure:"level"`
	Format      string `json:"format" mapstructure:"format"`
	ServiceName string `json:"service_name" mapstructure:"service_name"`
}

type StructuredLogger struct {
	*logrus.Logger
	ServiceName string
}

type CustomJSONFormatter struct {
	ServiceName string
}

func (f *CustomJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields)
	for k, v := range entry.Data {
		data[k] = v
	}

	data["global_event_timestamp"] = entry.Time.UTC().Format(time.RFC3339)
	data["level"] = entry.Level.String()
	data["message"] = entry.Message
	data["service_name"] = f.ServiceName

	formatter := &logrus.JSONFormatter{}
	return formatter.Format(&logrus.Entry{
		Logger:  entry.Logger,
		Data:    data,
		Time:    entry.Time,
		Level:   entry.Level,
		Message: "",
	})
}

func getServiceName(serviceName string) string {
	if serviceName == "" {
		return "go-modular-monolith"
	}
	return serviceName
}

func Init(config Config) {
	logger := logrus.New()

	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	logger.SetOutput(os.Stdout)

	switch config.Format {
	case "json":
		logger.SetFormatter(&CustomJSONFormatter{
			ServiceName: getServiceName(config.ServiceName),
		})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}

	Log = &StructuredLogger{
		Logger:      logger,
		ServiceName: getServiceName(config.ServiceName),
	}
}

func (sl *StructuredLogger) WithContext(ctx context.Context) *logrus.Entry {
	entry := sl.Logger.WithContext(ctx)

	if traceID := ctx.Value("trace_id"); traceID != nil {
		entry = entry.WithField("trace_id", traceID)
	}
	if sessionID := ctx.Value("session_id"); sessionID != nil {
		entry = entry.WithField("session_id", sessionID)
	}

	return entry
}

func (sl *StructuredLogger) WithEventName(eventName string) *logrus.Entry {
	return sl.Logger.WithField("global_event_name", eventName)
}

func (sl *StructuredLogger) WithContextFields(fields logrus.Fields) *logrus.Entry {
	if len(fields) > 0 {
		return sl.Logger.WithField("context", fields)
	}
	return sl.Logger.WithFields(logrus.Fields{})
}

func GetLogger() *StructuredLogger {
	if Log == nil {
		Init(Config{
			Level:       "info",
			Format:      "text",
			ServiceName: "go-modular-monolith",
		})
	}
	return Log
}

func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return GetLogger().WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return GetLogger().WithFields(fields)
}

func WithContext(ctx context.Context) *logrus.Entry {
	return GetLogger().WithContext(ctx)
}

func WithEventName(eventName string) *logrus.Entry {
	return GetLogger().WithEventName(eventName)
}

func WithContextFields(fields logrus.Fields) *logrus.Entry {
	return GetLogger().WithContextFields(fields)
}
