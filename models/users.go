package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"use-go/lenslocked.com/hash"
	"use-go/lenslocked.com/rand"
)

const userPwPepper = "secret-random-string"
const hmacSecretKey = "secret-hmac-ket"

var (
	// ErrNotFound is returned when a resource is not
	// found in the database
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID is
	// provided to a method like delete
	ErrInvalidID = errors.New("models: ID provided was invalid")

	// ErrInvalidPassword is return when a user cannot be
	// authenticated due to mismatched passwords
	ErrInvalidPassword = errors.New("models: incorrect password provided")
)

// NewUserService creates a connection to our db and returns
// a reference toUserService struct with the db connection
func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	hmac := hash.NewHMAC(hmacSecretKey)
	if err != nil {
		return nil, err
	}
	return &UserService{
		db:   db,
		hmac: hmac,
	}, nil
}

//UserService holds our db and methods to handle user table in db
type UserService struct {
	db   *gorm.DB
	hmac hash.HMAC
}

//ByID will look up by the id provided
//1 - user, nil
//2 - nil, ErrNotFound
//3 - nil, otherError
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id).First(&user)
	err := first(db, &user)
	return &user, err
}

//ByEmail will look up a user by their email and return a user
//1 - user, nil
//2 - nil, ErrNotFound
//3 - nil, otherError
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// ByRemember will look up a user by their remember token
// and returns that user.This method will handle hashing
// the token for us.
// Errors are the same as ByEmail
func (us *UserService) ByRemember(token string) (*User, error) {
	var user User
	hashedToken := us.hmac.Hash(token)
	err := first(us.db.Where("remember_hash = ?", hashedToken), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Authenticate will authenticate a user with the
// provided email and password
func (us *UserService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}
	return foundUser, nil
}

//InAgeRange will look up users in a specific age range return a slice of users
//1 - user, nil
//2 - nil, ErrNotFound
//3 - nil, otherError
func (us *UserService) InAgeRange(min, max uint) (*[]User, error) {
	var users []User
	err := us.db.Where("age >= ? AND age <= ?", min, max).Find(&users).Error
	return &users, err
}

//ByAge will look up a user by their age and return a user
//1 - user, nil
//2 - nil, ErrNotFound
//3 - nil, otherError
func (us *UserService) ByAge(age uint) (*User, error) {
	var user User
	db := us.db.Where("age = ?", age)
	err := first(db, &user)
	return &user, err
}

// first will query using the provided gorm.db and it will get
// the first item returned and place it into dst. If
// nothing is found in the query, it will return ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// Create will create the provided user
// and backfill data like ID, CreatedAt, and
// UpdatedAt fields
func (us *UserService) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = us.hmac.Hash(user.Remember)
	return us.db.Create(user).Error
}

// Update will update the provided user with all of the data
// in the provided user object.
func (us *UserService) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = us.hmac.Hash(user.Remember)
	}
	return us.db.Save(user).Error
}

// Delete will delete the user with the provided ID
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

// Close will close the UserService database connection
func (us *UserService) Close() error {
	return us.db.Close()
}

//DestructiveReset drops the user table and rebuilds it
func (us *UserService) DestructiveReset() error {
	if err := us.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return us.AutoMigrate()
}

//SetLogging allows me to turn logging on or off
func (us *UserService) SetLogging(isLogging bool) {
	us.db.LogMode(isLogging)
}

// AutoMigrate will attemp to automatically migrate the users table
func (us *UserService) AutoMigrate() error {
	return us.db.AutoMigrate(&User{}).Error
}

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Age          uint
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}
