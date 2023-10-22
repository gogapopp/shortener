// package batchsave contains an implementation of the PostBatchJSONhHandler handler
package batchsave

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/lib/models"
	"github.com/gogapopp/shortener/internal/app/lib/urlshortener"
	"github.com/gogapopp/shortener/internal/app/network-server/middlewares/auth"
	"go.uber.org/zap"
)

// BatchSave defines the batch method for saving URLs
type BatchSaver interface {
	BatchInsertURL(urls []models.BatchDatabaseResponse, userID string) error
}

// PostBatchJSONhHandler accepts an array of shortened link structures as input and in response an array of shortened link structures in json format
func PostBatchJSONhHandler(log *zap.SugaredLogger, batchSaver BatchSaver, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.api.save.PostBatchJSONhHandler"
		// getting the userID from the context that was set by the middleware UserIdentity
		userID, err := auth.GetUserIDFromCookie(r)
		if err != nil {
			userID = auth.GenerateUniqueUserID()
			auth.SetUserIDCookie(w, userID)
		}
		var req []models.BatchRequest
		var resp []models.BatchResponse
		var databaseResp []models.BatchDatabaseResponse
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		// we are starting to go through the request
		for k := range req {
			// check whether the passed value is a reference
			_, err = url.ParseRequestURI(req[k].OriginalURL)
			if err != nil {
				log.Infof("%s: %s", op, err)
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			// "compress" the link
			BatchShortURL := urlshortener.ShortenerURL(cfg.BaseAddr)
			// collecting data to send to the database
			databaseResp = append(databaseResp, models.BatchDatabaseResponse{
				ShortURL:      BatchShortURL,
				OriginalURL:   req[k].OriginalURL,
				CorrelationID: req[k].CorrelationID,
				UserID:        userID,
			})
			// collecting a json response
			resp = append(resp, models.BatchResponse{
				ShortURL:      BatchShortURL,
				CorrelationID: req[k].CorrelationID,
			})
		}
		err = batchSaver.BatchInsertURL(databaseResp, userID)
		if err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		// setting the Content-Type header and sending the response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
