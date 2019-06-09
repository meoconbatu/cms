package users

import (
	"errors"
	"sync"

	scrypt "github.com/elithrar/simple-scrypt"
)

var (
	// DB is the reference to our DB, which contains our user data.
	DB = newDB()

	// ErrUserAlreadyExists is the error thrown when a user attempts to create
	// a new user in the DB with a duplicate username.
	ErrUserAlreadyExists = errors.New("users: username already exists")
)

// Store is a very simple in memory database, that we'll use to store our users.
// It is protected by read-wrote mutexutex, so that two goroutines can't modify
// the underlying  map at the same time (since maps are not safe for concurrent use in GoGo)
type Store struct {
	rwm *sync.RWMutex
	m   map[string]string
}

func newDB() *Store {
	return &Store{
		rwm: &sync.RWMutex{},
		m:   make(map[string]string),
	}
}

// NewUser accepts a username and password and create a new user in our DB
func NewUser(username, password string) error {
	err := exists(username)
	if err != nil {
		return err
	}

	DB.rwm.Lock()
	defer DB.rwm.Unlock()
	hashedPassword, err := scrypt.GenerateFromPassword([]byte(password), scrypt.DefaultParams)
	if err != nil {
		return nil
	}
	DB.m[username] = string(hashedPassword)
	return nil
}

// AuthenticateUser accepts a username and password, and check that the given password
// matches the hashed password. It returns nil on success, and an error on failure.
func AuthenticateUser(username, password string) error {
	DB.rwm.RLock()
	defer DB.rwm.RUnlock()

	hashedPassword := DB.m[username]
	err := scrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}

// OverrideOldPassword overrides the old password with the new password.
// For use when resetting password
func OverrideOldPassword(username, password string) error {
	DB.rwm.Lock()
	defer DB.rwm.Unlock()
	hashedPassword, err := scrypt.GenerateFromPassword([]byte(password), scrypt.DefaultParams)
	if err != nil {
		return nil
	}
	DB.m[username] = string(hashedPassword)
	return nil
}
func exists(username string) error {
	DB.rwm.Lock()
	defer DB.rwm.Unlock()
	if DB.m[username] != "" {
		return ErrUserAlreadyExists
	}
	return nil
}
