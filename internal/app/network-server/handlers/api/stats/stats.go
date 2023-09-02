package stats

import (
	"encoding/json"
	"net/http"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/lib/models"
	"github.com/gogapopp/shortener/internal/app/network-server/middlewares/auth"
	"go.uber.org/zap"
)

// StatsGetter определяет метод GetStats
type StatsGetter interface {
	GetStats() (int, int, error)
}

// GetURLsHandler возвращает все сокращённые ссылки юзера
func GetStat(log *zap.SugaredLogger, statsGetter StatsGetter, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.api.userurls.GetURLsHandler"
		// получаем userID из контекста который был установлен мидлвеером userIdentity
		userID, err := auth.GetUserIDFromCookie(r)
		if err != nil {
			userID = auth.GenerateUniqueUserID()
			auth.SetUserIDCookie(w, userID)
		}
		_ = userID
		// получает статистику из хранилища из хранилища
		shortURLcount, userIDcount, err := statsGetter.GetStats()
		if err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		// формируем ответ
		resp := models.Stasts{
			URLs:    shortURLcount,
			UserIDs: userIDcount,
		}
		// устанавливаем заголовок Content-Type и отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
	}
}
