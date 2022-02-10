package app

import (
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func StartServer() {
	router := CreateRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

func CreateRouter() *chi.Mux {
	router := chi.NewRouter()
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method not found", http.StatusBadRequest)
	})
	router.Post("/", shortUrlHandler)
	router.Get("/{id}", retrieveShortUrlHandler)

	return router
}

func shortUrlHandler(w http.ResponseWriter, r *http.Request) {
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
}

func retrieveShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	pathID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(pathID)
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
