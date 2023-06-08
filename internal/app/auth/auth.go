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

	"github.com/gogapopp/shortener/internal/app/models"
)

var secretKey = []byte("my-secret-key")

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

var nextUserID = 0

// CreateNewUser создаёт userID для нового юзера
func CreateNewUser() string {
	nextUserID++
	return strconv.Itoa(nextUserID)
}

type Store struct {
	mu       sync.Mutex
	database map[string][]models.URL
}

// объявляем глобальную переменную, чтоб записывать в неё данные из пакетов handlers и concurrency
var GlobalStore *Store

func init() {
	GlobalStore = &Store{
		database: make(map[string][]models.URL),
	}
}

// SaveURLToDatabase получает userID и соответсвенно ему сохраняет models.URL
func (s *Store) SaveURLToDatabase(userID string, shortURL string, longURL string) {
	urls := s.GetURLsFromDatabase(userID)
	urls = append(urls, models.URL{
		ShortURL:    shortURL,
		OriginalURL: longURL,
		DeleteFlag:  false,
	})
	s.SaveURLsToDatabase(userID, urls)
}

// GetURLsFromDatabase получает структуры models.URL соответственно переданному userUD
func (s *Store) GetURLsFromDatabase(userID string) []models.URL {
	return s.database[userID]
}

// SetDeleteFlag принимает userID, проходит по структуре и ищет соответсвие между models.URL юзера и переданным ссылкам для удаления в теле метода GET
func (s *Store) SetDeleteFlag(userID string, shortURL string, baseAddr string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	urls := s.GetURLsFromDatabase(userID)
	for k, url := range urls {
		if url.ShortURL == fmt.Sprint(baseAddr+shortURL) {
			urls[k].DeleteFlag = true
			break
		}
	}
	s.SaveURLsToDatabase(userID, urls)
}

// SaveURLsToDatabase сохраняет структуру, которая была обработана функцией SetDeleteFlag
func (s *Store) SaveURLsToDatabase(userID string, urls []models.URL) {
	s.database[userID] = urls
}
