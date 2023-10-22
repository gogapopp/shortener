// package auth contains the authentication code
// TODO: need to ref.
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// secret Key secret key but no longer secret
var secretKey = []byte("secret-key")

// nextUserID records the user's ID
var nextUserID = 0

// AuthMiddleware authenticates the user using cookies
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

// GetUserIDFromCookie gets user ID from cookie
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

// SetUserIDCookie sets a cookie to the user
func SetUserIDCookie(w http.ResponseWriter, userID string) {
	signature := GenerateSignature(userID)
	value := fmt.Sprintf("%s|%s", userID, signature)

	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: value,
		Path:  "/",
	})
}

// GenerateSignature creates an encrypted signature
func GenerateSignature(userID string) string {
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(userID))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// GenerateUniqueUserID creates a unique ID for the user
func GenerateUniqueUserID() string {
	nextUserID++
	return strconv.Itoa(nextUserID)
}