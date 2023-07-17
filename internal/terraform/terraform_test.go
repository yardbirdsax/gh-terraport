package terraform

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromFile(t *testing.T) {
	testCases := []struct {
		name                             string
		expectedName                     string
		expectedRawSource                string
		expectedSource                   string
		expectedVersion                  string
		expectedGitHubRepositoryFullName string
		expectedPath                     string
		optFns                           []OptFn
	}{
		{
			name:                             "git_with_version",
			expectedName:                     "something",
			expectedRawSource:                "git::https://github.com/cloudposse/terraform-aws-vpc?ref=v1.2.0",
			expectedSource:                   "https://github.com/cloudposse/terraform-aws-vpc",
			expectedVersion:                  "1.2.0",
			expectedGitHubRepositoryFullName: "cloudposse/terraform-aws-vpc",
			expectedPath:                     "../../testdata",
		},
		{
			name:                             "registry_with_version",
			expectedName:                     "something_else",
			expectedRawSource:                "cloudposse/terraform-aws-vpc",
			expectedSource:                   "cloudposse/terraform-aws-vpc",
			expectedGitHubRepositoryFullName: "",
			expectedVersion:                  "1.2.0",
			expectedPath:                     "../../testdata",
		},
		{
			name:                             "with_local",
			expectedName:                     "remote",
			expectedRawSource:                "git::https://github.com/something/module?ref=1.0.0",
			expectedSource:                   "https://github.com/something/module",
			expectedVersion:                  "1.0.0",
			expectedGitHubRepositoryFullName: "something/module",
			optFns:                           []OptFn{WithExcludeLocalModules()},
			expectedPath:                     "../../testdata",
		},
		{
			name:                             "ssh",
			expectedName:                     "ssh",
			expectedRawSource:                "git::ssh://git@github.com/something/something?ref=1.0.0",
			expectedSource:                   "https://github.com/something/something",
			expectedVersion:                  "1.0.0",
			expectedGitHubRepositoryFullName: "something/something",
			expectedPath:                     "../../testdata",
		},
		{
			name:                             "with_dot_git",
			expectedName:                     "git",
			expectedRawSource:                "git::https://github.com/something/something.git?ref=1.0.0",
			expectedSource:                   "https://github.com/something/something",
			expectedVersion:                  "1.0.0",
			expectedGitHubRepositoryFullName: "something/something",
			expectedPath:                     "../../testdata",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			file := fmt.Sprintf("../../testdata/%s.tf", tc.name)
			terraform, err := FromFile(file, tc.optFns...)

			assert.NoError(t, err)
			assert.Equal(t, file, terraform.FileName)
			assert.Equal(t, tc.expectedRawSource, terraform.Modules[0].RawSource)
			assert.Equal(t, tc.expectedVersion, terraform.Modules[0].Version)
			assert.Equal(t, tc.expectedName, terraform.Modules[0].Name)
			assert.Equal(t, tc.expectedGitHubRepositoryFullName, terraform.Modules[0].GitHubRepositoryFullName)
			assert.Equal(t, tc.expectedSource, terraform.Modules[0].Source)
			assert.Equal(t, tc.expectedPath, terraform.Path)
		})
	}
}
