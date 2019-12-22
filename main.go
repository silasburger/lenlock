package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func handleFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/" {
		fmt.Fprintf(w, "<h1>homepage Endpoint Hit!</h1>")
	} else if r.URL.Path == "/contact" {
		fmt.Fprintf(w, "for help contact the support team: <a href=\"mailto:support@lenslocked.com\">support@lenslocked.com</a>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "<h1>Oops, looks like you're lost. If we are giving you bad directions please let us know: <a href=\"mailto:support@lenslocked.com\">support@lenslocked.com</a></h1>")
	}
}

func handleRequests() {
	router := mux.NewRouter()
	router.HandleFunc("/", handleFunc)
	http.ListenAndServe(":3000", router)
}

func main() {
	handleRequests()
}
