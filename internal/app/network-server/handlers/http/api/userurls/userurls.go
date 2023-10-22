package userurls

import (
	"encoding/json"
	"net/http"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/lib/models"
	"github.com/gogapopp/shortener/internal/app/network-server/middlewares/auth"
	"go.uber.org/zap"
)

// UserURLsGetter defines the getUser URLs method
type UserURLsGetter interface {
	GetUserURLs(userID string) ([]models.UserURLs, error)
}

// GetURLsHandler returns all shortened user links
func GetURLsHandler(log *zap.SugaredLogger, userURLsGetter UserURLsGetter, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.api.userurls.GetURLsHandler"
		// get the userID from the context that was set by the middleware UserIdentity
		userID, err := auth.GetUserIDFromCookie(r)
		if err != nil {
			userID = auth.GenerateUniqueUserID()
			auth.SetUserIDCookie(w, userID)
		}
		// gets links from the repository
		userURLs, err := userURLsGetter.GetUserURLs(userID)
		if err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "url not found", http.StatusBadRequest)
			return
		}
		if len(userURLs) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// setting the Content-Type header and sending the response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(userURLs); err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
	}
}
