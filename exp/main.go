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
	defer us.Close()
	us.DestructiveReset()
	user := models.User{
		Name:  "Micael Scott",
		Email: "michael@dundermilflen.com",
	}
	if err := us.Create(&user); err != nil {
		panic(err)
	}
	user.Email = "michaelscott@michaelscottpaper.com"
	if err := us.Update(&user); err != nil {
		panic(err)
	}
	userByID, err := us.ByID(user.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(userByID)
}
