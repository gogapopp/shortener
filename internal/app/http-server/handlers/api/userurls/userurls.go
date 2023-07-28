package userurls

import (
	"encoding/json"
	"net/http"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/http-server/middlewares/auth"
	"github.com/gogapopp/shortener/internal/app/lib/models"
	"go.uber.org/zap"
)

//go:generate mockgen -source=redirect.go -destination=mocks/mock.go
type UserURLsGetter interface {
	GetUserURLs(userID string) ([]models.UserURLs, error)
}

func GetURLsHandler(log *zap.SugaredLogger, userURLsGetter UserURLsGetter, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.api.userurls.GetURLsHandler"
		// получаем userID из контекста который был установлен мидлвеером userIdentity
		userID, err := auth.GetUserIDFromCookie(r)
		if err != nil {
			userID = auth.GenerateUniqueUserID()
			auth.SetUserIDCookie(w, userID)
		}
		// получает ссылки из хранилища
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
		// устанавливаем заголовок Content-Type и отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		log.Info(userURLs)
		if err := json.NewEncoder(w).Encode(userURLs); err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
	}
}
