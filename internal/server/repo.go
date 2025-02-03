package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func repoLog(w http.ResponseWriter, r *http.Request) {
	repoName := r.PathValue("repo")

	var commits []*object.Commit

	repo, err := git.PlainOpen(repoPath + repoName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("ERRO: Failed to open repo %s: %v\n", repoName, err)
		return
	}

	head, err := repo.Head()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("ERRO: Failed to get repo %s HEAD: %v\n", repoName, err)
	}

	objs, err := repo.Log(&git.LogOptions{
		From:  head.Hash(),
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("ERRO: Failed to get commits from repo %s: %v\n", repoName, err)
		return
	}

	objs.ForEach(func(o *object.Commit) error {
		commits = append(commits, o)
		return nil
	})

	data := struct {
		Commits []*object.Commit
		Name    string
	}{
		Commits: commits,
		Name:    repoName,
	}

	tmpl := template.Must(template.ParseFiles("templates/repo_log.tmpl"))
	tmpl.Execute(w, data)
}

func repoIndex(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Name string
	}{
		Name: r.PathValue("repo"),
	}

	tmpl := template.Must(template.ParseFiles("templates/repo_index.tmpl"))
	tmpl.Execute(w, data)
}
