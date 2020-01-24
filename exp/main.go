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

	rows, err := db.Query(`
		SELECT users.id, users.email, users.name,
		orders.id AS order_id,
		orders.amount AS order_amount,
		orders.description AS order_description
		FROM users
		INNER JOIN orders
		ON users.id = orders.user_id`)

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var userID, orderID, amount int
		var email, name, desc string
		if err := rows.Scan(&userID, &name, &email, &orderID, &amount, &desc); err != nil {
			panic(err)
		}
		fmt.Println(userID, name, email, orderID, amount, desc)
	}
	if err != nil {
		panic(err)
	}
}
