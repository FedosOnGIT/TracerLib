package tracer

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	ThrownMessageKey = "thrown_message"
	RequestIDKey     = "request_id"
)

const (
	MainTagKey = "project"
	MainTagVal = "godzen"
)

type Logger struct {
	Tracer *zap.Logger

	requestID *string
}

func New(module, service string, flags *Configuration, tags map[string]string) *Logger {
	if tags == nil {
		tags = make(map[string]string)
	}
	tags[MainTagKey] = MainTagVal
	tracerEncoder, err := NewEncoder(module, service, tags, flags.flags)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to initialize Tracer logger: %w", err))
		os.Exit(1)
	}
	tracerCore := newTracerLoggerCore(flags.tracerFile, tracerEncoder)
	tracerLogger := zap.New(tracerCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return &Logger{
		Tracer:    tracerLogger,
		requestID: nil,
	}
}

func newTracerLoggerCore(filePath string, encoder *Encoder) zapcore.Core {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	writeSyncer := zapcore.AddSync(file)

	levelEnabler := zap.LevelEnablerFunc(
		func(lvl zapcore.Level) bool {
			return lvl >= zapcore.WarnLevel
		},
	)
	return zapcore.NewCore(encoder, writeSyncer, levelEnabler)
}

func (logger *Logger) log(level zapcore.Level, message string, fields ...zap.Field) {
	tracerFields := []zap.Field{zap.String(ThrownMessageKey, message)}
	if logger.requestID != nil {
		tracerFields = append(tracerFields, zap.String(RequestIDKey, *logger.requestID))
	}
	logger.Tracer.Check(level, message).Write(tracerFields...)
}

func (logger *Logger) logf(level zapcore.Level, message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(message, args...)
	tracerFields := []zap.Field{zap.String(ThrownMessageKey, formattedMessage)}
	if logger.requestID != nil {
		tracerFields = append(tracerFields, zap.String(RequestIDKey, *logger.requestID))
	}
	logger.Tracer.Check(level, message).Write(tracerFields...)
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

func (logger *Logger) Sync() error {
	_ = logger.Tracer.Sync()
	return nil
}

func (logger *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{
		Tracer:    logger.Tracer,
		requestID: &requestID,
	}
}
