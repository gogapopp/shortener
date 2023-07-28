package save

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/lib/models"
	"github.com/gogapopp/shortener/internal/app/lib/urlshortener"
	"go.uber.org/zap"
)

type URLSaver interface {
	SaveURL(baseURL, longURL, shortURL, correlationID string) error
}

//go:generate mockgen -source=save.go -destination=mocks/mock.go
func PostSaveJSONHandler(log *zap.SugaredLogger, urlSaver URLSaver, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.api.save.PostSaveJSONHandler"
		// декодируем данные из тела запроса
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
		err = urlSaver.SaveURL(cfg.BaseAddr, req.URL, shortURL, "")
		if err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		// передаём значение в ответ
		var resp models.Response
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
