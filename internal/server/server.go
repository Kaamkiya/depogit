package server

import "net/http"

var repoPath = "/home/"

func Start(addr string, scanPath string) {
	repoPath = scanPath
	http.ListenAndServe(addr, route())
}

func route() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", index)
	mux.HandleFunc("GET /static/{file}", serveStatic)

	return mux
}
