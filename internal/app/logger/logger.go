package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

// NewLogger создаём логгер
func NewLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()
	Sugar := logger.Sugar()
	Log = Sugar

	return Sugar, nil
}

func RequestLogger(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Log.Infow("request logger",
			"method", r.Method,
			"url", fmt.Sprintf("%s%s", r.Host, r.URL.Path),
		)
		h(w, r)
	})
}

func ResponseLogger(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// читаем боди запоса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		// возвращаем данные обратно
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		var data interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Fatal(err)
		}

		Log.Infow("response logger",
			"method", r.Method,
			"host", r.Host,
			"body", data,
		)
		h(w, r)
	})
}
