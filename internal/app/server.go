package app

import (
	"encoding/json"
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

type UrlShortenerRequest struct {
	URL string `json:"url"`
}

type UrlShortenerResponse struct {
	URL string `json:"result"`
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
	router.Post("/api/shorten", h.alternativeShortUrlHandler)
	router.Get("/{id}", h.retrieveShortURLHandler)

	return router
}

func (h *HandlersCollection) alternativeShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	var request UrlShortenerRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if request.URL == "" {
		http.Error(w, "No Url provided", http.StatusBadRequest)
		return
	}
	id := h.Storage.ShortRoute(request.URL)
	var newURL strings.Builder
	newURL.WriteString("http://localhost:8080/")
	newURL.WriteString(strconv.Itoa(id))
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	result := UrlShortenerResponse{URL: newURL.String()}
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
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
