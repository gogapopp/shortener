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
var CookieURLSMap = make(map[string]string)
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
	userID, err := auth.GetUserIDFromCookie(r)
	if err != nil {
		userID = auth.CreateNewUser()
		auth.SetUserIDCookie(w, userID)
	}
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
	// auth
	auth.SaveURLToDatabase(userID, fmt.Sprint("http://"+parsedURL.Host+parsedURL.Path), mainURL)
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
	userID, err := auth.GetUserIDFromCookie(r)
	if err != nil {
		userID = auth.CreateNewUser()
		auth.SetUserIDCookie(w, userID)
	}
	id := r.URL.Path
	// проверяем есть ли значение в мапе
	urls := auth.GetURLsFromDatabase(userID)
	for _, url := range urls {
		if url.DeleteFlag == true {
			w.WriteHeader(http.StatusGone)
			return
		}
	}
	fmt.Println(URLSMap)
	// проверяем "удаления ли ссылка"
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
	userID, err := auth.GetUserIDFromCookie(r)
	if err != nil {
		userID = auth.CreateNewUser()
		auth.SetUserIDCookie(w, userID)
	}
	var responseHeader = http.StatusCreated
	ctx := r.Context()
	// декодируем данные из тела запроса
	var req models.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusMethodNotAllowed)
		return
	}
	// проверяем имеет ли body в себе url ссылку
	_, err = url.ParseRequestURI(req.URL)
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
	// auth
	auth.SaveURLToDatabase(userID, fmt.Sprint("http://"+parsedURL.Host+parsedURL.Path), req.URL)
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
	userID, err := auth.GetUserIDFromCookie(r)
	if err != nil {
		userID = auth.CreateNewUser()
		auth.SetUserIDCookie(w, userID)
	}
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
			// auth
			auth.SaveURLToDatabase(userID, fmt.Sprint("http://"+parsedURL.Host+parsedURL.Path), req[k].OriginalURL)

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

func GetURLs(w http.ResponseWriter, r *http.Request) {
	// проверяем наличие куки с идентификатором пользователя
	userID, err := auth.GetUserIDFromCookie(r)
	if err != nil {
		userID = auth.CreateNewUser()
		auth.SetUserIDCookie(w, userID)
	}
	// получаем все сокращенные пользователем URL из базы данных
	urls := auth.GetURLsFromDatabase(userID)

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// устанавливаем заголовок Content-Type и отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(urls); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func DeleteShortURLs(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromCookie(r)
	if err != nil {
		userID = auth.CreateNewUser()
		auth.SetUserIDCookie(w, userID)
	}
	// читаем тело реквеста
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusMethodNotAllowed)
		return
	}
	log.Println("ids", string(body))
	var IDs []string
	err = json.Unmarshal(body, &IDs)
	if err != nil {
		log.Println("err", err)

	}
	log.Println("ids", IDs)
	// создаем канал для отправки идентификаторов на удаление
	idCh := make(chan string)

	// запускаем горутины для обработки удаления URL
	go deleteURLs(idCh, userID, fmt.Sprint("http://"+r.Host))
	go deleteURLs(idCh, userID, fmt.Sprint("http://"+r.Host))

	// отправляем идентификаторы в канал idCh
	for _, id := range IDs {
		idCh <- id
	}

	// закрываем канал idCh
	close(idCh)
	w.WriteHeader(http.StatusAccepted)
}

func deleteURLs(idCh chan string, userID string, baseAddr string) {
	log.Println("deleteURLs")
	for id := range idCh {
		log.Println("SetDeleteFlag")
		auth.SetDeleteFlag(userID, id, true, baseAddr)
	}
}
