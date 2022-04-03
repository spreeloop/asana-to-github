package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	pathAsanaJSON      string
	githubOrganization string
	githubRepository   string
	githubToken        string

	// rootCmd represents the base command when called without any subcommands.
	rootCmd = &cobra.Command{
		Use:   "asana-to-github",
		Short: "Migration from Asana to Github Project",
		Long: `Command line tool for migrating tasks from Asana
	to Github project. e.g:

	asana-to-github --asana-json </path/to/exported_asana_tasks.json> \
		--github-token <github-personal-access-token> \
		--github-organization <github-org-name> \
		--github-repository <github-repo-name> \

	How to export asana tasks to JSON:
	https://asana.com/guide/help/faq/security`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Echo: " + strings.Join(args, " "))
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&pathAsanaJSON, "asana-json", "", "Relative or absolute path to the exported Asana tasks in JSON format")
	rootCmd.MarkFlagRequired("asana-json")

	rootCmd.Flags().StringVar(&githubToken, "github-token", "", "Github personal access token with access to create and modify Projects")
	rootCmd.MarkFlagRequired("github-token")

	rootCmd.Flags().StringVar(&githubOrganization, "github-organization", "", "Name of the Github organization")
	rootCmd.MarkFlagRequired("github-organization")

	rootCmd.Flags().StringVar(&githubRepository, "github-repository", "", "Name of the Github repository")
	rootCmd.MarkFlagRequired("github-repository")
}
