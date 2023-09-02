package subnet

import (
	"net"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// SubnetMiddleware проверяет IP адресс клиента, входит ли он в доверенную подсеть
func SubnetMiddleware(log *zap.SugaredLogger, subnet string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log.Info("subnet middleware enabled")
		fn := func(w http.ResponseWriter, r *http.Request) {
			if subnet == "" {
				http.Error(w, "untrusted ip", http.StatusForbidden)
				return
			}
			// парсим subnet в IP
			subnetIP := net.ParseIP(subnet)
			// смотрим заголовок X-Real-IP
			ipStr := r.Header.Get("X-Real-IP")
			ip := net.ParseIP(ipStr)
			// если заголовок пуст
			if ip == nil {
				// если заголовок X-Real-IP пуст, пробуем X-Forwarded-For
				ips := r.Header.Get("X-Forwarded-For")
				// разделяем цепочку адресов
				ipStrings := strings.Split(ips, ",")
				// интересует только первый
				ipStr = ipStrings[0]
				// парсим
				ip = net.ParseIP(ipStr)
			}
			if ip == nil {
				http.Error(w, "untrusted ip", http.StatusForbidden)
				return
			} else if !ip.Equal(subnetIP) {
				http.Error(w, "untrusted ip", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
