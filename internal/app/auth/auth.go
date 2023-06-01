package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
)

var secretKey = []byte("my-secret-key")

type URL struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

func GetUserIDFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		return "", err
	}

	parts := strings.Split(cookie.Value, "|")
	if len(parts) != 2 {
		return "", http.ErrNoCookie
	}

	userID := parts[0]
	signature := parts[1]

	expectedSignature := GenerateSignature(userID)
	if signature != expectedSignature {
		return "", http.ErrNoCookie
	}

	return userID, nil
}

func SetUserIDCookie(w http.ResponseWriter, userID string) {
	signature := GenerateSignature(userID)
	value := userID + "|" + signature

	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: value,
	})
}

func GenerateSignature(data string) string {
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func CreateNewUser() string {
	nextUserID++
	return strconv.Itoa(nextUserID)
}

func SaveURLToDatabase(userID string, shortURL string, longURL string) {
	urls := GetURLsFromDatabase(userID)
	urls = append(urls, URL{
		ShortURL:    shortURL,
		OriginalURL: longURL,
	})
	SaveURLsToDatabase(userID, urls)
}

var nextUserID = 0

var database = make(map[string][]URL)

func GetURLsFromDatabase(userID string) []URL {
	urls, ok := database[userID]
	if !ok {
		return []URL{}
	}
	return urls
}

func SaveURLsToDatabase(userID string, urls []URL) {
	database[userID] = urls
}
