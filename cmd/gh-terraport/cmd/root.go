/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	terraport "github.com/yardbirdsax/gh-terraport"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "terraport",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
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
