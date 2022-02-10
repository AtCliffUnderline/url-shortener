package app_test

import (
	"github.com/AtCliffUnderline/url-shortener/internal/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	router := app.CreateRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	// Testing if wrong method on existing path returns Bad Request
	resp, _ := testRequest(t, server, http.MethodPatch, "/", strings.NewReader(""))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Testing if existing method on wrong path returns 400
	resp, _ = testRequest(t, server, http.MethodPost, "/route-not-exists", strings.NewReader(""))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Testing symbolic IDs
	resp, _ = testRequest(t, server, http.MethodGet, "/id", strings.NewReader(""))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Testing unexciting ID
	resp, _ = testRequest(t, server, http.MethodGet, "/0", strings.NewReader(""))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Testing successful scenario
	resp, body := testRequest(t, server, http.MethodPost, "/", strings.NewReader("https://google.com"))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Contains(t, body, "/0")

	resp, body = testRequest(t, server, http.MethodGet, "/0", strings.NewReader(""))
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	assert.Equal(t, "https://google.com", resp.Header.Get("Location"))
}
