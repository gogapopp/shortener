package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gogapopp/shortener/internal/app/encryptor"
	"github.com/gogapopp/shortener/internal/app/models"
	"github.com/gogapopp/shortener/internal/app/storage"
)

var URLSMap = storage.URLSMap
var writeToFile bool = false
var writeToDatabase bool = false

func WriteToFile(b bool) {
	writeToFile = b
}
func WriteToDatabase(c bool) {
	writeToDatabase = c
}

// PostShortURL получает ссылку в body и присваивает ей уникальный ключ, значение хранит в мапе "key": "url"
func PostShortURL(w http.ResponseWriter, r *http.Request) {
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

	// сохраняем в базу данных
	if writeToDatabase {
		// делаем запись в виде id (primary key), short_url, long_url
		err := storage.InsertURL(ctx, parsedURL.Path, URLSMap[parsedURL.Path], "")
		if err != nil {
			log.Fatal(err)
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
		if err := storage.Save(); err != nil {
			log.Fatal(err)
		}
	}

	// отправляем ответ
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
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
			log.Fatal(err)
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
		if err := storage.Save(); err != nil {
			log.Fatal(err)
		}
	}

	// устанавливаем заголовок Content-Type и отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// PostBatchJSONhHandler
func PostBatchJSONhHandler(w http.ResponseWriter, r *http.Request) {
	var req []models.BatchRequest
	var resp []models.BatchResponse
	if writeToDatabase {
		ctx := r.Context()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "error decoding request body", http.StatusMethodNotAllowed)
			return
		}

		// начинаем проходить по реквесту
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
			// делаем запись в базу данных
			if err := storage.InsertURL(ctx, shortURL, req[k].OriginalURL, req[k].CorrelationID); err != nil {
				fmt.Println(err)
				http.Error(w, "error insert to database", http.StatusMethodNotAllowed)
				return
			}

			resp = append(resp, models.BatchResponse{
				ShortURL:      shortURL,
				CorrelationID: req[k].CorrelationID,
			})
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
			if err := storage.Save(); err != nil {
				log.Fatal(err)
			}
		}
	}

	// устанавливаем заголовок Content-Type и отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
