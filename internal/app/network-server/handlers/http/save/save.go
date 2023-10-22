// package save contains the PostSaveHandler handler code
package save

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gogapopp/shortener/internal/app/config"
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

// PostSaveHandler accepts a link as a string and returns an abbreviated
func PostSaveHandler(log *zap.SugaredLogger, urlSaver URLSaver, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.PostSaveHandler"
		// get the userID from the context that was set by the middleware UserIdentity
		userID, err := auth.GetUserIDFromCookie(r)
		if err != nil {
			userID = auth.GenerateUniqueUserID()
			auth.SetUserIDCookie(w, userID)
		}
		// reading the body of the request
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		bodyURL := string(body)
		// check whether the link passed to the body is valid
		_, err = url.ParseRequestURI(bodyURL)
		if err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		log.Infof("%s", body)
		// making a compressed link from a regular link
		shortURL := urlshortener.ShortenerURL(cfg.BaseAddr)
		// saving a short link
		err = urlSaver.SaveURL(bodyURL, shortURL, "", userID)
		if err != nil {
			log.Infof("%s: %s", op, err)
			if errors.Is(postgres.ErrURLExists, err) {
				shortURL = urlSaver.GetShortURL(bodyURL)
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusConflict)
				fmt.Fprint(w, shortURL)
				return
			}
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		// sending a response
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, shortURL)
	}
}
