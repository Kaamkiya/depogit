package server

import "net/http"

var repoPath = "/var/www/git"

func Start(addr string, scanPath string) {
	repoPath = scanPath
	http.ListenAndServe(addr, route())
}

func route() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", index)
	mux.HandleFunc("GET /static/{file}", serveStatic)
	mux.HandleFunc("GET /{repo}", repoIndex)
	mux.HandleFunc("GET /{repo}/log/{$}", repoLog)

	return mux
}
