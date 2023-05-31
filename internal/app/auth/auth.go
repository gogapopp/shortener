package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var SecretKey = []byte("supersecretkey")

type User struct {
	ID   int
	URLs []URL
}

type URL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

var Users = make(map[int]User)

func GetUserIDFromCookie(w http.ResponseWriter, r *http.Request) (int, error) {
	c, err := r.Cookie("user_id")
	if err != nil {
		return createNewUser(w)
	}

	parts := strings.Split(c.Value, "|")
	if len(parts) != 2 {
		return createNewUser(w)
	}

	userID, err := strconv.Atoi(parts[0])
	if err != nil {
		return createNewUser(w)
	}

	mac := hmac.New(sha256.New, SecretKey)
	mac.Write([]byte(parts[0]))
	expectedMAC := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(parts[1]), []byte(expectedMAC)) {
		return createNewUser(w)
	}

	return userID, nil
}

func createNewUser(w http.ResponseWriter) (int, error) {
	userID := len(Users) + 1
	Users[userID] = User{ID: userID}
	mac := hmac.New(sha256.New, SecretKey)
	mac.Write([]byte(strconv.Itoa(userID)))
	macStr := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	c := &http.Cookie{
		Name:  "user_id",
		Value: strconv.Itoa(userID) + "|" + macStr,
	}
	http.SetCookie(w, c)
	return userID, nil
}

func AddURL(userID int, shortURL string, originalURL string) {
	user, ok := Users[userID]
	if !ok {
		fmt.Println("такого юзера не существует")
	}
	// добавляем новый URL в список URL пользователя
	user.URLs = append(user.URLs, URL{ShortURL: shortURL, OriginalURL: originalURL})
	Users[userID] = user
}

func DeleteURLs(userID int) {
	user, ok := Users[userID]
	if ok {
		user.URLs = nil
		Users[userID] = user
	}
}
