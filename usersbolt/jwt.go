package usersbolt

import (
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	identityURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	provider    = New()
	signingKey  = genRandBytes()
)

// New creates a new oauth2 config
func New() *oauth2.Config {
	return &oauth2.Config{
		ClientID: "703558633182-6c57la34kieji5i97dsc24a873akf4o4.apps.googleusercontent.com",
		//os.Getenv("GOOGLE_KEY"),
		ClientSecret: "aSA4tV8EiGXwUSqlv2UNKY0A",
		//os.Getenv("GOOGLE_SECRET"),
		Endpoint:    google.Endpoint,
		RedirectURL: "http://localhost:3000/auth/gplus/callback",
		Scopes:      []string{"email", "profile"},
	}
}

// AuthcodeURL return a URL that asks for permission
func AuthcodeURL() string {
	return provider.AuthCodeURL("", oauth2.AccessTypeOffline)
}

// GetToken return a URL that asks for permission
func GetToken(code string) ([]byte, error) {
	token, err := provider.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}
	client := provider.Client(oauth2.NoContext, token)
	resp, err := client.Get(identityURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	user := make(map[string]string)
	json.NewDecoder(resp.Body).Decode(&user)

	email := user["email"]
	return genToken(email)
}
func genToken(email string) ([]byte, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["sub"] = email
	token.Claims["exp"] = time.Now().Add(time.Hour * 72)
	token.Claims["iat"] = time.Now().Unix()

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return nil, err
	}
	return []byte(tokenString), nil
}

// VerifyToken gets the token from an HTTP request, and ensures that it's valid.
// It'll return the user's username as a string
func VerifyToken(r *http.Request) (string, error) {
	token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return signingKey, nil
	})
	if err != nil {
		return "", err
	}
	if token.Valid == false {
		return "", jwt.ErrInvalidKey
	}
	return token.Claims["sub"].(string), nil
}
