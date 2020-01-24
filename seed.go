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
	defer db.Close()

	if err != nil {
		panic(err)
	}
	createNewTables(db)
	addUsers(db)
	addOrders(db)

	fmt.Println("successfully seeded db")
}

func createNewTables(db *sql.DB) {
	_, err := db.Exec(`DROP TABLE users`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`DROP TABLE orders`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			name TEXT,
			email TEXT NOT NULL
		)`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`
		CREATE TABLE orders (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			amount INT,
			description TEXT
		)`)
	if err != nil {
		panic(err)
	}

}

func addUsers(db *sql.DB) {
	type User struct {
		Name  string
		Email string
	}

	users := []User{
		User{Name: "silas", Email: "silas@burger.com"},
		User{Name: "john", Email: "john@burger.com"},
		User{Name: "robert", Email: "robert@burger.com"},
		User{Name: "jon", Email: "jon@calhoun.com"},
	}

	for _, user := range users { //could use concurrency to speed this up
		_, err := db.Exec(`
			INSERT INTO users(name, email)
			VALUES ($1, $2)
			RETURNING id`,
			user.Name, user.Email)
		if err != nil {
			panic(err)
		}
	}
}

func addOrders(db *sql.DB) {
	for i := 1; i < 6; i++ {
		userID := 1
		if i > 3 {
			userID = 3
		}
		amount := i * 100
		description := fmt.Sprintf("USB-c Adapter x%d", amount)
		_, err := db.Exec(`
				INSERT INTO orders(user_id, amount, description)
				VALUES($1, $2, $3)`, userID, amount, description)
		if err != nil {
			panic(err)
		}
	}
}
