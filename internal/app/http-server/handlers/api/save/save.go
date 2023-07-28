package save

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/lib/models"
	"github.com/gogapopp/shortener/internal/app/lib/urlshortener"
	"github.com/gogapopp/shortener/internal/app/storage/postgres"
	"go.uber.org/zap"
)

//go:generate mockgen -source=save.go -destination=mocks/mock.go
type URLSaver interface {
	SaveURL(longURL, shortURL, correlationID string) error
	GetShortURL(longURL string) string
}

func PostSaveJSONHandler(log *zap.SugaredLogger, urlSaver URLSaver, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.api.save.PostSaveJSONHandler"
		// декодируем данные из тела запроса
		var resp models.Response
		var req models.Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		// проверяем является ли ссылка переданная в body валидной
		_, err := url.ParseRequestURI(req.URL)
		if err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		// делаем из обычной ссылки сжатую
		shortURL := urlshortener.ShortenerURL(cfg.BaseAddr, req.URL)
		// сохраняем короткую ссылку
		err = urlSaver.SaveURL(req.URL, shortURL, "")
		if err != nil {
			log.Infof("%s: %s", op, err)
			if errors.Is(postgres.ErrURLExists, err) {
				shortURL := urlSaver.GetShortURL(req.URL)
				// передаём значение в ответ
				resp.ShortURL = shortURL
				// устанавливаем заголовок Content-Type и отправляем ответ
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					log.Infof("%s: %s", op, err)
					http.Error(w, "something went wrong", http.StatusInternalServerError)
					return
				}
				return
			}
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		// передаём значение в ответ
		resp.ShortURL = shortURL
		// устанавливаем заголовок Content-Type и отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
	}
}
