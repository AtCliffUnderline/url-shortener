package app_test

import (
	"github.com/AtCliffUnderline/url-shortener/internal/app"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUrlShortenerWrongMethod(t *testing.T) {
	request := httptest.NewRequest(http.MethodPatch, "/", nil)
	w := httptest.NewRecorder()
	h := http.HandlerFunc(app.UrlShortenerHandler)
	h.ServeHTTP(w, request)
	result := w.Result()
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
}

func TestUrlShortenerSymbolicId(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/asd", nil)
	w := httptest.NewRecorder()
	h := http.HandlerFunc(app.UrlShortenerHandler)
	h.ServeHTTP(w, request)
	result := w.Result()
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
}

func TestUrlShortenerWrongId(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/123", nil)
	w := httptest.NewRecorder()
	h := http.HandlerFunc(app.UrlShortenerHandler)
	h.ServeHTTP(w, request)
	result := w.Result()
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
}

func TestUrlShorteningSuccessful(t *testing.T) {
	// step 1: shortening link
	body := strings.NewReader("https://yc.gl")
	request := httptest.NewRequest(http.MethodPost, "/", body)
	w := httptest.NewRecorder()
	h := http.HandlerFunc(app.UrlShortenerHandler)
	h.ServeHTTP(w, request)
	result := w.Result()
	assert.Equal(t, http.StatusCreated, result.StatusCode)
	newUrl, err := ioutil.ReadAll(result.Body)
	assert.NoError(t, err)
	err = result.Body.Close()
	assert.NoError(t, err)
	assert.Contains(t, string(newUrl), "/0")

	// step 2: redirect to
	request = httptest.NewRequest(http.MethodGet, "/0", body)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, request)
	result = w.Result()
	assert.Equal(t, http.StatusTemporaryRedirect, result.StatusCode)
	assert.Equal(t, "https://yc.gl", result.Header.Get("Location"))
}
