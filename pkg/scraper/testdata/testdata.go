package testdata

import (
	"embed"
	"net/http"
)

//go:embed *.html
var files embed.FS

func FileHandler() http.Handler {

	return http.FileServer(http.FS(files))

	// server := &http.Server{
	// 	Addr: "0.0.0.0:9977",
	// }

	// server.HandleFunc("/", fs.ServeHTTP)

	// return server

	// return fs
}
