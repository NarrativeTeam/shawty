package storages

import (
	"os"
	"testing"
)

func GetTestPostgres() *Postgres {
	// We need to fetch the actuall address for postgres, it seems like the driver doesn't respect /etc/hosts
	host := os.Getenv("POSTGRES_PORT_5432_TCP_ADDR")
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
