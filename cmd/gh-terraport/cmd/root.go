/*
Copyright © 2022 Josh Feierman <josh@sqljosh.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	terraport "github.com/yardbirdsax/gh-terraport"
)

var rootCmd = &cobra.Command{
	Use:   "terraport",
	Short: "An extension for the GitHub CLI for getting information about Terraform module usage in GitHub",
	Long: `An extension for the GitHub CLI that retrieves and presents information about Terraform module usage within repositories in GitHub.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		optFns := []terraport.OptFn{
			terraport.WithConcurrency(maxConcurrency),
		}
		if ignoreLocal {
			optFns = append(optFns, terraport.WithExcludeLocalModules())
		}
		results, err := terraport.FromSearch(search, optFns...)
		if results != nil {
			fmt.Print(results.AsCSV())
		} else {
			fmt.Printf("error: %v", err)
		}
		return err
	},
}

var search string
var ignoreLocal bool
var maxConcurrency int

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&search, "search", "s", "", "The search query to use when searching for repositories")
	//nolint:errcheck
	rootCmd.MarkFlagRequired("search")
	rootCmd.Flags().BoolVarP(&ignoreLocal, "ignore-local", "", false, "Ignore modules with local paths")
	rootCmd.Flags().IntVarP(&maxConcurrency, "max-concurrency", "m", 10, "The maximum number of concurrent asynchronous operations")
}
