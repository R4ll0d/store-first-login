package logs

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func InitLog() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	logPath := viper.GetString("log.path")
	if logPath == "" {
		logPath = "./log"
	}

	if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
		panic(fmt.Sprintf("Failed to create log directory: %v", err))
	}

	dateStr := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("%s/%s_%s.log", logPath, hostname, dateStr)

	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(fmt.Sprintf("Failed to open log file: %v", err))
	}

	writeSyncer := zapcore.AddSync(logFile)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = ""

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		zap.InfoLevel,
	)

	log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

// Logging functions
func Info(message string, fields ...zap.Field) {
	log.Info(message, fields...)
	log.Sync()
}

func Debug(message string, fields ...zap.Field) {
	log.Debug(message, fields...)
	log.Sync()
}

func Error(message interface{}, fields ...zap.Field) {
	switch v := message.(type) {
	case error:
		log.Error(v.Error(), zap.String("error", v.Error()))
	case string:
		log.Error(v)
	}
	log.Sync()
}
