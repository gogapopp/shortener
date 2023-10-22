// package redirect contains the GetURLGetterHandler handler
package redirect

import (
	"fmt"
	"net/http"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/network-server/middlewares/auth"
	"go.uber.org/zap"
)

// URLGetter defines the getURL method
type URLGetter interface {
	GetURL(shortURL, userID string) (bool, string, error)
}

// GetURLGetterHandler redirects the user to the link that corresponds to the abbreviated
func GetURLGetterHandler(log *zap.SugaredLogger, urlGetter URLGetter, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.GetURLGetterHandler"
		// get the userID from the context that was set by the middleware UserIdentity
		userID, err := auth.GetUserIDFromCookie(r)
		if err != nil {
			userID = auth.GenerateUniqueUserID()
			auth.SetUserIDCookie(w, userID)
		}
		url := fmt.Sprintf("http://%s%s", r.Host, r.URL.Path)
		// gets a link from the repository
		isDelete, longURL, err := urlGetter.GetURL(url, userID)
		if isDelete {
			http.Error(w, "url is deleted", http.StatusGone)
			return
		}
		if err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "url not found", http.StatusBadRequest)
			return
		}
		w.Header().Add("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
