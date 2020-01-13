package main

import (
	"fmt"
	"net/http"

	"use-go/lenslocked.com/controllers"

	"github.com/gorilla/mux"
	"use-go/lenslocked.com/views"
)

var (
	homeView    *views.View
	contactView *views.View
	faqView     *views.View
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w, nil))
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w, nil))
}

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(faqView.Render(w, nil))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "<h1>Oops, looks like you're lost. If we are giving you bad directions please let us know: <a href=\"mailto:support@lenslocked.com\">support@lenslocked.com</a></h1>")
}

func handleRequests() {
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")
	faqView = views.NewView("bootstrap", "views/faq.gohtml")
	usersC := controllers.NewUsers()

	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.HandleFunc("/", home)
	router.HandleFunc("/contact", contact)
	router.HandleFunc("/faq", faq)
	router.HandleFunc("/signup", usersC.New)
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
