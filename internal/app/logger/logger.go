package logger

import (
	"go.uber.org/zap"
)

// В принципе и логгер можно при желание сделать мидлварью
// NewLogger создаём логгер
func NewLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()
	Sugar := logger.Sugar()

	return Sugar, nil
}
