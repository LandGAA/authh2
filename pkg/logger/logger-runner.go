package logger

import (
	"go.uber.org/zap"
	"time"
)

func LoggerRun() {
	if err := InitLogger(); err != nil {
		panic(err)
	}
	Logger.Info("Логгер запущен", zap.String("date", time.Now().String()))
	defer Logger.Sync()
}
