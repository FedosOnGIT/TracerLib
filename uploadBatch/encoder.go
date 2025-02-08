package uploadBatch

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"hash/fnv"
	"time"
)

const (
	LoggerName = "UploadBatchTracerLogger"
)

type Encoder struct {
	zapcore.Encoder
	osVersion   string
	vendor      string
	host        string
	dc          string
	versionName string
	versionCode int
	deviceID    string
	environment string
	module      string
	service     string
}

func newEncoder(configuration Configuration) (*Encoder, error) {
	jsonEncoder := zapcore.NewJSONEncoder(
		zapcore.EncoderConfig{
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			StacktraceKey:  "stackTrace",
		},
	)

	versionCode, err := generateVersionCode(configuration.VersionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get version code: %w", err)
	}

	environment := configuration.Environment

	return &Encoder{
		Encoder:     jsonEncoder,
		osVersion:   configuration.OsVersion,
		vendor:      configuration.Vendor,
		host:        configuration.Host,
		dc:          configuration.DataCenter,
		versionName: configuration.VersionName,
		versionCode: versionCode,
		deviceID:    configuration.DeviceID,
		environment: environment,
		module:      configuration.Module,
		service:     configuration.Service,
	}, nil
}

func generateVersionCode(versionName string) (int, error) {
	versionHash := fnv.New32a()
	_, err := versionHash.Write([]byte(versionName))
	if err != nil {
		return 0, err
	}
	return int(versionHash.Sum32()), nil
}

func (encoder *Encoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	fieldMap := extractFields(fields)

	properties := LogProperties{
		Level:         entry.Level.String(),
		Logger:        LoggerName,
		Language:      "go",
		Environment:   encoder.environment,
		Service:       encoder.service,
		Hostname:      encoder.host,
		DataCenter:    encoder.dc,
		Message:       entry.Message,
		ThrownMessage: getStringField(fieldMap, ThrownMessageKey),
		RequestID:     getStringField(fieldMap, RequestIDKey),
	}

	logMessage := LogMessage{
		UploadBean: UploadBean{
			OSVersion:     encoder.osVersion,
			Vendor:        encoder.vendor,
			VersionName:   encoder.versionName,
			VersionCode:   encoder.versionCode,
			DeviceID:      encoder.deviceID,
			Module:        encoder.module,
			Properties:    properties,
			Tags:          getStringArrayField(fieldMap, TagsKey),
			CrashIDSource: "message",
		},
		StackTrace: entry.Stack,
		Type:       "CRASH",
		Timestamp:  entry.Time.Format(time.RFC3339),
	}

	return marshalLogMessage(logMessage)
}

func extractFields(fields []zapcore.Field) map[string]interface{} {
	objectMap := zapcore.NewMapObjectEncoder()
	for _, field := range fields {
		field.AddTo(objectMap)
	}
	return objectMap.Fields
}

func getStringField(fieldMap map[string]interface{}, key string) string {
	if value, exists := fieldMap[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

func getStringArrayField(fieldMap map[string]interface{}, key string) []string {
	if value, exists := fieldMap[key]; exists {
		if strArray, ok := value.([]string); ok {
			return strArray
		}
	}
	return nil
}

func marshalLogMessage(logMessage LogMessage) (*buffer.Buffer, error) {
	logData, err := json.Marshal(logMessage)
	if err != nil {
		return nil, err
	}

	pool := buffer.NewPool()
	logBuffer := pool.Get()
	logBuffer.AppendBytes(logData)
	logBuffer.AppendString("\n")
	return logBuffer, nil
}

func (encoder *Encoder) Clone() zapcore.Encoder {
	return &Encoder{
		Encoder:     encoder.Encoder.Clone(),
		osVersion:   encoder.osVersion,
		vendor:      encoder.vendor,
		host:        encoder.host,
		dc:          encoder.dc,
		versionName: encoder.versionName,
		versionCode: encoder.versionCode,
		deviceID:    encoder.deviceID,
		environment: encoder.environment,
		module:      encoder.module,
		service:     encoder.service,
	}
}
