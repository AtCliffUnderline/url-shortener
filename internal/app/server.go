package app

import (
	"encoding/json"
	"errors"
	"github.com/AtCliffUnderline/url-shortener/internal/app/models"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ShortenerService struct {
	Config         ApplicationConfig
	Storage        RouteStorage
	UserRepository UserRepository
}

type URLShortenerRequest struct {
	URL string `json:"url"`
}

type URLShortenerResponse struct {
	URL string `json:"result"`
}

func StartServer(config ApplicationConfig) {
	service := &ShortenerService{
		Storage: &DefaultRouteStorage{},
		Config:  config,
	}
	if config.StoragePath != "" {
		service = &ShortenerService{
			Storage: &FileRouteStorage{
				FilePath: config.StoragePath,
			},
			Config: config,
		}
	}
	service.UserRepository = UserRepository{}
	router := service.CreateRouter()
	log.Fatal(http.ListenAndServe(service.Config.ServerAddress, router))
}

func (service *ShortenerService) CreateRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(gzipHandle)
	router.Use(service.authMiddleware)
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method not found", http.StatusBadRequest)
	})
	router.Post("/", service.shortURLHandler)
	router.Post("/api/shorten", service.alternativeShortURLHandler)
	router.Get("/{id}", service.retrieveShortURLHandler)
	router.Get("/api/user/urls", service.getUserURLs)

	return router
}

func (service *ShortenerService) alternativeShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var request URLShortenerRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if request.URL == "" {
		http.Error(w, "No Url provided", http.StatusBadRequest)
		return
	}
	id, err := service.Storage.ShortRoute(request.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var newURL strings.Builder
	newURL.WriteString(service.Config.BaseURL)
	newURL.WriteString("/")
	newURL.WriteString(strconv.Itoa(id))
	w.Header().Set("Content-Type", "application/json")
	err = service.saveRouteForUser(r, newURL.String(), request.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := URLShortenerResponse{URL: newURL.String()}
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
}

func (service *ShortenerService) shortURLHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	id, err := service.Storage.ShortRoute(string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var newURL strings.Builder
	newURL.WriteString(service.Config.BaseURL)
	newURL.WriteString("/")
	newURL.WriteString(strconv.Itoa(id))
	err = service.saveRouteForUser(r, newURL.String(), string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(newURL.String()))
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
}

func (service *ShortenerService) getUserURLs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, isOk := ctx.Value(UserContext{}).(models.User)
	if !isOk {
		http.Error(w, "unable to retrieve user", http.StatusInternalServerError)
	}
	routes := service.UserRepository.GetUserRoutes(user)
	if routes == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	response, err := json.Marshal(routes)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
}

func (service *ShortenerService) retrieveShortURLHandler(w http.ResponseWriter, r *http.Request) {
	pathID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(pathID)
	if err != nil {
		http.Error(w, "Bad ID", http.StatusBadRequest)
	}
	route, err := service.Storage.GetRouteByID(id)
	if err != nil {
		http.Error(w, "Bad ID", http.StatusBadRequest)
	}
	w.Header().Set("Location", route)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (service *ShortenerService) saveRouteForUser(r *http.Request, newRoute string, originalRoute string) error {
	ctx := r.Context()
	user, isOk := ctx.Value(UserContext{}).(models.User)
	if !isOk {
		return errors.New("unable to save route")
	}
	route := UserRoute{
		ShortURL:    newRoute,
		OriginalURL: originalRoute,
	}
	service.UserRepository.AddRouteForUser(user, route)

	return nil
}
