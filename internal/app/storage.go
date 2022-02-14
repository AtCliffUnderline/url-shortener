package app

import "errors"

type RouteStorage interface {
	ShortRoute(fullRoute string) int
	GetRouteByID(id int) (string, error)
}

type DefaultRouteStorage struct {
}

var routeMap = map[int]string{}

func (DefaultRouteStorage) ShortRoute(fullRoute string) int {
	id := len(routeMap)
	routeMap[id] = fullRoute

	return id
}

func (DefaultRouteStorage) GetRouteByID(id int) (string, error) {
	if routeMap[id] == "" {
		return "", errors.New("no route with this ID found")
	}

	return routeMap[id], nil
}
