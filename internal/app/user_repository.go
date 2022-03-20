package app

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/AtCliffUnderline/url-shortener/internal/app/models"
	"strconv"
)

var tokenUserLink = map[string]models.User{}
var userRoutesMap = map[int][]UserRoute{}

type UserRepository struct {
}
type UserRoute struct {
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

func (repo *UserRepository) AddRouteForUser(user models.User, route UserRoute) {
	userRoutesMap[user.Id] = append(userRoutesMap[user.Id], route)
}

func (repo *UserRepository) GetUserRoutes(user models.User) []UserRoute {
	return userRoutesMap[user.Id]
}

func (repo *UserRepository) CreateUser() models.User {
	id := len(tokenUserLink) + 1
	hash := sha256.Sum256([]byte(strconv.Itoa(id)))
	token := hex.EncodeToString(hash[:])
	user := models.User{
		Id:    id,
		Token: token,
	}
	tokenUserLink[token] = user

	return user
}

func (repo *UserRepository) GetUserByToken(token string) (models.User, error) {
	if user, err := tokenUserLink[token]; err {
		return user, nil
	}

	return models.User{}, errors.New("User with such token not found")
}
