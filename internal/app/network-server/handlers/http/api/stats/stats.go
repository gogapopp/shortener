package stats

import (
	"encoding/json"
	"net/http"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/lib/models"
	"github.com/gogapopp/shortener/internal/app/network-server/middlewares/auth"
	"go.uber.org/zap"
)

// StatsGetter defines the GetStats method
type StatsGetter interface {
	GetStats() (int, int, error)
}

// GetStat returns all abbreviated links and the number of users
func GetStat(log *zap.SugaredLogger, statsGetter StatsGetter, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.api.userurls.GetURLsHandler"
		// getting the userID from the context that was set by the middleware UserIdentity
		userID, err := auth.GetUserIDFromCookie(r)
		if err != nil {
			userID = auth.GenerateUniqueUserID()
			auth.SetUserIDCookie(w, userID)
		}
		_ = userID
		// gets statistics from the repository from the repository
		shortURLcount, userIDcount, err := statsGetter.GetStats()
		if err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		// forming a response
		resp := models.Stasts{
			URLs:    shortURLcount,
			UserIDs: userIDcount,
		}
		// setting the Content-Type header and sending the response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
	}
}
