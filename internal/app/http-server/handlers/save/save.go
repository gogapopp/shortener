package save

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/lib/urlshortener"
	"go.uber.org/zap"
)

//go:generate mockgen -source=save.go -destination=mocks/mock.go
type URLSaver interface {
	SaveURL(baseURL, longURL, shortURL, correlationID string) error
}

func PostSaveHandler(log *zap.SugaredLogger, urlSaver URLSaver, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.PostSaveHandler"
		// читаем тело реквеста
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		bodyURL := string(body)
		// проверяем является ли ссылка переданная в body валидной
		_, err = url.ParseRequestURI(bodyURL)
		if err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		// делаем из обычной ссылки сжатую
		shortURL := urlshortener.ShortenerURL(cfg.BaseAddr, bodyURL)
		// сохраняем короткую ссылку
		err = urlSaver.SaveURL(cfg.BaseAddr, bodyURL, shortURL, "")
		if err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		// отправляем ответ
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, shortURL)
	}
}
