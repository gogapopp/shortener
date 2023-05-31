package auth

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
)

const SECRET_KEY = "supersecretkey"

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

	tokenString := c.Value
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil || !token.Valid {
		return createNewUser(w)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userIDstring := claims["userID"].(string)
		userID, err := strconv.Atoi(userIDstring)
		if err != nil {
			return createNewUser(w)
		}
		return userID, nil
	}
	return createNewUser(w)
}

func createNewUser(w http.ResponseWriter) (int, error) {
	userID := len(Users) + 1
	Users[userID] = User{ID: userID}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": strconv.Itoa(userID),
	})
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return -1, err
	}
	c := &http.Cookie{
		Name:  "user_id",
		Value: tokenString,
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
