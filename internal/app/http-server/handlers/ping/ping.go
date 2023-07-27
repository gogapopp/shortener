package ping

import (
	"fmt"
	"net/http"

	"github.com/gogapopp/shortener/internal/app/config"
	"go.uber.org/zap"
)

type DBPinger interface {
	Ping() error
}

func GetPingDBHandler(log *zap.SugaredLogger, dbPinger DBPinger, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.ping.GetPingDBHandler"
		err := dbPinger.Ping()
		if err != nil {
			log.Info(fmt.Sprintf("%s: %s", op, err))
			http.Error(w, "error ping DB", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	}
}
