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

	// go vet не дает использовать здесь defer, якобы все body остаются не закрытыми
	/**
	internal/app/server_test.go:41:24: response body must be closed
	internal/app/server_test.go:45:23: response body must be closed
	internal/app/server_test.go:49:23: response body must be closed
	internal/app/server_test.go:53:23: response body must be closed
	internal/app/server_test.go:57:27: response body must be closed
	internal/app/server_test.go:61:26: response body must be closed
	*/
	// defer resp.Body.Close()

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	router := app.CreateRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	// Testing if wrong method on existing path returns Bad Request
	resp, _ := testRequest(t, server, http.MethodPatch, "/", strings.NewReader(""))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()

	// Testing if existing method on wrong path returns 400
	resp, _ = testRequest(t, server, http.MethodPost, "/route-not-exists", strings.NewReader(""))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()

	// Testing symbolic IDs
	resp, _ = testRequest(t, server, http.MethodGet, "/id", strings.NewReader(""))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()

	// Testing unexciting ID
	resp, _ = testRequest(t, server, http.MethodGet, "/0", strings.NewReader(""))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()

	// Testing successful scenario
	resp, body := testRequest(t, server, http.MethodPost, "/", strings.NewReader("https://google.com"))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Contains(t, body, "/0")
	resp.Body.Close()

	resp, _ = testRequest(t, server, http.MethodGet, "/0", strings.NewReader(""))
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	assert.Equal(t, "https://google.com", resp.Header.Get("Location"))
	resp.Body.Close()
}
