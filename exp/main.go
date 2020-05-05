package main

import (
	"fmt"

	_ "github.com/lib/pq"
	"use-go/lenslocked.com/models"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "silas.burger"
	dbname = "lenslocked"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)

	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	us.DestructiveReset()
	us.SetLogging(true)
	defer us.Close()

	user := models.User{
		Name:     "silas burger",
		Email:    "silas@burger.io",
		Password: "jon",
		Remember: "abc123",
	}

	err = us.Create(&user)
	if err != nil {
		panic(err)
	}
	user2, err := us.ByRemember("abc123")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", *user2)
}
