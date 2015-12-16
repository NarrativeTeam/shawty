package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockStorage struct{}

func (s *MockStorage) Save(url string) string {
	return "short-" + url
}

func (s *MockStorage) Load(token string) (string, error) {
	if strings.HasPrefix(token, "short-") {
		return token[len("short-"):len(token)], nil
	} else {
		return "", fmt.Errorf("No short-url found")
	}
}

func TestEncodeURL(t *testing.T) {
	cases := []struct {
		input      string
		statusCode int
		expected   string
	}{
		{
			`{"url": "http://foo.bar"}`,
			200,
			`{"url":"http://foo.bar","short_url":"short-http://foo.bar"}`,
		}, {
			`{}`,
			400,
			`{"error":"Missing required parameter 'url'"}`,
		}, {
			``,
			400,
			`{"error":"Unable to parse json-input"}`,
		},
	}
	handler := MainHandler(&MockStorage{})
	for _, c := range cases {
		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "", strings.NewReader(c.input))
		handler.ServeHTTP(resp, req)
		if resp.Code != c.statusCode {
			t.Errorf("Incorrect status-code: %v expected %v", resp.Code, c.statusCode)
		}
		if resp.Body.String() != c.expected {
			t.Errorf("Unexpected output: %v expected %v", resp.Body.String(), c.expected)
		}
	}
}

func TestRedirect(t *testing.T) {
	cases := []struct {
		path           string
		statusCode     int
		locationHeader string
	}{
		{"short-http://google.com", 301, "http://google.com"},
		{"", 404, ""},
	}

	handler := MainHandler(&MockStorage{})
	for _, c := range cases {
		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "http://shawty.io/"+c.path, nil)
		handler.ServeHTTP(resp, req)
		if resp.Code != c.statusCode {
			t.Errorf("Incorrect status-code: %v expected %v", resp.Code, c.statusCode)
		}
		locHeader := resp.HeaderMap.Get("Location")
		if locHeader != c.locationHeader {
			t.Errorf("Incorrect value in location-header: %v expected %v", locHeader, c.locationHeader)
		}
	}
}
