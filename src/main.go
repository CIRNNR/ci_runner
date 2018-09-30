package main

import (
	"dockerdemo/src/fetcher"
	"dockerdemo/src/runner"
	"os"
)

func main() {
	var repo = "https://github.com/kh3dr0n/testrepo.git"
	var commit = "6e179f4ebe7267331df3aa28ecc5f5a2f2ecfd0c"
	ci(repo, commit)
}

func ci(repo string, commit string) {
	os.RemoveAll("./" + commit)
	fetcher.Fetch(repo, commit)
	runner.Run(commit)
	os.RemoveAll("./" + commit)
}
