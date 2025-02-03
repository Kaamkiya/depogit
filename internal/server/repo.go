package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/dustin/go-humanize"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func repoLog(w http.ResponseWriter, r *http.Request) {
	repoName := r.PathValue("repo")

	type commit struct {
		Message      string
		Time         string
		LinesAdded   int
		LinesRemoved int
	}
	var commits []commit

	repo, err := git.PlainOpen(repoPath + repoName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("ERRO: Failed to open repo %s: %v\n", repoName, err)
		return
	}

	objs, err := repo.CommitObjects()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("ERRO: Failed to get commits from repo %s: %v\n", repoName, err)
		return
	}

	objs.ForEach(func(o *object.Commit) error {
		stats, err := o.Stats()
		if err != nil {
			return err
		}

		linesAdded := 0
		linesRemoved := 0
		for _, s := range stats {
			linesAdded += s.Addition
			linesRemoved += s.Deletion
		}

		commits = append(commits, commit{
			Message:      o.Message,
			Time:         humanize.Time(o.Author.When),
			LinesAdded:   linesAdded,
			LinesRemoved: linesRemoved,
		})
		return nil
	})

	data := struct {
		Commits []commit
		Name    string
	}{
		Commits: commits,
		Name:    repoName,
	}

	tmpl := template.Must(template.ParseFiles("templates/repo_log.tmpl"))
	tmpl.Execute(w, data)
}
