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
	Config  ApplicationConfig
	Storage RouteStorage
}

type URLShortenerRequest struct {
	URL string `json:"url"`
}

type URLShortenerResponse struct {
	URL string `json:"result"`
}

func StartServer() {
	config := getConfig()
	handlerCollection := &HandlersCollection{
		Storage: &DefaultRouteStorage{},
		Config:  config,
	}
	if config.StoragePath != "" {
		handlerCollection = &HandlersCollection{
			Storage: &FileRouteStorage{
				FilePath: config.StoragePath,
			},
			Config: getConfig(),
		}
	}
	router := handlerCollection.CreateRouter()
	log.Fatal(http.ListenAndServe(handlerCollection.Config.ServerAddress, router))
}

func (h *HandlersCollection) CreateRouter() *chi.Mux {
	router := chi.NewRouter()
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method not found", http.StatusBadRequest)
	})
	router.Post("/", h.shortURLHandler)
	router.Post("/api/shorten", h.alternativeShortURLHandler)
	router.Get("/{id}", h.retrieveShortURLHandler)

	return router
}

func (h *HandlersCollection) alternativeShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var request URLShortenerRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if request.URL == "" {
		http.Error(w, "No Url provided", http.StatusBadRequest)
		return
	}
	id, err := h.Storage.ShortRoute(request.URL)
	var newURL strings.Builder
	newURL.WriteString(h.Config.BaseURL)
	newURL.WriteString("/")
	newURL.WriteString(strconv.Itoa(id))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	result := URLShortenerResponse{URL: newURL.String()}
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
	id, err := h.Storage.ShortRoute(string(b))
	var newURL strings.Builder
	newURL.WriteString(h.Config.BaseURL)
	newURL.WriteString("/")
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
