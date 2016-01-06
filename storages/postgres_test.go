package storages

import (
	"os"
	"testing"
)

func GetTestPostgres() *Postgres {
	host := os.Getenv("POSTGRES_HOST")
	pg, err := NewPostgres(host, "postgres", "postgres", "postgres", false)
	if err != nil {
		panic(err)
	}
	pg.Teardown() // Remove old data
	pg.Setup()
	return pg
}

func TestStore(t *testing.T) {
	cases := []struct {
		url   string
		token string
	}{
		{"http://google.com", "WodTB2"},
	}
	for _, c := range cases {
		storage := GetTestPostgres()
		token, err := storage.Save(c.url)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if token != c.token {
			t.Errorf("Unexpected token: %v expected %v", token, c.token)
		}
	}
}

func TestStoreRestore(t *testing.T) {
	type tokenUrl struct {
		token string
		url   string
	}
	tokenUrls := []tokenUrl{}

	baseUrls := []string{"http://google.com", "http://facebook.com", "http://getnarrative.com", "http://narrativeapp.com/foobar"}
	storage := GetTestPostgres()
	for _, url := range baseUrls {
		token, err := storage.Save(url)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		tokenUrls = append(tokenUrls, tokenUrl{token, url})
	}

	for _, c := range tokenUrls {
		url, err := storage.Load(c.token)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if url != c.url {
			t.Errorf("Got unexpected url from storage: %v expected %v", url, c.url)
		}
	}
}

func TestRestoreStats(t *testing.T) {
	storage := GetTestPostgres()
	token, err := storage.Save("http://google.com")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	storage.Load(token)
	storage.Load(token)
	storage.Load(token)
	storage.Load(token)

	stats, err := storage.GetStats(token)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if stats.hits != 4 {
		t.Errorf("Unexpected number of hits: %v expected 4", stats.hits)
	}

}
