package main

import (
	"fmt"

	_ "github.com/lib/pq"
	"use-go/lenslocked.com/hash"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "silas.burger"
	dbname = "lenslocked"
)

func main() {
	hmac := hash.NewHMAC("my-secret-key")
	// This should print out:
	//   4waUFc1cnuxoM2oUOJfpGZLGP1asj35y7teuweSFgPY=
	fmt.Println(hmac.Hash("this is my string to hash"))
}
