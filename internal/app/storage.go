package app

import (
	"errors"
)

var RouteDeletedError = errors.New("route has been deleted")

const DeletedRoute = "deleted"

type RouteStorage interface {
	ShortRoute(fullRoute string, userID int) (int, error)
	GetRouteByID(id int) (string, error)
	SaveBatchRoutes(routes []BatchURLShortenerRequest, userID int) ([]BatchURLShortenerURLIDs, error)
	DeleteRouteByIDForUser(routeID int, userID int) error
}

type ShortenedURL struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type DefaultRouteStorage struct {
}

var routeMap = map[int]string{}

func (DefaultRouteStorage) DeleteRouteByIDForUser(routeID int, _ int) error {
	routeMap[routeID] = DeletedRoute

	return nil
}

func (DefaultRouteStorage) ShortRoute(fullRoute string, _ int) (int, error) {
	id := len(routeMap)
	routeMap[id] = fullRoute

	return id, nil
}

func (DefaultRouteStorage) GetRouteByID(id int) (string, error) {
	if routeMap[id] == "" {
		return "", errors.New("no route with this ID found")
	}

	if routeMap[id] == DeletedRoute {
		return "", RouteDeletedError
	}

	return routeMap[id], nil
}

func (st DefaultRouteStorage) SaveBatchRoutes(routes []BatchURLShortenerRequest, userID int) ([]BatchURLShortenerURLIDs, error) {
	var result []BatchURLShortenerURLIDs
	for _, URLToShort := range routes {
		id, _ := st.ShortRoute(URLToShort.URL, userID)
		result = append(result, BatchURLShortenerURLIDs{ID: id, CorrelationID: URLToShort.ID, OriginalURL: URLToShort.URL})
	}

	return result, nil
}
