package app

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
)

type RouteStorage interface {
	ShortRoute(fullRoute string) (int, error)
	GetRouteByID(id int) (string, error)
}

type ShortenedURL struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type DefaultRouteStorage struct {
}

type FileRouteStorage struct {
	FilePath string
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

func (storage FileRouteStorage) ShortRoute(fullRoute string) (int, error) {
	file, err := os.OpenFile(storage.FilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return 0, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	scanner := bufio.NewScanner(file)
	lastShortenedURL := ShortenedURL{}
	for {
		scanner.Scan()
		bytes := scanner.Bytes()
		if bytes == nil {
			break
		}
		err := json.Unmarshal(bytes, &lastShortenedURL)
		if err != nil {
			return 0, err
		}
	}
	id := lastShortenedURL.ID + 1
	newRoute := ShortenedURL{
		URL: fullRoute,
		ID:  id,
	}
	data, err := json.Marshal(&newRoute)
	if err != nil {
		return 0, err
	}
	data = append(data, '\n')

	_, err = file.Write(data)

	return id, err
}

func (storage FileRouteStorage) GetRouteByID(id int) (string, error) {
	file, err := os.OpenFile(storage.FilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	scanner := bufio.NewScanner(file)
	shortenedURL := ShortenedURL{ID: 0, URL: ""}
	for {
		scanner.Scan()
		bytes := scanner.Bytes()
		if bytes == nil {
			break
		}
		err := json.Unmarshal(bytes, &shortenedURL)
		if err != nil {
			return "", err
		}
		if shortenedURL.ID == id {
			return shortenedURL.URL, nil
		}
	}

	return "", errors.New("no route with this ID found")
}
