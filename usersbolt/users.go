package usersbolt

import (
	"errors"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	scrypt "github.com/elithrar/simple-scrypt"
)

var (
	// DB is the reference to our DB, which contains our user data.
	DB = newDB()

	// ErrUserAlreadyExists is the error thrown when a user attempts to create
	// a new user in the DB with a duplicate username.
	ErrUserAlreadyExists = errors.New("users: username already exists")

	// ErrUserNotFound is the error thrown when a user can't be found in the
	// database.
	ErrUserNotFound = errors.New("users: user not found")
)

// Store is a very simple in memory database, that we'll use to store our users.
// It is protected by read-wrote mutexutex, so that two goroutines can't modify
// the underlying  map at the same time (since maps are not safe for concurrent use in GoGo)
type Store struct {
	DB              *bolt.DB
	Users, Sessions string
}

func newDB() *Store {
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil
	}
	s := &Store{DB: db, Users: "users", Sessions: "sessions"}
	err = s.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(s.Users))
		if err != nil {
			return fmt.Errorf("create users: %s", err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte(s.Sessions))
		if err != nil {
			return fmt.Errorf("create sessions: %s", err)
		}
		return nil
	})
	return s
}

// NewUser accepts a username and password and create a new user in our DB
func NewUser(username, password string) error {
	err := exists(username)
	if err != nil {
		return err
	}
	hashedPassword, err := scrypt.GenerateFromPassword([]byte(password), scrypt.DefaultParams)
	if err != nil {
		return nil
	}
	return DB.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB.Users))
		return b.Put([]byte(username), hashedPassword)
	})
}

// AuthenticateUser accepts a username and password, and check that the given password
// matches the hashed password. It returns nil on success, and an error on failure.
func AuthenticateUser(username, password string) error {
	var hashedPassword []byte
	DB.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB.Users))
		hashedPassword = b.Get([]byte(username))
		return nil
	})
	return scrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}

// OverrideOldPassword overrides the old password with the new password.
// For use when resetting password
func OverrideOldPassword(username, password string) error {
	hashedPassword, err := scrypt.GenerateFromPassword([]byte(password), scrypt.DefaultParams)
	if err != nil {
		return nil
	}
	return DB.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB.Users))
		return b.Put([]byte(username), hashedPassword)
	})
}
func exists(username string) error {
	var hashedPassword []byte
	DB.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB.Users))
		hashedPassword = b.Get([]byte(username))
		return nil
	})
	if hashedPassword != nil {
		return ErrUserAlreadyExists
	}
	return nil
}
