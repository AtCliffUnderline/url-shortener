package app

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"os"
)

type FileRouteStorage struct {
	FilePath string
}

func (storage FileRouteStorage) DeleteRouteByIDForUser(_ int, _ int) error {
	return nil
}

func (storage FileRouteStorage) ShortRoute(fullRoute string, _ int) (int, error) {
	file, err := os.OpenFile(storage.FilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return 0, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
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

func (storage FileRouteStorage) SaveBatchRoutes(routes []BatchURLShortenerRequest, userID int) ([]BatchURLShortenerURLIDs, error) {
	var result []BatchURLShortenerURLIDs
	for _, URLToShort := range routes {
		id, _ := storage.ShortRoute(URLToShort.URL, userID)
		result = append(result, BatchURLShortenerURLIDs{ID: id, CorrelationID: URLToShort.ID, OriginalURL: URLToShort.URL})
	}

	return result, nil
}
