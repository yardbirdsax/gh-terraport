package terraport

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hashicorp/go-getter"
	"github.com/yardbirdsax/gh-terraport/internal/github"
	"github.com/yardbirdsax/gh-terraport/internal/result"
	"github.com/yardbirdsax/gh-terraport/internal/terraform"
)

type Terraport struct {
	ignoreLocal bool
	// A map / cache of modules and their most up to date versions
	moduleVersionCache map[string]string
	// A mutex used to make the method that updates the cache thread safe
	moduleVersionCacheMutex sync.Mutex
	// The max concurrency for asynchronous operations
	maxConcurrency int
}

type OptFn func(*Terraport) error

// WithExcludeLocalModules excludes any modules with a local file path
// as the source.
func WithExcludeLocalModules() OptFn {
	return func(t *Terraport) error {
		t.ignoreLocal = true
		return nil
	}
}

// WithConcurrency sets the max concurrency for asynchronous operations
func WithConcurrency(maxConcurrency int) OptFn {
	return func(t *Terraport) error {
		t.maxConcurrency = maxConcurrency
		return nil
	}
}

func FromSearch(search string, optFns ...OptFn) (*result.Results, error) {
	t := &Terraport{}

	for _, f := range optFns {
		err := f(t)
		if err != nil {
			return &result.Results{}, err
		}
	}

	return t.FromSearch(search)
}

func (t *Terraport) FromSearch(search string) (*result.Results, error) {
	rawResults := [][]interface{}{}
	var results *result.Results
	optFns := []terraform.OptFn{}
	mu := &sync.Mutex{}
	wg := sync.WaitGroup{}
	coordinatorChan := make(chan bool, t.maxConcurrency)
	defer close(coordinatorChan)

	if t.ignoreLocal {
		optFns = append(optFns, terraform.WithExcludeLocalModules())
	}

	repos, err := github.ReposFromSearch(search,)
	if err != nil {
		return results, err
	}
	wg.Add(len(repos))

	for i, repo := range repos {
		log.Printf("beginning repo: %s, number %d of %d", repo.FullName, i, len(repos))
		coordinatorChan <- true
		go func(repo github.RepoSearchResult) {
			defer wg.Done()
			defer func() { <- coordinatorChan }()

			tmpDir := filepath.Join(os.TempDir(), repo.FullName)
			gitURL := fmt.Sprintf("git::%s", repo.CloneURL)

			err := getter.Get(tmpDir, gitURL)
			if err != nil {
				log.Printf("error processing repo '%s': %v", repo.FullName, err)
				return
			}
			defer os.RemoveAll(tmpDir)

			err = filepath.Walk(tmpDir, func(path string, info fs.FileInfo, err error) error {
				if filepath.Ext(path) == ".tf" {
					terraformFile, err := terraform.FromFile(path, optFns...)
					if err != nil {
						return err
					}

					pathWithoutBaseDirectory, err := filepath.Rel(tmpDir, path)
					if err != nil {
						return err
					}

					for _, m := range terraformFile.Modules {
						mu.Lock()
						var moduleLatestVersion string = ""
						var isUpToDate string = ""
						if m.GitHubRepositoryFullName != "" {
							moduleLatestVersion = t.getLatestRelease(m.GitHubRepositoryFullName)
							if m.Version == moduleLatestVersion {
								isUpToDate = "true"
							} else {
								isUpToDate = "false"
							}
						}
						rawResults = append(rawResults, []interface{}{ repo.FullName, pathWithoutBaseDirectory, m.Source, m.Version, moduleLatestVersion, isUpToDate})
						mu.Unlock()
					}
				}
				return nil
			})
			if err != nil {
				log.Printf("error processing repo '%s': %v", repo.FullName, err)
				return
			}
		}(repo)
	}

	wg.Wait()

	if len(rawResults) > 0 {
		results, err = result.FromSlice(
			[]interface{}{ "Repo Full Name", "File Path", "Module Formatted Source", "Module Version", "Latest Version", "Up To Date"},
			rawResults,
		)
	}
	return results, err
}

func (t *Terraport) getLatestRelease(repositoryFullName string) string {
	if version, exists := t.moduleVersionCache[repositoryFullName]; exists {
		return version
	} else {
		t.moduleVersionCacheMutex.Lock()
		defer t.moduleVersionCacheMutex.Unlock()

		releases, err := github.ListReleases(repositoryFullName)
		if err != nil {
			log.Printf("error retrieving latest module version: %v", err)
			return ""
		}
		if len(releases) == 0 {
			return ""
		}
		latestRelease := releases[0]

		return strings.ReplaceAll(latestRelease.Name, "v", "")
	}
}
