package app

import "errors"

var routeMap = map[int]string{}

func ShortRoute(fullRoute string) int {
	id := len(routeMap)
	routeMap[id] = fullRoute

	return id
}

func GetRouteById(id int) (string, error) {
	if routeMap[id] == "" {
		return "", errors.New("No route with this ID found")
	}

	return routeMap[id], nil
}
