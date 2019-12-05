package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/" {
		fmt.Fprintf(w, "<h1>homepage Endpoint Hit!</h1>")
	} else if r.URL.Path == "/contact" {
		fmt.Fprintf(w, "for help contact the support team: <a href=\"mailto:support@lenslocked.com\">support@lenslocked.com</a>")
	}
}

func handleRequests() {
	http.HandleFunc("/", handleFunc)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	handleRequests()
}
