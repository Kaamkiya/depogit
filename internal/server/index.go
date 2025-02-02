package server

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
)

func index(w http.ResponseWriter, r *http.Request) {
	log.Println(repoPath)
	allPaths, err := os.ReadDir(repoPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to read path: %v\n", err)
		return
	}

	var repoPaths []os.DirEntry
	for _, path := range allPaths {
		if path.IsDir() {
			repoPaths = append(repoPaths, path)
		}
	}

	type repo struct {
		Name string
		Time time.Time
	}
	var repos []repo

	for _, rp := range repoPaths {
		if rp.IsDir() {
			r, err := git.PlainOpen(repoPath + rp.Name())
			if err != nil {
				continue
			}

			head, err := r.Head()
			if err != nil {
				log.Printf("ERRO: %s: Failed to get repo head: %v\n", rp.Name(), err)
				continue
			}

			commit, err := r.CommitObject(head.Hash())
			if err != nil {
				log.Printf("Failed to get most recent repo commit: %v\n", err)
			}

			repos = append(repos, repo{
				Name: rp.Name(),
				Time: commit.Author.When,
			})
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/index.tmpl"))

	tmpl.Execute(w, repos)
}
