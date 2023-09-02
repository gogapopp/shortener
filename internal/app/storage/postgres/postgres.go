// package postgres реализация интерфейса Storage для записи в файл
package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/gogapopp/shortener/internal/app/lib/globalstore"
	"github.com/gogapopp/shortener/internal/app/lib/models"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// ErrURLExists возращаемая если ссылка уже сохранена в хранилище
var ErrURLExists = errors.New("url exists")

// storage хранилище ссылок
type storage struct {
	db *sql.DB
}

// NewStorage создаёт хранилище storage
func NewStorage(databaseDSN string) (*storage, error) {
	const op = "storage.postgres.NewStorage"

	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS urls (
        id serial PRIMARY KEY,
        short_url TEXT,
        long_url TEXT,
        correlation_id TEXT,
		user_id TEXT,
		is_delete BOOLEAN
    );
    CREATE UNIQUE INDEX IF NOT EXISTS long_url_id ON urls(long_url);
`)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	return &storage{
		db: db,
	}, nil
}

// SaveURL сохраняет ссылки в хранилище
func (s *storage) SaveURL(longURL, shortURL, correlationID string, userID string) error {
	const op = "storage.postgres.SaveURL"
	var isDelete = false
	result, err := s.db.Exec("INSERT INTO urls (short_url, long_url, correlation_id, user_id, is_delete) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (long_url) DO NOTHING", shortURL, longURL, correlationID, userID, isDelete)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}
	if rowsAffected == 0 {
		return ErrURLExists
	}
	globalstore.GlobalStore.SaveURLToDatabase(userID, shortURL, longURL)
	return nil
}

// GetURL получает ссылку из хранилища
func (s *storage) GetURL(shortURL, userID string) (bool, string, error) {
	const op = "storage.postgres.GetURL"
	var longURL string
	var isDelete bool
	row := s.db.QueryRow("SELECT long_url, is_delete FROM urls WHERE short_url = $1", shortURL)
	err := row.Scan(&longURL, &isDelete)
	if err != nil {
		return false, "", fmt.Errorf("%s: %s", op, err)
	}
	return isDelete, longURL, nil
}

// Ping() проверяет подключение к базе данных
func (s *storage) Ping() (*sql.DB, error) {
	err := s.db.Ping()
	return s.db, err
}

// BatchInsertURL реализует batch запись скоращённых ссылок в хранилище
func (s *storage) BatchInsertURL(urls []models.BatchDatabaseResponse, userID string) error {
	const op = "storage.postgres.BatchInsertURL"
	var isDelete = false
	// собираем запрос
	query := "INSERT INTO urls (short_url, long_url, correlation_id, user_id, is_delete) VALUES "
	values := []interface{}{}

	for i, url := range urls {
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5)
		values = append(values, url.ShortURL, url.OriginalURL, url.CorrelationID, userID, isDelete)
		// ...
		globalstore.GlobalStore.SaveURLToDatabase(userID, url.ShortURL, url.OriginalURL)
	}
	// удаляем последнюю запятую и обновляем поля
	query = query[:len(query)-1]
	query = fmt.Sprintf("%sON CONFLICT (long_url) DO NOTHING", query)

	// выполняем запрос
	_, err := s.db.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}
	return nil
}

// GetShortURL получает короткую ссылку из хранилища
func (s *storage) GetShortURL(longURL string) string {
	var shortURL string
	row := s.db.QueryRow("SELECT short_url FROM urls WHERE long_url = $1", longURL)
	row.Scan(&shortURL)
	return shortURL
}

// GetUserURLs возвращает ссылки которые сохранял определённый пользователь
func (s *storage) GetUserURLs(userID string) ([]models.UserURLs, error) {
	// const op = "storage.postgres.GetUserURLs"
	// rows, err := s.db.Query("SELECT long_url, short_url FROM urls WHERE user_id = $1", userID)
	// if err != nil {
	// 	return nil, fmt.Errorf("%s: %s", op, err)
	// }
	// defer rows.Close()

	// var urls []models.UserURLs
	// for rows.Next() {
	// 	var url models.UserURLs
	// 	if err := rows.Scan(&url.OriginalURL, &url.ShortURL); err != nil {
	// 		return nil, fmt.Errorf("%s: %s", op, err)
	// 	}
	// 	urls = append(urls, url)
	// }

	// if err := rows.Err(); err != nil {
	// 	return nil, fmt.Errorf("%s: %s", op, err)
	// }

	// return urls, nil

	// получаем все сокращенные пользователем URL из базы данных
	urls := globalstore.GlobalStore.GetURLsFromDatabase(userID)
	return urls, nil
}

// SetDeleteFlag реализует логику удаления ссылок из хранилища
func (s *storage) SetDeleteFlag(IDs []string, userID string) error {
	const op = "storage.postgres.SetDeleteFlag"
	query := `
		UPDATE urls
		SET is_delete = true
		WHERE short_url = ANY($1) AND user_id = $2
	`
	if _, err := s.db.Exec(query, "{"+strings.Join(IDs, ",")+"}", userID); err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}
	return nil
}

// GetStats получаем кол-во юзеров и коротких ссылок в бд
func (s *storage) GetStats() (int, int, error) {
	const op = "storage.postgres.GetStats"
	statsQuery := "SELECT COUNT(short_url) AS short_url_count, COUNT(DISTINCT user_id) AS user_id_count FROM urls"
	var shortURLcount, userIDcount int
	err := s.db.QueryRow(statsQuery).Scan(&shortURLcount, &userIDcount)
	if err != nil {
		return 0, 0, fmt.Errorf("%s: %s", op, err)
	}
	return shortURLcount, userIDcount, nil
}
