package uploadBatch

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
)

const (
	ThrownMessageKey = "thrown_message"
	RequestIDKey     = "request_id"
	TagsKey          = "tags"
)

type Logger struct {
	*zap.Logger

	requestID *string
	tags      map[string]string
}

func New(out io.Writer, level zapcore.Level, configuration Configuration) (*Logger, error) {
	encoder, err := newEncoder(configuration)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracer encoder: %w", err)
	}
	syncer := zapcore.AddSync(out)

	levelEnabler := zap.LevelEnablerFunc(
		func(lvl zapcore.Level) bool {
			return lvl >= level
		},
	)
	core := zapcore.NewCore(encoder, syncer, levelEnabler)
	return &Logger{
		Logger:    zap.New(core),
		requestID: nil,
		tags:      map[string]string{},
	}, nil
}

func (logger *Logger) log(level zapcore.Level, message string, fields ...zap.Field) {
	tracerFields := []zap.Field{
		zap.String(ThrownMessageKey, message),
		zap.Strings(TagsKey, logger.encodeTags()),
	}
	if logger.requestID != nil {
		tracerFields = append(tracerFields, zap.String(RequestIDKey, *logger.requestID))
	}
	tracerFields = append(tracerFields, fields...)
	logger.Check(level, message).Write(tracerFields...)
}

func (logger *Logger) logf(level zapcore.Level, message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(message, args...)
	tracerFields := []zap.Field{
		zap.String(ThrownMessageKey, formattedMessage),
		zap.Strings(TagsKey, logger.encodeTags()),
	}
	if logger.requestID != nil {
		tracerFields = append(tracerFields, zap.String(RequestIDKey, *logger.requestID))
	}
	logger.Check(level, message).Write(tracerFields...)
}

func (logger *Logger) encodeTags() []string {
	var tags []string
	for key, val := range logger.tags {
		tags = append(tags, fmt.Sprintf("%s=%s", key, val))
	}
	return tags
}

func (logger *Logger) Info(message string, fields ...zap.Field) {
	logger.log(zapcore.InfoLevel, message, fields...)
}

func (logger *Logger) Infof(message string, args ...interface{}) {
	logger.logf(zapcore.InfoLevel, message, args...)
}

func (logger *Logger) Warn(message string, fields ...zap.Field) {
	logger.log(zapcore.WarnLevel, message, fields...)
}

func (logger *Logger) Warnf(message string, args ...interface{}) {
	logger.logf(zapcore.WarnLevel, message, args...)
}

func (logger *Logger) Error(message string, fields ...zap.Field) {
	logger.log(zapcore.ErrorLevel, message, fields...)
}

func (logger *Logger) Errorf(message string, args ...interface{}) {
	logger.logf(zapcore.ErrorLevel, message, args...)
}

func (logger *Logger) Fatal(message string, fields ...zap.Field) {
	logger.log(zapcore.FatalLevel, message, fields...)
}

func (logger *Logger) Fatalf(message string, args ...interface{}) {
	logger.logf(zapcore.FatalLevel, message, args...)
}

func (logger *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{
		requestID: &requestID,
		tags:      logger.tags,
	}
}

func (logger *Logger) WithTag(tag, value string) *Logger {
	tags := make(map[string]string)
	for key, val := range logger.tags {
		tags[key] = val
	}
	tags[tag] = value
	return &Logger{
		requestID: logger.requestID,
		tags:      tags,
	}
}
