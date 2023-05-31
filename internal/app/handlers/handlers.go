package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gogapopp/shortener/internal/app/auth"
	"github.com/gogapopp/shortener/internal/app/encryptor"
	"github.com/gogapopp/shortener/internal/app/models"
	"github.com/gogapopp/shortener/internal/app/storage"
)

var URLSMap = storage.URLSMap
var writeToFile bool = false
var writeToDatabase bool = true

var ErrParseURL = errors.New("ErrParseURL")

func WriteToFile(b bool) {
	writeToFile = b
}
func WriteToDatabase(c bool) {
	writeToDatabase = c
}

// PostShortURL получает ссылку в body и присваивает ей уникальный ключ, значение хранит в мапе "key": "url"
func PostShortURL(w http.ResponseWriter, r *http.Request) {
	var responseHeader = http.StatusCreated
	ctx := r.Context()
	// читаем тело реквеста
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusMethodNotAllowed)
		return
	}
	mainURL := string(body)
	// делаем из обычной ссылки сжатую
	shortURL := encryptor.ShortenerURL(mainURL)
	// проверяем имеет ли body в себе url ссылку
	_, err = url.ParseRequestURI(mainURL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	// получаем url path от новой сжатой ссылки /{id} и заполняем мапу
	parsedURL, err := url.Parse(shortURL)
	if err != nil {
		log.Fatal(err)
	}
	// сохраняем значение в мапу
	URLSMap[parsedURL.Path] = mainURL
	// получаем request адресс с которого происходит запрос
	host := r.Host
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	requestURL := fmt.Sprintf("%s://%s%s", scheme, host, parsedURL.Path)
	// сохраняем в базу данных
	if writeToDatabase {
		// делаем запись в виде id (primary key), short_url, long_url
		err := storage.InsertURL(ctx, requestURL, URLSMap[parsedURL.Path], "")
		if err != nil {
			if errors.Is(err, storage.ErrConflict) {
				// собираем shortURL
				shortURL := storage.FindShortURL(ctx, mainURL)
				responseHeader = http.StatusConflict
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(responseHeader)
				fmt.Fprint(w, shortURL)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	// сохраняем в файл
	if writeToFile {
		storage.UUIDCounter++
		storage.ShortURLStorage = append(storage.ShortURLStorage, models.ShortURL{
			UUID:        storage.UUIDCounter,
			ShortURL:    parsedURL.Path,
			OriginalURL: URLSMap[parsedURL.Path],
		})
		storage.Save()
	}
	// отправляем ответ
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(responseHeader)
	fmt.Fprint(w, shortURL)
}

// GetPingDatabase пингует PostgreSQL
func GetPingDatabase(w http.ResponseWriter, r *http.Request) {
	if db := storage.DB(); db == nil {
		http.Error(w, "База данных не инициализирована", http.StatusInternalServerError)
		return
	}
	if err := storage.DB().Ping(); err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetHandleURL проверяет валидная ссылка или нет, если валидная то редиректит по адрессу
func GetHandleURL(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path
	// проверяем есть ли значение в мапе
	fmt.Println(URLSMap)
	if _, ok := URLSMap[id]; ok {
		w.Header().Add("Location", URLSMap[id])
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	} else {
		http.Error(w, "Link does not exist", http.StatusBadRequest)
	}
}

// PostJSONHandler принимает url: url и возвращает result: shortUrl
func PostJSONHandler(w http.ResponseWriter, r *http.Request) {
	var responseHeader = http.StatusCreated
	ctx := r.Context()
	// декодируем данные из тела запроса
	var req models.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusMethodNotAllowed)
		return
	}
	// проверяем имеет ли body в себе url ссылку
	_, err := url.ParseRequestURI(req.URL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	// "сжимаем" строку
	shortURL := encryptor.ShortenerURL(req.URL)
	// получаем url path от новой сжатой ссылки /{id} и заполняем мапу
	parsedURL, err := url.Parse(shortURL)
	if err != nil {
		log.Fatal(err)
	}
	// сохраняем значение в мапу
	URLSMap[parsedURL.Path] = req.URL
	// передаём значение в ответ
	var resp models.Response
	resp.ShortURL = shortURL

	// сохраняем в базу данных
	if writeToDatabase {
		// делаем запись в виде id (primary key), short_url, long_url
		err := storage.InsertURL(ctx, shortURL, req.URL, "")
		if err != nil {
			if errors.Is(err, storage.ErrConflict) {
				responseHeader = http.StatusConflict
				shortURL := storage.FindShortURL(ctx, req.URL)
				resp.ShortURL = shortURL
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	// сохраняем в файл
	if writeToFile {
		storage.UUIDCounter++
		storage.ShortURLStorage = append(storage.ShortURLStorage, models.ShortURL{
			UUID:        storage.UUIDCounter,
			ShortURL:    parsedURL.Path,
			OriginalURL: URLSMap[parsedURL.Path],
		})
		storage.Save()
	}

	// устанавливаем заголовок Content-Type и отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseHeader)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// PostBatchJSONhHandler
func PostBatchJSONhHandler(w http.ResponseWriter, r *http.Request) {
	var responseHeader = http.StatusCreated
	var req []models.BatchRequest
	var resp []models.BatchResponse
	var databaseResp []models.BatchDatabaseResponse
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusMethodNotAllowed)
		return
	}
	if writeToDatabase {
		// начинаем проходить по реквесту
		for k := range req {
			// проверяем является ли переданное значение ссылкой
			_, err := url.ParseRequestURI(req[k].OriginalURL)
			if err != nil {
				http.Error(w, "Invalid URL", http.StatusBadRequest)
				return
			}
			// "сжимаем" ссылку
			BatchShortURL := encryptor.ShortenerURL(req[k].OriginalURL)
			// получаем url path от новой сжатой ссылки /{id} и заполняем мапу
			parsedURL, err := url.Parse(BatchShortURL)
			if err != nil {
				log.Fatal(err)
			}
			// сохраняем значение в мапу
			URLSMap[parsedURL.Path] = req[k].OriginalURL

			databaseResp = append(databaseResp, models.BatchDatabaseResponse{
				ShortURL:      BatchShortURL,
				OriginalURL:   req[k].OriginalURL,
				CorrelationID: req[k].CorrelationID,
			})

			resp = append(resp, models.BatchResponse{
				ShortURL:      BatchShortURL,
				CorrelationID: req[k].CorrelationID,
			})
		}
		err := storage.BatchInsertURL(ctx, databaseResp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if writeToFile {
		for k := range req {
			// проверяем является ли переданное значение ссылкой
			_, err := url.ParseRequestURI(req[k].OriginalURL)
			if err != nil {
				http.Error(w, "Invalid URL", http.StatusBadRequest)
				return
			}
			// "сжимаем" ссылку
			shortURL := encryptor.ShortenerURL(req[k].OriginalURL)
			// получаем url path от новой сжатой ссылки /{id} и заполняем мапу
			parsedURL, err := url.Parse(shortURL)
			if err != nil {
				log.Fatal(err)
			}
			// сохраняем значение в мапу
			URLSMap[parsedURL.Path] = req[k].OriginalURL
			storage.UUIDCounter++
			storage.ShortURLStorage = append(storage.ShortURLStorage, models.ShortURL{
				UUID:        storage.UUIDCounter,
				ShortURL:    parsedURL.Path,
				OriginalURL: URLSMap[parsedURL.Path],
			})
			storage.Save()
		}
	}
	// устанавливаем заголовок Content-Type и отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseHeader)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func GetJSONURLS(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromCookie(w, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	auth.DeleteURLs(userID)
	for k, v := range URLSMap {
		auth.AddURL(userID, fmt.Sprint("http://"+r.Host+r.URL.Host+k), v)
	}

	user, ok := auth.Users[userID]
	if !ok || len(user.URLs) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	data, err := json.Marshal(user.URLs)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
