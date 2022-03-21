package app

import (
	"errors"
)

type RouteStorage interface {
	ShortRoute(fullRoute string) (int, error)
	GetRouteByID(id int) (string, error)
	SaveBatchRoutes(routes []BatchURLShortenerRequest) ([]BatchURLShortenerURLIDs, error)
}

type ShortenedURL struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type DefaultRouteStorage struct {
}

var routeMap = map[int]string{}

func (DefaultRouteStorage) ShortRoute(fullRoute string) (int, error) {
	id := len(routeMap)
	routeMap[id] = fullRoute

	return id, nil
}

func (DefaultRouteStorage) GetRouteByID(id int) (string, error) {
	if routeMap[id] == "" {
		return "", errors.New("no route with this ID found")
	}

	return routeMap[id], nil
}

func (st DefaultRouteStorage) SaveBatchRoutes(routes []BatchURLShortenerRequest) ([]BatchURLShortenerURLIDs, error) {
	var result []BatchURLShortenerURLIDs
	for _, URLToShort := range routes {
		id, _ := st.ShortRoute(URLToShort.URL)
		result = append(result, BatchURLShortenerURLIDs{ID: id, CorrelationID: URLToShort.ID, OriginalURL: URLToShort.URL})
	}

	return result, nil
}
