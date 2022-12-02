package testhelper

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// This code is taken almost verbatim from https://github.com/hashicorp/go-getter/.

// GitRepo is a helper struct which controls a single temp git repo.
type GitRepo struct {
	t   *testing.T
	Url *url.URL
  Owner string
  Name string
  FullName string
	dir string
}

// TestGitRepo creates a new test git repository.
func TestGitRepo(t *testing.T, owner string, name string) *GitRepo {
	dir, err := os.MkdirTemp("", "go-getter")
	if err != nil {
		t.Fatal(err)
	}
	dir = filepath.Join(dir, name)
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}

	r := &GitRepo{
		t:   t,
		dir: dir,
    Name: name,
    Owner: owner,
    FullName: fmt.Sprintf("%s/%s", owner, name),
	}

	url, err := url.Parse("file://" + r.dir)
	if err != nil {
		t.Fatal(err)
	}
	r.Url = url

	t.Logf("initializing git repo in %s", dir)
	r.git("init")
	r.git("config", "user.name", "go-getter")
	r.git("config", "user.email", "go-getter@hashicorp.com")

	return r
}

// git runs a git command against the repo.
func (r *GitRepo) git(args ...string) {
	cmd := exec.Command("git", args...)
	cmd.Dir = r.dir
	bfr := bytes.NewBuffer(nil)
	cmd.Stderr = bfr
	if err := cmd.Run(); err != nil {
		r.t.Fatal(err, bfr.String())
	}
}

// CommitFile writes and commits a text file to the repo.
func (r *GitRepo) CommitFile(fileName, content string) {
  path := filepath.Join(r.dir, fileName)
  if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
    r.t.Fatal(err)
  }
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		r.t.Fatal(err)
	}
	r.git("add", fileName)
	r.git("commit", "-m", "Adding "+fileName)
}

func (r *GitRepo) CleanUp() {
  err := os.RemoveAll(r.dir)
  if err != nil {
    r.t.Fatal(err)
  }
}