package server

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/go-git/go-git/v5"
)

func index(w http.ResponseWriter, r *http.Request) {
	allPaths, err := os.ReadDir(repoPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("ERRO: Failed to read path: %v\n", err)
		return
	}

	type repo struct {
		Name string
		Time string
	}

	var repos []repo

	for _, rp := range allPaths {
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
				log.Printf("ERRO: Failed to get most recent repo commit: %v\n", err)
			}

			repos = append(repos, repo{
				Name: rp.Name(),
				Time: humanize.Time(commit.Author.When),
			})
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/index.tmpl"))

	tmpl.Execute(w, repos)
}
