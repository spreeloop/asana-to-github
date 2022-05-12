package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"spreeloop.com/asana-to-github/asana"
	"spreeloop.com/asana-to-github/github"
)

var (
	asanaToken         string
	asanaProjectId     string
	githubOrganization string
	githubRepository   string
	githubToken        string
	dryRun             bool

	// rootCmd represents the base command when called without any subcommands.
	rootCmd = &cobra.Command{
		Use:   "asana-to-github",
		Short: "Migration from Asana to Github Project",
		Long: `Command line tool for migrating tasks from Asana
	to Github project. e.g:

	asana-to-github --asana-token <asana-personal-access-token> \
		--asana-project-id <asana-project-id> \
		--github-token <github-personal-access-token> \
		--github-organization <github-org-name> \
		--github-repository <github-repo-name> \
		--dry-run true|false

	How to export asana tasks to JSON:
	https://asana.com/guide/help/faq/security`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			asanaClient := asana.New(ctx, asanaToken)
			tasks, err := asana.FetchTasks(ctx, asanaClient, asanaProjectId)
			if err != nil {
				fmt.Printf("failed to retrieve asana tasks: %v\n", err)
				return
			}

			githubClient := github.NewFake()
			if !dryRun {
				githubClient = github.New(ctx, githubToken)
			}

			successCount := 0
			completedStateCount := 0
			for _, t := range tasks {
				labels := []string{}
				for _, tag := range t.Tags {
					labels = append(labels, tag.Name)
				}

				err := github.CreateIssue(ctx, githubClient, githubOrganization, githubRepository, t.Name, t.Notes, labels, t.Completed)
				if err != nil {
					fmt.Printf("failed to create issue %v: %v\n", t.Name, err)
					continue
				}

				if t.Completed {
					completedStateCount++
				}
				successCount++
			}

			fmt.Printf("Successfully migrated %v out of %v tasks\n", successCount, len(tasks))
			fmt.Printf("%v out of %v tasks are closed\n", completedStateCount, len(tasks))
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
	rootCmd.Flags().StringVar(&asanaToken, "asana-token", "", "Asana personal access token with access to read project tasks")
	rootCmd.MarkFlagRequired("asana-token")

	rootCmd.Flags().StringVar(&asanaProjectId, "asana-project-id", "", "ID of the project whose task we want to export")
	rootCmd.MarkFlagRequired("asana-project-id")

	rootCmd.Flags().StringVar(&githubToken, "github-token", "", "Github personal access token with access to create and modify Projects")
	rootCmd.MarkFlagRequired("github-token")

	rootCmd.Flags().StringVar(&githubOrganization, "github-organization", "", "Name of the Github organization")
	rootCmd.MarkFlagRequired("github-organization")

	rootCmd.Flags().StringVar(&githubRepository, "github-repository", "", "Name of the Github repository")
	rootCmd.MarkFlagRequired("github-repository")

	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print out tasks that will be migrated without actually migrating them")
}
