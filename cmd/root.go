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
	asanaToken              string
	asanaProjectId          string
	githubOrganization      string
	githubRepository        string
	githubToken             string
	dryRun                  bool
	githubDelayForRateLimit int
	forceUpdate             bool

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
		--github-repository <github-repo-name>

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
				githubClient = github.New(ctx, githubToken, githubDelayForRateLimit)
			}

			migratedIssues, err := github.ListMigratedIssues(ctx, githubClient, githubToken, githubOrganization, githubRepository)
			if err != nil {
				fmt.Printf("failed to retrieve already migrated issues: %v\n", err)
				return
			}

			successCount := 0
			completedStateCount := 0
			for i, t := range tasks {
				fmt.Printf("i: %v\n", i)

				labels := []string{}
				for _, tag := range t.Tags {
					labels = append(labels, tag.Name)
				}

				err := github.MigrateIssue(ctx, githubClient, migratedIssues, githubOrganization, githubRepository, t.Name, t.Notes, labels, t.Completed, t.PermalinkURL, forceUpdate)
				if err != nil {
					fmt.Printf("failed to migrate issue %v: %v\n", t.Name, err)
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
	if err := rootCmd.MarkFlagRequired("asana-token"); err != nil {
		fmt.Println("failed to mark flag as required")
	}

	rootCmd.Flags().StringVar(&asanaProjectId, "asana-project-id", "", "ID of the project whose task we want to export")
	if err := rootCmd.MarkFlagRequired("asana-project-id"); err != nil {
		fmt.Println("failed to mark flag as required")
	}

	rootCmd.Flags().StringVar(&githubToken, "github-token", "", "Github personal access token with access to create and modify Projects")
	if err := rootCmd.MarkFlagRequired("github-token"); err != nil {
		fmt.Println("failed to mark flag as required")
	}

	rootCmd.Flags().StringVar(&githubOrganization, "github-organization", "", "Name of the Github organization")
	if err := rootCmd.MarkFlagRequired("github-organization"); err != nil {
		fmt.Println("failed to mark flag as required")
	}

	rootCmd.Flags().StringVar(&githubRepository, "github-repository", "", "Name of the Github repository")
	if err := rootCmd.MarkFlagRequired("github-repository"); err != nil {
		fmt.Println("failed to mark flag as required")
	}

	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print out tasks that will be migrated without actually migrating them")
	rootCmd.Flags().BoolVar(&forceUpdate, "force-update", false, "Always override the issue in github with the data from asana")
	rootCmd.Flags().IntVar(&githubDelayForRateLimit, "github-rate-limit-delay", 10, "Delay to apply between github API requests to avoid hitting rate limits")
}
