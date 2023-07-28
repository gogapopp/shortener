package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"

	"go.uber.org/zap"
)

var secretKey = []byte("secret-key")
var userIDCounter uint64

func AuthMiddleware(log *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log.Info("auth middleware enabled")
		fn := func(w http.ResponseWriter, r *http.Request) {
			userID, err := GetUserIDFromCookie(r)
			if err != nil {
				userID = GenerateUniqueUserID()
				SetUserIDCookie(w, userID)
			}
			SetUserIDCookie(w, userID)

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func GetUserIDFromCookie(r *http.Request) (string, error) {
	const op = "middlewares.auth.GetUserIDFromCookie"

	cookie, err := r.Cookie("user_id")
	if err != nil {
		return "", fmt.Errorf("%s: %s", op, err)
	}

	parts := strings.Split(cookie.Value, "|")
	if len(parts) != 2 {
		return "", http.ErrNoCookie
	}

	userID := parts[0]
	signature := parts[1]

	exceptedSignature := GenerateSignature(userID)
	if signature != exceptedSignature {
		return "", http.ErrNoCookie
	}

	return userID, nil
}

func SetUserIDCookie(w http.ResponseWriter, userID string) {
	signature := GenerateSignature(userID)
	value := fmt.Sprintf("%s|%s", userID, signature)

	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: value,
	})
}

func GenerateSignature(userID string) string {
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(userID))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func GenerateUniqueUserID() string {
	atomic.AddUint64(&userIDCounter, 1)
	return "user_" + strconv.FormatUint(userIDCounter, 10)
}
