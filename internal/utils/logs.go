package utils

import "go.uber.org/zap"

var Logger = &zap.SugaredLogger{}

func NewLogger() {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	Logger = logger.Sugar()
}
