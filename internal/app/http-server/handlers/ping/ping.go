// package ping содержит код хендлера GetPingDBHandler
package ping

import (
	"database/sql"
	"net/http"

	"github.com/gogapopp/shortener/internal/app/config"
	"go.uber.org/zap"
)

// DBPinger определяет метод Ping()
type DBPinger interface {
	Ping() (*sql.DB, error)
}

// GetPingDBHandler проверяет соединение с БД
func GetPingDBHandler(log *zap.SugaredLogger, dbPinger DBPinger, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.ping.GetPingDBHandler"
		_, err := dbPinger.Ping()
		if err != nil {
			log.Infof("%s: %s", op, err)
			http.Error(w, "error ping DB", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	}
}
