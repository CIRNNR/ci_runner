package fetcher

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
)

func Fetch(repoURL string, commit string) {
	repo, err := git.PlainClone(commit, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})
	if err != nil {
		panic(err)
	}
	if commit != ""{
		w, err := repo.Worktree()
		if err != nil {
			panic(err)
		}
		err = w.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(commit),
		})
	}
}
