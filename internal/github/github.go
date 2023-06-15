package github

import (
	"fmt"
	"log"
	"net/url"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/cli/go-gh/pkg/auth"
)

type RepoSearchResult struct {
  FullName string `json:"full_name"`
  CloneURL string `json:"clone_url"`
}

type RepoSearchResults struct {
  TotalCount int `json:"total_count"`
  Items []RepoSearchResult `json:"items"`
}

func ReposFromSearch(search string) ([]RepoSearchResult, error) {
  results := []RepoSearchResult{}
  rawResults := RepoSearchResults{}


  client, err := getClient()
  if err != nil {
    return results, err
  }

  sanitizedSearch := url.QueryEscape(search)
  totalCount := 0
  page := 1
  for {
    if page > 10 {
      log.Printf("[WARN] retrieved the maximum number of search results, consider revising your search criteria")
      log.Printf("[WARN] incomplete results may be shown")
      return results, nil
    }
    err := client.Get(fmt.Sprintf("search/repositories?q=%s&page=%d&per_page=100", sanitizedSearch, page), &rawResults)
    if err != nil {
      return results, err
    }
    results = append(results, rawResults.Items...)
    log.Printf("retrieved %d items of %d total", len(results), rawResults.TotalCount)

    totalCount = len(results)
    if totalCount >= rawResults.TotalCount {
      break
    }
    page++
  }

  return results, nil
}

type ListReleaseResult struct {
  Name string `json:"name"`
}

type ListReleaseResults []ListReleaseResult

func ListReleases(repositoryFullName string) (ListReleaseResults, error) {
  results := ListReleaseResults{}

  client, err := getClient()
  if err != nil {
    return results, err
  }

  listPath := fmt.Sprintf("repos/%s/releases", repositoryFullName)
  err = client.Get(listPath, &results)

  return results, err
}

func getClient() (api.RESTClient, error) {
	token, _ := auth.TokenForHost("github.com")
	return gh.RESTClient(&api.ClientOptions{
		AuthToken: token,
	})
}
