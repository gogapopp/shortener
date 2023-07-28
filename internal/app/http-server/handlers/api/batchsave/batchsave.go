package batchsave

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/lib/models"
	"github.com/gogapopp/shortener/internal/app/lib/urlshortener"
	"go.uber.org/zap"
)

type BatchSaver interface {
	BatchInsertURL(urls []models.BatchDatabaseResponse) error
}

func PostBatchJSONhHandler(log *zap.SugaredLogger, batchSaver BatchSaver, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.api.save.PostBatchJSONhHandler"
		var req []models.BatchRequest
		var resp []models.BatchResponse
		var databaseResp []models.BatchDatabaseResponse
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		// начинаем проходить по реквесту
		for k := range req {
			log.Info(req[k])
			// проверяем является ли переданное значение ссылкой
			_, err := url.ParseRequestURI(req[k].OriginalURL)
			if err != nil {
				log.Infof("%s: %s", op, err)
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			// "сжимаем" ссылку
			BatchShortURL := urlshortener.ShortenerURL(cfg.BaseAddr, req[k].OriginalURL)

			// собираем данные для отправки в бд
			databaseResp = append(databaseResp, models.BatchDatabaseResponse{
				ShortURL:      BatchShortURL,
				OriginalURL:   req[k].OriginalURL,
				CorrelationID: req[k].CorrelationID,
			})
			// собираем json ответ
			resp = append(resp, models.BatchResponse{
				ShortURL:      BatchShortURL,
				CorrelationID: req[k].CorrelationID,
			})
		}
		err := batchSaver.BatchInsertURL(databaseResp)
		if err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		// устанавливаем заголовок Content-Type и отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
