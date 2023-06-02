package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var secretKey = []byte("my-secret-key")

type URL struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
	DeleteFlag  bool
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
		DeleteFlag:  false,
	})
	SaveURLsToDatabase(userID, urls)
}

func GetURLsFromDatabase(userID string) []URL {
	return database[userID]
}

var nextUserID = 0

var (
	database = make(map[string][]URL)
	mu       sync.Mutex
)

func SetDeleteFlag(userID string, shortURL string, deleteFlag bool, baseAddr string) {
	mu.Lock()
	defer mu.Unlock()
	urls := GetURLsFromDatabase(userID)
	for k, url := range urls {
		if url.ShortURL == fmt.Sprint(baseAddr+shortURL) {
			urls[k].DeleteFlag = deleteFlag
			break
		}
	}
	SaveURLsToDatabase(userID, urls)
}

func SaveURLsToDatabase(userID string, urls []URL) {
	database[userID] = urls
}
