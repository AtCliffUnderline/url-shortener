package app

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func StartServer() {
	createRouting()
	http.ListenAndServe(":8080", nil)
}

func createRouting() {
	http.HandleFunc("/", UrlShortenerHandler)
}

func UrlShortenerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.URL.Path == "/" {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		id := ShortRoute(string(b))
		var newURL strings.Builder
		newURL.WriteString("http://localhost:8080/")
		newURL.WriteString(strconv.Itoa(id))
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(newURL.String()))

		return
	} else if r.Method == http.MethodGet {
		pathID := r.URL.Path
		id, err := strconv.Atoi(pathID[1:])
		if err != nil {
			http.Error(w, "Bad ID", http.StatusBadRequest)
		}
		route, err := GetRouteById(id)
		if err != nil {
			http.Error(w, "Bad ID", http.StatusBadRequest)
		}
		w.Header().Set("Location", route)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	http.Error(w, "Method not found", http.StatusBadRequest)
}
