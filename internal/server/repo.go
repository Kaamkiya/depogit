package server

import (
	"log"
	"net/http"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"

	"codeberg.org/Kaamkiya/depogit/internal/tmpl"
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

	tmpl.Templates.ExecuteTemplate(w, "repo_log", data)
}

func repoIndex(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Name": r.PathValue("repo"),
	}

	tmpl.Templates.ExecuteTemplate(w, "repo_index", data)
}

func repoTree(w http.ResponseWriter, r *http.Request) {
	repoName := r.PathValue("repo")
	ref := r.PathValue("ref")
	path := r.PathValue("path")

	repo, err := git.PlainOpen(repoPath + repoName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to open repo %s: %v\n", repoName, err)
		return
	}

	hash, err := repo.ResolveRevision(plumbing.Revision(ref))

	commit, err := repo.CommitObject(*hash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to get repo %s ref: %v\n", repoName, err)
		return
	}

	tree, err := commit.Tree()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to get repo %s tree: %v\n", repoName, err)
		return
	}

	var files []object.TreeEntry

	if path == "" {
		files = tree.Entries
	} else {
		obj, err := tree.FindEntry(path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Failed to find path entry: %v\n", err)
			return
		}

		if !obj.Mode.IsFile() {
			subtree, err := tree.Tree(path)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("Failed to find path entry: %v\n", err)
				return
			}

			files = subtree.Entries
		}
	}

	data := make(map[string]any)

	data["Name"] = repoName
	data["Parent"] = path
	data["Files"] = files
	data["Ref"] = ref
	data["Path"] = path

	tmpl.Templates.ExecuteTemplate(w, "tree", data)
}

func repoFile(w http.ResponseWriter, r *http.Request) {
	repoName := r.PathValue("repo")
	path := r.PathValue("path")
	ref := r.PathValue("ref")

	repo, err := git.PlainOpen(repoPath + repoName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to open repo %s: %v\n", repoName, err)
		return
	}

	hash, err := repo.ResolveRevision(plumbing.Revision(ref))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to get repo %s ref: %v\n", repoName, err)
		return
	}

	commit, err := repo.CommitObject(*hash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to get repo %s commit: %v\n", repoName, err)
		return
	}

	tree, err := commit.Tree()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to open repo %s tree: %v\n", repoName, err)
		return
	}

	f, err := tree.File(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to get repo %s file: %v\n", repoName, err)
		return
	}

	contents, err := f.Contents()
	if r.URL.Query().Get("raw") == "true" {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Failed to get repo %s file: %v\n", repoName, err)
			return
		}
		w.Write([]byte(contents))
	} else {
		data := make(map[string]any)

		data["Name"] = repoName

		// I don't care about the error. If there is an error, just assume it's not
		// a binary file.
		if isBin, _ := f.IsBinary(); !isBin {
			data["contents"] = contents
		} else {
			data["contents"] = "Not displaying binary file."
		}

		data["Path"] = path

		tmpl.Templates.ExecuteTemplate(w, "repo_file", data)
	}
}
