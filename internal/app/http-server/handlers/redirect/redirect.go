package redirect

import (
	"fmt"
	"net/http"

	"github.com/gogapopp/shortener/internal/app/config"
	"go.uber.org/zap"
)

type URLGetter interface {
	GetURL(longURL string) (string, error)
}

func GetURLGetterHandler(log *zap.SugaredLogger, urlGetter URLGetter, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.GetURLGetterHandler"
		url := fmt.Sprintf("http://%s%s", r.Host, r.URL.Path)
		// получает ссылку из хранилища
		shortURL, err := urlGetter.GetURL(url)
		if err != nil {
			log.Info(fmt.Sprintf("%s: %s", op, err))
			http.Error(w, "url not found", http.StatusBadRequest)
			return
		}
		// отправляем ответ
		w.Header().Add("Location", shortURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
