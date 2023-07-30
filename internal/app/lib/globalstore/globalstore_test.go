package globalstore_test

import (
	"testing"

	"github.com/gogapopp/shortener/internal/app/lib/globalstore"
)

func TestSaveURLToDatabase(t *testing.T) {
	userID := "testUser"
	shortURL := "shortURL"
	longURL := "longURL"

	globalstore.GlobalStore.SaveURLToDatabase(userID, shortURL, longURL)

	urls := globalstore.GlobalStore.GetURLsFromDatabase(userID)

	if len(urls) != 1 {
		t.Errorf("expected 1 URL, got %d", len(urls))
	}

	if urls[0].ShortURL != shortURL {
		t.Errorf("expected ShortURL to be %s, got %s", shortURL, urls[0].ShortURL)
	}

	if urls[0].OriginalURL != longURL {
		t.Errorf("expected OriginalURL to be %s, got %s", longURL, urls[0].OriginalURL)
	}
}
