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

type APIError struct {
	Msg        string `json:"error"`
	statusCode int
}

func (apiErr APIError) writeTo(w http.ResponseWriter) {
	out, err := json.Marshal(apiErr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(apiErr.statusCode)
		w.Write(out)
	}
}

func MainHandler(storage storages.IStorage) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "GET" {
			// We need to remove the / in the beginning of path
			code := r.URL.Path[1:]

			url, err := storage.Load(code)
			if err != nil {
				APIError{"No url found for token", http.StatusNotFound}.writeTo(w)
				return
			}
			http.Redirect(w, r, string(url), 301)
		} else if r.Method == "POST" {
			dec := json.NewDecoder(r.Body)
			var data APIURL
			err := dec.Decode(&data)
			if err != nil {
				APIError{"Unable to parse json-input", http.StatusBadRequest}.writeTo(w)
				return
			}

			if data.Url != "" {
				token, err := storage.Save(data.Url)
				if err != nil {
					APIError{"Internal server error", http.StatusInternalServerError}.writeTo(w)
					return
				}
				shortUrl := *r.URL
				shortUrl.Path = token
				data.ShortUrl = shortUrl.String()
				out, err := json.Marshal(data)
				if err != nil {
					APIError{"Internal server error", http.StatusInternalServerError}.writeTo(w)
					return
				}
				w.Write(out)
				return
			} else {
				APIError{"Missing required parameter 'url'", http.StatusBadRequest}.writeTo(w)
				return
			}
		}

	}

	return http.HandlerFunc(handleFunc)
}
