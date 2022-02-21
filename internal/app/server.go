package app

import (
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type HandlersCollection struct {
	Storage RouteStorage
}

func StartServer() {
	handlerCollection := &HandlersCollection{Storage: &DefaultRouteStorage{}}
	router := handlerCollection.CreateRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

func (h *HandlersCollection) CreateRouter() *chi.Mux {
	router := chi.NewRouter()
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method not found", http.StatusBadRequest)
	})
	router.Post("/", h.shortURLHandler)
	router.Get("/{id}", h.retrieveShortURLHandler)

	return router
}

func (h *HandlersCollection) shortURLHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	id := h.Storage.ShortRoute(string(b))
	var newURL strings.Builder
	newURL.WriteString("http://localhost:8080/")
	newURL.WriteString(strconv.Itoa(id))
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(newURL.String()))
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
}

func (h *HandlersCollection) retrieveShortURLHandler(w http.ResponseWriter, r *http.Request) {
	pathID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(pathID)
	if err != nil {
		http.Error(w, "Bad ID", http.StatusBadRequest)
	}
	route, err := h.Storage.GetRouteByID(id)
	if err != nil {
		http.Error(w, "Bad ID", http.StatusBadRequest)
	}
	w.Header().Set("Location", route)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
