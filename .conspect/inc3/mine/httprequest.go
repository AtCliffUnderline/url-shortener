package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"sort"
	"time"
)

type (
	User struct {
		ID       int     `json:"id"`
		Name     string  `json:"name"`
		Username string  `json:"username"`
		Email    string  `json:"email"`
		Address  Address `json:"address"`
		Phone    string  `json:"phone"`
		Website  string  `json:"website"`
		Company  Company `json:"company"`
	}
	Address struct {
		Street  string `json:"street"`
		Suite   string `json:"suite"`
		City    string `json:"city"`
		Zipcode string `json:"zipcode"`
	}

	Company struct {
		Name        string `json:"name"`
		CatchPhrase string `json:"catchPhrase"`
		Bs          string `json:"bs"`
	}
)

func main() {
	var users []User
	url := "https://jsonplaceholder.typicode.com/users"
	client := resty.New()

	client.
		SetRetryCount(5).
		SetRetryWaitTime(10 * time.Second)

	_, err := client.R().
		SetResult(&users).
		Get(url)
	if err != nil {
		panic(err)
		return
	}
	sort.SliceStable(users, func(i, j int) bool {
		return users[i].Name < users[j].Name
	})
	fmt.Println(users)
}
