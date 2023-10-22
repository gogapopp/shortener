// package url delete contains the handler code DeleteHandler
package urlsdelete

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/lib/concurrency"
	"github.com/gogapopp/shortener/internal/app/network-server/middlewares/auth"
	"go.uber.org/zap"
)

// URLDeleter defines the SetDeleteFlag method
type URLDeleter interface {
	SetDeleteFlag(IDs []string, userID string) error
}

// DeleteHandler эндлер кото-рый распознает мас - сив и идентификатор сокра - ченных строк для удаления
func DeleteHandler(log *zap.SugaredLogger, urlDeleter URLDeleter, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.api.urlsdelete.PostSaveHandler"
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
		var IDs []string
		err = json.Unmarshal(body, &IDs)
		if err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		reqURL := fmt.Sprintf("http://%s", r.Host)
		go concurrency.ProcessIDs(IDs, reqURL, urlDeleter, userID)
		w.WriteHeader(http.StatusAccepted)
	}
}
