// package save содержит в себе код хендлера PostSaveJSONHandler
package save

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/http-server/middlewares/auth"
	"github.com/gogapopp/shortener/internal/app/lib/models"
	"github.com/gogapopp/shortener/internal/app/lib/urlshortener"
	"github.com/gogapopp/shortener/internal/app/storage/postgres"
	"go.uber.org/zap"
)

// URLSaver определяет метод SaveURL и GetShortURL
type URLSaver interface {
	SaveURL(longURL, shortURL, correlationID string, userID string) error
	GetShortURL(longURL string) string
}

// PostSaveJSONHandler принимает в JSON формате url и возвращает сокращенный URL
func PostSaveJSONHandler(log *zap.SugaredLogger, urlSaver URLSaver, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.api.save.PostSaveJSONHandler"
		// получаем userID из контекста который был установлен мидлвеером userIdentity
		userID, err := auth.GetUserIDFromCookie(r)
		if err != nil {
			userID = auth.GenerateUniqueUserID()
			auth.SetUserIDCookie(w, userID)
		}
		// декодируем данные из тела запроса
		var resp models.Response
		var req models.Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		// проверяем является ли ссылка переданная в body валидной
		_, err = url.ParseRequestURI(req.URL)
		if err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		// делаем из обычной ссылки сжатую
		shortURL := urlshortener.ShortenerURL(cfg.BaseAddr)
		// сохраняем короткую ссылку
		err = urlSaver.SaveURL(req.URL, shortURL, "", userID)
		if err != nil {
			log.Infof("%s: %s", op, err)
			if errors.Is(postgres.ErrURLExists, err) {
				shortURL := urlSaver.GetShortURL(req.URL)
				// передаём значение в ответ
				resp.ShortURL = shortURL
				// устанавливаем заголовок Content-Type и отправляем ответ
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict)
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
