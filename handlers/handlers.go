// Package handlers provides HTTP request handlers.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NarrativeTeam/shawty/storages"
)

type APIURL struct {
	Url      string `json:"url"`
	ShortUrl string `json:"short_url"`
}

func MainHandler(storage storages.IStorage) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			code := r.URL.Path

			url, err := storage.Load(code)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("URL Not Found. Error: " + err.Error() + "\n"))
				return
			}

			http.Redirect(w, r, string(url), 301)
		} else if r.Method == "POST" {
			dec := json.NewDecoder(r.Body)
			var data APIURL
			err := dec.Decode(&data)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if data.Url != "" {
				data.ShortUrl = storage.Save(data.Url)
				out, err := json.Marshal(data)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Write(out)
				return
			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

	}

	return http.HandlerFunc(handleFunc)
}
