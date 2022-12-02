package terraport

// import (
// 	"fmt"
// 	"os"
// 	"path"
// 	"path/filepath"
// 	"regexp"
// 	"testing"

// 	"github.com/h2non/gock"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"github.com/yardbirdsax/gh-terraport/internal/github"
// 	"github.com/yardbirdsax/gh-terraport/internal/testhelper"
// )

// func setupGitRepos(t *testing.T, testName string, owner string) []*testhelper.GitRepo {
//   testDataFolders := []string{
//     "testdata/repo1",
//     "testdata/repo2",
//   }
//   repos := []*testhelper.GitRepo{}

//   for _, dir := range testDataFolders {
//     name := path.Base(dir)
//     gitRepoHelper := testhelper.TestGitRepo(t, owner, name)
//     repos = append(repos, gitRepoHelper)
//     addFilesInFolder(t, ".", dir, gitRepoHelper)
//   }

//   return repos
// }

// /*
// Adds files at a given path to the created Git repository.

// Files should be added at the same path relative to the base directory they reside in.

// Examples where the root folder being added is at path '/folder':

// - original: folder/something.tf, final: something.tf
// - original: folder/folder2/something.tf, final: folder2/something.tf
// */
// func addFilesInFolder(t *testing.T, path string, basePath string, helper *testhelper.GitRepo) {
//   dirContents, err := os.ReadDir(path)
//   require.NoError(t, err)
//   for _, item := range dirContents {
//     if item.IsDir() {
//       addFilesInFolder(t, filepath.Join(path, item.Name()), path, helper)
//       break
//     }

//     filePath := filepath.Join(path, item.Name())
//     contents, err := os.ReadFile(filePath)
//     require.NoError(t, err)
//     helper.CommitFile(filePath, string(contents))
//   }
// }

// func cleanupFolders(repos []*testhelper.GitRepo) {
//   for _, repo := range repos {
//     repo.CleanUp()
//   }
// }

// func TestResultsFromSearch(t *testing.T) {
// 	expectedOwnerName := "owner"
//   repos := setupGitRepos(t, "ResultsFromSearch", expectedOwnerName)
//   defer cleanupFolders(repos)
// 	// expectedRepoName1 := "repo1"
// 	expectedModuleRepoURL1 := "https://github.com/owner/repo1"
// 	// expectedRepoName2 := "repo2"
// 	expectedModuleRepoURL2 := "https://github.com/owner/repo2"
// 	expectedSearchPath := "search/repositories"
// 	searchResultItems := []github.RepoSearchResult{}
//   for _, repo := range repos {
//     searchResultItems = append(searchResultItems, github.RepoSearchResult{
//       CloneURL: repo.Url.String(),
//       FullName: repo.FullName,
//     })
//   }
//   searchResults := github.RepoSearchResults{
//     TotalCount: 2,
//     Items: searchResultItems,
//   }
//   gock.New("https://api.github.com").
// 		Get(expectedSearchPath).
// 		MatchParam("q", fmt.Sprintf(regexp.QuoteMeta("user:%s"), expectedOwnerName)).
// 		MatchParam("page", "1").
// 		MatchParam("per_page", "100").
// 		Reply(200).JSON(
// 		searchResults,
// 	)
// 	expectedResults := &Results{
// 		Modules{
// 			{
// 				RepositoryURL: expectedModuleRepoURL1,
// 				RawURL:        "git::https://github.com/someone/terraform-aws-s3?ref=v1.2.0",
// 				Version:       "1.2.0",
// 				FileName:      "git_with_version.tf",
// 			},
// 			{
// 				RepositoryURL: expectedModuleRepoURL1,
// 				RawURL:        "someone-else/terraform-aws-sqs",
// 				Version:       "1.2.0",
// 				FileName:      "registry_with_version.tf",
// 			},
// 			{
// 				RepositoryURL: expectedModuleRepoURL2,
// 				RawURL:        "git::https://github.com/someone/terraform-aws-vpc?ref=v1.2.0",
// 				Version:       "1.2.0",
// 				FileName:      "git_with_version.tf",
// 			},
// 			{
// 				RepositoryURL: expectedModuleRepoURL2,
// 				RawURL:        "someone-else/terraform-aws-rds",
// 				Version:       "1.2.0",
// 				FileName:      "registry_with_version.tf",
// 			},
// 		},
// 		nil,
// 	}

// 	actualResults, err := FromSearch(fmt.Sprintf("user:%s", expectedOwnerName))

// 	assert.Equal(t, expectedResults, actualResults)
// 	assert.NoError(t, err)
// 	assert.False(t, gock.IsPending())
// }
