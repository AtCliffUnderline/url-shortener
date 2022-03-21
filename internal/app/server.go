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
	Database       BaseDB
}

type URLShortenerRequest struct {
	URL string `json:"url"`
}

type BatchURLShortenerRequest struct {
	ID  string `json:"correlation_id"`
	URL string `json:"original_url"`
}

type BatchURLShortenerResponse struct {
	ID  string `json:"correlation_id"`
	URL string `json:"short_url"`
}

type BatchURLShortenerURLIDs struct {
	ID            int
	OriginalURL   string
	CorrelationID string
}

type URLShortenerResponse struct {
	URL string `json:"result"`
}

func StartServer(config ApplicationConfig, db BaseDB) {
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
	if db.IsConnectionEstablished() {
		service = &ShortenerService{
			Storage: &DatabaseRouteStorage{
				baseDB: &db,
			},
			Config: config,
		}
	}
	service.UserRepository = UserRepository{}
	service.Database = db
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
	router.Post("/api/shorten/batch", service.batchShortURLHandler)
	router.Get("/{id}", service.retrieveShortURLHandler)
	router.Get("/api/user/urls", service.getUserURLs)
	router.Get("/ping", service.pingDatabase)

	return router
}

func (service *ShortenerService) batchShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var shortenedRoutes []BatchURLShortenerResponse
	var routes []BatchURLShortenerRequest
	if err := json.NewDecoder(r.Body).Decode(&routes); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	savedRoutes, err := service.Storage.SaveBatchRoutes(routes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, route := range savedRoutes {
		var newURL strings.Builder
		newURL.WriteString(service.Config.BaseURL)
		newURL.WriteString("/")
		newURL.WriteString(strconv.Itoa(route.ID))
		shortenedRoutes = append(shortenedRoutes, BatchURLShortenerResponse{URL: newURL.String(), ID: route.CorrelationID})
		err = service.saveRouteForUser(r, newURL.String(), route.OriginalURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(shortenedRoutes)
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

func (service *ShortenerService) alternativeShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var request URLShortenerRequest
	var httpStatus int

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if request.URL == "" {
		http.Error(w, "No Url provided", http.StatusBadRequest)
		return
	}
	id, err := service.Storage.ShortRoute(request.URL)
	if errors.As(err, &ErrRouteAlreadyShortened) {
		httpStatus = http.StatusConflict
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		httpStatus = http.StatusCreated
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
	w.WriteHeader(httpStatus)
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
}

func (service *ShortenerService) shortURLHandler(w http.ResponseWriter, r *http.Request) {
	var httpStatus int

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	id, err := service.Storage.ShortRoute(string(b))
	if errors.As(err, &ErrRouteAlreadyShortened) {
		httpStatus = http.StatusConflict
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		httpStatus = http.StatusCreated
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
	w.WriteHeader(httpStatus)
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

func (service *ShortenerService) pingDatabase(w http.ResponseWriter, _ *http.Request) {
	err := service.Database.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
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
