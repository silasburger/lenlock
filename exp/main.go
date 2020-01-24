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

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("successfully connected to the db")

	defer db.Close()

	type User struct {
		Name  string
		Email string
		ID    int
	}

	var users []User
	rows, err := db.Query(`
		SELECT id, name, email
		FROM users`)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("no rows")
			}
			panic(err)
		}
		users = append(users, user)
	}

	fmt.Println(users)
}
