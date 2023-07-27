package save

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/lib/models"
	"github.com/gogapopp/shortener/internal/app/lib/urlshortener"
	"go.uber.org/zap"
)

type URLSaver interface {
	SaveURL(baseURL, longURL, shortURL string)
}

func PostSaveJSONHandler(log *zap.SugaredLogger, urlSaver URLSaver, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.api.save.PostSaveJSONHandler"
		// декодируем данные из тела запроса
		var req models.Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Info(fmt.Sprintf("%s: %s", op, err))
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		// проверяем является ли ссылка переданная в body валидной
		_, err := url.ParseRequestURI(req.URL)
		if err != nil {
			log.Info(fmt.Sprintf("%s: %s", op, err))
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		// делаем из обычной ссылки сжатую
		shortURL := urlshortener.ShortenerURL(cfg.BaseAddr, req.URL)
		// сохраняем короткую ссылку
		urlSaver.SaveURL(cfg.BaseAddr, req.URL, shortURL)
		// передаём значение в ответ
		var resp models.Response
		resp.ShortURL = shortURL
		// устанавливаем заголовок Content-Type и отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Info(fmt.Sprintf("%s: %s", op, err))
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
	}
}
