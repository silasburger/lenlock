package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"use-go/lenslocked.com/views"
)

var (
	homeView    *views.View
	contactView *views.View
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w, nil))
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w, nil))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "<h1>Oops, looks like you're lost. If we are giving you bad directions please let us know: <a href=\"mailto:support@lenslocked.com\">support@lenslocked.com</a></h1>")
}

func handleRequests() {
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")

	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.HandleFunc("/", home)
	router.HandleFunc("/contact", contact)
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
