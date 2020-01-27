package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jinzhu/gorm"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "silas.burger"
	dbname = "lenslocked"
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
	Color string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)

	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.DB().Ping()
	if err != nil {
		panic(err)
	}

	db.LogMode(true)
	db.AutoMigrate(&User{})
	// db.DropTableIfExists(&User{})
	name, email, color := getInfo()
	u := User{
		Name:  name,
		Email: email,
		Color: color,
	}
	if err = db.Create(&u).Error; err != nil {
		panic(err)
	}
	fmt.Println("successfully connected to the db")
}
func getInfo() (name, email, color string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("What is your name?")
	name, _ = reader.ReadString('\n')
	fmt.Println("What is your email?")
	email, _ = reader.ReadString('\n')
	fmt.Println("What is your color?")
	color, _ = reader.ReadString('\n')
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	color = strings.TrimSpace(color)
	return name, email, color
}
