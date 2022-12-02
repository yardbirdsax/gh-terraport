// Package terraform contains structures related to parsing Terraform files.
package terraform

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

// TerraformFile represents a single Terraform file.
type TerraformFile struct {
	// Modules are the 'module' blocks present in the file.
	Modules []Module `hcl:"module,block"`
	// The name of the file
	FileName string
}

type OptFn func(*TerraformFile) error

// WithExcludeLocalModules will ignore any module blocks that reference local paths.
func WithExcludeLocalModules() OptFn {
	return func(tf *TerraformFile) error {
		for i := 0; i < len(tf.Modules); i++ {
			if i >= len(tf.Modules) {
				break
			}
			if strings.HasPrefix(tf.Modules[i].RawSource, ".") {
				tf.Modules = append(tf.Modules[0:i], tf.Modules[i+1:]...)
				i--
			}
		}
		return nil
	}
}

func (t *TerraformFile) parse() error {
	for i := range t.Modules {
		// version
		if t.Modules[i].Version == "" {
			if matches := regexp.MustCompile("ref=(v)?([^&]+)").FindAllStringSubmatch(t.Modules[i].RawSource, -1); matches != nil {
				t.Modules[i].Version = matches[0][2]
			}
		}
		// Source
		if regexp.MustCompile("^(git::|git@)").MatchString(t.Modules[i].RawSource) {
			moduleURL, err := url.Parse(
				regexp.MustCompile("(git::https://|git@|git::ssh://git@)").ReplaceAllString(
					strings.ReplaceAll(
						strings.ReplaceAll(t.Modules[i].RawSource, "github.com:", "github.com/"),
						".git",
						"",
					),
					"https://",
				),
			)
			if err != nil {
				return err
			}
			t.Modules[i].Source = fmt.Sprintf("https://%s%s", moduleURL.Host, moduleURL.Path)
		} else {
			t.Modules[i].Source = t.Modules[i].RawSource
		}
		// GitHub full name
		if strings.HasPrefix(t.Modules[i].Source, "https://github.com") {
			moduleURL, err := url.Parse(strings.ReplaceAll(t.Modules[i].Source, ".git", ""))
			if err != nil {
				return err
			}

			splitPath := strings.Split(moduleURL.Path, "/")
			repoFullName := strings.Join(splitPath[1:3], "/")
			t.Modules[i].GitHubRepositoryFullName = repoFullName
		}
	}

	return nil
}

// Module is a single "module" block
type Module struct {
	Name                     string `hcl:",label"`
	RawSource                string `hcl:"source"`
	Source                   string
	Version                  string `hcl:"version,optional"`
	SourceType               string
	GitHubRepositoryFullName string
}

// FromFile creates a new TerraformFile struct from a file at a given path
func FromFile(filePath string, optFns ...OptFn) (*TerraformFile, error) {
	var terraformFile *TerraformFile = &TerraformFile{}

	f, err := os.Open(filePath)
	if err != nil {
		return terraformFile, err
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return terraformFile, err
	}

	//nolint:errcheck
	hclsimple.Decode("placeholder.hcl", content, nil, terraformFile)
	terraformFile.FileName = filePath

	err = terraformFile.parse()
	if err != nil {
		return terraformFile, err
	}

	for _, f := range optFns {
		err = f(terraformFile)
		if err != nil {
			return terraformFile, err
		}
	}

	return terraformFile, err
}
