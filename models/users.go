package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	//ErrNotFound is returned when a resource is not found in the database
	ErrNotFound = errors.New("models: resource not found")
)

// NewUserService creates a connection to our db and returns
// a reference toUserService struct with the db connection
func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	return &UserService{
		db: db,
	}, nil
}

//UserService holds our db and methods to handle user table in db
type UserService struct {
	db *gorm.DB
}

//ByID will look up by the id provided
//1 - user, nil
//2 - nil, ErrNotFound
//3 - nil, otherError
func (us *UserService) ByID(id uint) (*User, error) {
	us.db.LogMode(true)
	var user User
	err := us.db.Where("id = ?", id).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Create will create the provided user
// and backfill data like ID, CreatedAt, and
// UpdatedAt fields
func (us *UserService) Create(user *User) error {
	us.db.LogMode(true)
	return us.db.Create(user).Error
}

// Update will update the provided user with all of the data
// in the provided user object.
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

// Close will close the UserService database connection
func (us *UserService) Close() error {
	return us.db.Close()
}

//DestructiveReset drops the user table and rebuilds it
func (us *UserService) DestructiveReset() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
}

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index`
}
