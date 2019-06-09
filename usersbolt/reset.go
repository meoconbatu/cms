package usersbolt

import (
	"crypto/rand"
	"encoding/base64"
)

func genRandBytes() []byte {
	b := make([]byte, 24)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return []byte(base64.URLEncoding.EncodeToString(b))
}
