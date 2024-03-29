// package save contains the PostSaveJSONHandler handler code
package save

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/lib/models"
	"github.com/gogapopp/shortener/internal/app/lib/urlshortener"
	"github.com/gogapopp/shortener/internal/app/network-server/middlewares/auth"
	"github.com/gogapopp/shortener/internal/app/storage/postgres"
	"go.uber.org/zap"
)

// URLSaver defines the SaveURL and GetShortURL methods
type URLSaver interface {
	SaveURL(longURL, shortURL, correlationID string, userID string) error
	GetShortURL(longURL string) string
}

// PostSaveJSONHandler accepts a url in JSON format and returns an abbreviated URL
func PostSaveJSONHandler(log *zap.SugaredLogger, urlSaver URLSaver, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.api.save.PostSaveJSONHandler"
		// getting the user ID from the context that was set by the middleware userIdentity
		userID, err := auth.GetUserIDFromCookie(r)
		if err != nil {
			userID = auth.GenerateUniqueUserID()
			auth.SetUserIDCookie(w, userID)
		}
		// decoding data from the request body
		var resp models.Response
		var req models.Request
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		// check whether the link passed to the body is valid
		_, err = url.ParseRequestURI(req.URL)
		if err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		// making a compressed link from a regular link
		shortURL := urlshortener.ShortenerURL(cfg.BaseAddr)
		// saving a short link
		err = urlSaver.SaveURL(req.URL, shortURL, "", userID)
		if err != nil {
			log.Infof("%s: %s", op, err)
			if errors.Is(postgres.ErrURLExists, err) {
				shortURL = urlSaver.GetShortURL(req.URL)
				// passing the value in response
				resp.ShortURL = shortURL
				// setting the Content-Type header and sending the response
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
		// passing the value in response
		resp.ShortURL = shortURL
		// setting the Content-Type header and sending the response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
	}
}
