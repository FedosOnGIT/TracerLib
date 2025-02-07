package tracer

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"hash/fnv"
	"time"
)

const (
	Vendor                = "Linux"
	Language              = "go"
	ProductionEnvironment = "production"
	MessageType           = "CRASH"
	CrashIDSource         = "message"
	StacktraceKey         = "stackTrace"
	LoggerName            = "TracerLogger"
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
	tags        []string
}

func NewEncoder(module, service string, tags map[string]string, flags Flags) (*Encoder, error) {
	jsonEncoder := zapcore.NewJSONEncoder(
		zapcore.EncoderConfig{
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			StacktraceKey:  StacktraceKey,
		},
	)

	versionCode, err := generateVersionCode(flags.VersionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get version code: %w", err)
	}

	environment := flags.Environment
	if environment == "" {
		environment = ProductionEnvironment
	}

	var tracerTags []string
	for tagKey, tagVal := range tags {
		tracerTags = append(tracerTags, fmt.Sprintf("%s=%s", tagKey, tagVal))
	}

	return &Encoder{
		Encoder:     jsonEncoder,
		osVersion:   flags.OsVersion,
		vendor:      Vendor,
		host:        flags.Host,
		dc:          flags.DC,
		versionName: flags.VersionName,
		versionCode: versionCode,
		deviceID:    flags.DeviceID,
		environment: environment,
		module:      module,
		service:     service,
		tags:        tracerTags,
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
		Language:      Language,
		Env:           encoder.environment,
		Service:       encoder.service,
		Hostname:      encoder.host,
		DC:            encoder.dc,
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
			Tags:          encoder.tags,
			CrashIDSource: CrashIDSource,
		},
		StackTrace: entry.Stack,
		Type:       MessageType,
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
