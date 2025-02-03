package server

import "net/http"

func serveStatic(w http.ResponseWriter, r *http.Request) {
	f := r.PathValue("file")

	w.Header().Add("Content-Type", "text/css")

	http.ServeFile(w, r, f)
}
