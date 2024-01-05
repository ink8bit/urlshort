package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func New() (*zap.SugaredLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("cannot create zap logger: %w", err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	sugar := *logger.Sugar()
	return &sugar, nil
}
