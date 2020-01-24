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

	addUsers(db)

	fmt.Println("successfully seeded db")
	db.Close()
}

func addUsers(db *sql.DB) {
	type User struct {
		name  string
		email string
	}

	users := []User{
		User{name: "silas", email: "silas@burger.com"},
		User{name: "john", email: "john@burger.com"},
		User{name: "robert", email: "robert@burger.com"},
		User{name: "jon", email: "jon@calhoun.com"},
	}

	for _, user := range users { //could use concurrency to speed this up
		db.QueryRow(`
			INSERT INTO users(name, email)
			VALUES ($1, $2)
			RETURNING id`,
			user.name, user.email)
	}

}
