package subnet

import (
	"net"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// SubnetMiddleware checks the client's IP address, whether it is included in a trusted subnet
func SubnetMiddleware(log *zap.SugaredLogger, subnet string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log.Info("subnet middleware enabled")
		fn := func(w http.ResponseWriter, r *http.Request) {
			if subnet == "" {
				http.Error(w, "untrusted ip", http.StatusForbidden)
				return
			}
			// parse subnet to IP
			subnetIP := net.ParseIP(subnet)
			// look at the X-Real-IP header
			ipStr := r.Header.Get("X-Real-IP")
			ip := net.ParseIP(ipStr)
			// if the header is empty
			if ip == nil {
				// if the X-Real-IP header is empty, try X-Forwarded-For
				ips := r.Header.Get("X-Forwarded-For")
				// separating the chain of addresses
				ipStrings := strings.Split(ips, ",")
				// interested only in the first
				ipStr = ipStrings[0]
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
