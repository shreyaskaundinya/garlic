package utils

import "go.uber.org/zap"

func NewLogger() *zap.SugaredLogger {
	dev, _ := zap.NewDevelopment()
	return dev.Sugar()
}
