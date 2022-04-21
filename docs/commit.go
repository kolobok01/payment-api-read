package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func rangeDate(start, end time.Time) func() time.Time {
	y, m, d := start.Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)

		return addRandomTime(date, 300)
	}
}

func addRandomTime(d time.Time, interval int) time.Time {
	r := rand.Intn(interval)
	added := d.Add(time.Duration(r) * time.Minute)
	if added.IsZero() {
		return d
	}
	return added
}

// Basic example of how to commit changes to the current branch to an existing
// repository.
func main() {
	/*	rand.Seed(time.Now().UnixNano())
		r := rand.Intn(900)
		start := time.Date(2017, 11, 9, 7, 0, 0, 0, time.UTC)
		t1 := start.Add(time.Duration(r) * time.Minute)
		fmt.Println(t1)*/
	//start := time.Now()
	rand.Seed(time.Now().UnixNano())
	CheckArgs("<directory>")
	directory := os.Args[1]
	start := time.Date(2021, 07, 5, 1, 0, 0, 0, time.UTC)
	end := time.Date(2021, 07, 7, 1, 0, 0, 0, time.UTC)
	fmt.Println(start.Format("2006-01-02"), "-", end.Format("2006-01-02"))
	rand.Seed(time.Now().UnixNano())
	for rd := rangeDate(start, end); ; {
		date := rd()
		for i := 0; i < retNumOfItterations(); i++ {
			if date.IsZero() {
				return
			}
			date = addRandomTime(date, 100)
			// Opens an already existing repository.
			r, err := git.PlainOpen(directory)
			CheckIfError(err)

			w, err := r.Worktree()
			CheckIfError(err)

			// ... we need a file to commit so let's create a new file inside of the
			// worktree of the project using the go standard library.
			Info("echo \"hello world!\" > example-git-file")
			filename := filepath.Join(directory, "log.txt")
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				err = ioutil.WriteFile(filename, []byte("hello world!"), 0644)
				CheckIfError(err)
			}

			// Adds the new file to the staging area.
			Info("git add example-git-file")
			_, err = w.Add("log.txt")
			CheckIfError(err)
			// We can verify the current status of the worktree using the method Status.
			Info("git status --porcelain")
			status, err := w.Status()
			CheckIfError(err)

			fmt.Println(status)

			// Commits the current staging area to the repository, with the new file
			// just created. We should provide the object.Signature of Author of the
			// commit Since version 5.0.1, we can omit the Author signature, being read
			// from the git config files.
			Info("git commit -m \"example go-git commit\"")

			commit, err := w.Commit(fmt.Sprintf("%+v", date), &git.CommitOptions{
				Author: &object.Signature{
					Name:  "kolobok01",
					Email: "kolobublik@gmail.com",
					When:  date,
				},
			})

			CheckIfError(err)

			// Prints the current HEAD to verify that all worked well.
			Info("git show -s")
			obj, err := r.CommitObject(commit)
			CheckIfError(err)

			fmt.Println(obj)
			fmt.Println(date.Format("2006-01-02"))
		}

	}

}

// CheckArgs should be used to ensure the right command line arguments are
// passed before executing an example.
func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		Warning("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
		os.Exit(1)
	}
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

func retNumOfItterations() int {
	return rand.Intn(18)
}
