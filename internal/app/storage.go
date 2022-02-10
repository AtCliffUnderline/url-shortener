package app

import "errors"

var routeMap = map[int]string{}

func ShortRoute(fullRoute string) int {
	id := len(routeMap)
	routeMap[id] = fullRoute

	return id
}

func GetRouteByID(id int) (string, error) {
	if routeMap[id] == "" {
		return "", errors.New("no route with this ID found")
	}

	return routeMap[id], nil
}
