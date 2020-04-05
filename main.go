package main

import (
	"fmt"
	"net/http"

	"use-go/lenslocked.com/controllers"
	"use-go/lenslocked.com/models"

	"github.com/gorilla/mux"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "silas.burger"
	dbname = "lenslocked"
)

func handleRequests(us *models.UserService) {
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)

	router := mux.NewRouter()
	// router.NotFoundHandler = http.HandlerFunc(notFound)
	router.Handle("/", staticC.Home).Methods("GET")
	router.Handle("/contact", staticC.Contact).Methods("GET")
	router.HandleFunc("/signup", usersC.New).Methods("GET")
	router.HandleFunc("/signup", usersC.Create).Methods("POST")
	http.ListenAndServe(":8080", router)
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)

	us, err := models.NewUserService(psqlInfo)
	us.DestructiveReset()
	must(err)
	defer us.Close()
	us.AutoMigrate()
	handleRequests(us)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
