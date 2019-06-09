package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

// CreateLogger creates a new loggerlogger that writeswrites to the given filename
func CreateLogger(filename string) *log.Logger {
	file, err := os.OpenFile(filename+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	logger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	return logger
}

// Time runs the next function in the chain
func Time(logger *log.Logger, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)
		logger.Println(elapsed)
	})
}

// PassContext is use to pass values between middlewares
type PassContext func(ctx context.Context, w http.ResponseWriter, r *http.Request)

func (fn PassContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(context.Background(), "foo", "bar")
	fn(ctx, w, r)
}
