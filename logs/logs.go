package logs

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func InitLog() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = ""
	encoderConfig.CallerKey = "caller" // เพิ่ม key สำหรับ caller
	encoderConfig.MessageKey = "msg"   // เปลี่ยนให้ใช้ "msg" แทนข้อความ

	// ใช้ stdout แทนการเขียนไฟล์
	consoleSyncer := zapcore.Lock(os.Stdout)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // ใช้ JSON Encoder
		consoleSyncer,
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
