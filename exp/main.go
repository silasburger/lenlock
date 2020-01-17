package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
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

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	var id int
	row := db.QueryRow(`
		INSERT INTO users(name, email)
		VALUES ($1, $2)
		RETURNING id`,
		"Jon Calhoun", "jon@calhoun.io")
	err = row.Scan(&id)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("successfully connected to the db")
	db.Close()
}
