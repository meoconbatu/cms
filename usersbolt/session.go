package usersbolt

import (
	"net/http"
	"time"

	"github.com/boltdb/bolt"
)

const cookieName = "_goproj_sess"

// GetSession gets the current session from the cookie
func GetSession(w http.ResponseWriter, r *http.Request) string {
	s, err := r.Cookie(cookieName)
	if err != nil {
		http.Error(w, "Please login to view this page", http.StatusUnauthorized)
		return ""
	}
	user, err := get(s.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}
	return user
}

// SetSession sets the session for the given user
func SetSession(w http.ResponseWriter, user string) {
	bytes := string(genRandBytes())
	err := save(string(bytes), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    bytes,
		Expires:  time.Now().Add(time.Hour * 72),
		HttpOnly: true,
	})
}
func save(id, user string) error {
	return DB.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB.Sessions))
		return b.Put([]byte(id), []byte(user))
	})
}

func get(id string) (string, error) {
	var user []byte
	DB.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB.Sessions))
		user = b.Get([]byte(id))
		return nil
	})
	if user == nil {
		return "", ErrUserNotFound
	}
	return string(user), nil
}
