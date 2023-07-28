package redirect

import (
	"fmt"
	"net/http"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/http-server/middlewares/auth"
	"go.uber.org/zap"
)

//go:generate mockgen -source=redirect.go -destination=mocks/mock.go
type URLGetter interface {
	GetURL(shortURL, userID string) (string, error)
}

func GetURLGetterHandler(log *zap.SugaredLogger, urlGetter URLGetter, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.GetURLGetterHandler"
		// получаем userID из контекста который был установлен мидлвеером userIdentity
		userID, err := auth.GetUserIDFromCookie(r)
		if err != nil {
			userID = auth.GenerateUniqueUserID()
			auth.SetUserIDCookie(w, userID)
		}
		url := fmt.Sprintf("http://%s%s", r.Host, r.URL.Path)
		// получает ссылку из хранилища
		longURL, err := urlGetter.GetURL(url, userID)
		if err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "url not found", http.StatusBadRequest)
			return
		}
		// отправляем ответ
		w.Header().Add("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
