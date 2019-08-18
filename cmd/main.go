package main

import (
	"net/http"
	"os"

	"golang.org/x/net/context"

	"github.com/meoconbatu/cms"
	"github.com/meoconbatu/cms/middleware"
)

func withContext(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	bar := ctx.Value("foo")
	w.Write([]byte(bar.(string)))
}
func main() {
	port := os.Getenv("PORT")
	logger := middleware.CreateLogger("server")
	http.Handle("/", middleware.Time(logger, cms.ServeIndex))
	http.HandleFunc("/new", cms.HandleNew)
	http.HandleFunc("/page/", cms.ServePage)
	http.HandleFunc("/post", cms.ServePost)
	http.HandleFunc("/login", cms.ServeLogin)
	http.HandleFunc("/restrict", cms.ServeRestricted)

	http.Handle("/context", middleware.PassContext(withContext))

	http.HandleFunc("/auth/gplus/authorize", cms.AuthURLHandler)
	http.HandleFunc("/auth/gplus/callback", cms.CallbackURLHandler)
	http.HandleFunc("/oauth", cms.ServeOAuthRestricted)
	if port == "" {
		port = "3000"
	}
	http.ListenAndServe(":"+port, nil)
}
