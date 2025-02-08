package main

import (
	"Tracer/uploadBatch"
	"go.uber.org/zap/zapcore"
	"os"
)

func main() {
	logger, _ := uploadBatch.New(os.Stdout, zapcore.WarnLevel, uploadBatch.Configuration{
		Module:      "test",
		Service:     "test",
		OsVersion:   "123",
		Vendor:      "Linux",
		Host:        "localhost",
		DataCenter:  "datacenter",
		CloudMinion: "runner",
		VersionName: "1.0",
		Environment: "test",
		DeviceID:    "123",
	})
	logger.Warnf("first warning, %s", "test")
	logger.Errorf("first error, %s", "test")
	logger.Fatalf("first fatal, %s", "test")
}
