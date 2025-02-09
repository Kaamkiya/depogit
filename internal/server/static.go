package server

import "net/http"

func serveStatic(w http.ResponseWriter, r *http.Request) {
	f := r.PathValue("file")

	http.ServeFile(w, r, f)
}
