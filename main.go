package main

import (
	"net/http"

	"use-go/lenslocked.com/controllers"

	"github.com/gorilla/mux"
)

func handleRequests() {
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers()

	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.Handle("/", staticC.HomeView).Methods("GET")
	router.Handle("/contact", staticC.ContactView).Methods("GET")
	router.HandleFunc("/signup", usersC.New).Methods("GET")
	router.HandleFunc("/signup", usersC.Create).Methods("POST")
	http.ListenAndServe(":3000", router)
}

func main() {
	handleRequests()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
