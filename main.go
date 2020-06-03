package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/v32/github"
	"github.com/spf13/cobra"
)

var (
	ghRepository        string
	ghOrganization      string
	targetTag           string
	branch              string
	releaseNoteTemplate string

	//bavery_todo:cleanup
	rootCmd = cobra.Command{
		Use:   "irmt",
		Short: "Istio release management tool",
	}

	version    = "-0.0.1"
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run:   printVersion,
	}

	releaseNotesCmd = &cobra.Command{
		Use:   "releasenotes",
		Short: "Generate release notes",
		Run:   genReleaseNotes,
	}

	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Generate release status",
		Run:   releaseStatus,
	}

	cutReleaseCmd = &cobra.Command{
		Use:   "cutrelease",
		Short: "Cut release",
		Run:   cutRelease,
	}
)

func printVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("Version: %v", version)
}

func httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func genReleaseNotes(cmd *cobra.Command, args []string) {
	fmt.Printf("Generating release notes:%+v\n\n", args)

	ghClient := github.NewClient(nil)
	prOptions := &github.PullRequestListOptions{
		State: "all",
		Base:  branch,
	}
	pulls, _, err := ghClient.PullRequests.List(context.Background(), ghOrganization, ghRepository, prOptions)
	if err != nil {
		fmt.Printf("Failed to list pulls: %+v", err)
	}

	for _, pull := range pulls {
		mergedAt := ""
		if pull.MergedAt != nil {
			mergedAt = pull.MergedAt.Format("2006 01/02-15:04:05.000")
		}

		labels := func(labels []*github.Label) []string {
			validLabels := make([]string, 0)
			for _, label := range labels {
				if label != nil {
					validLabels = append(validLabels, *(*label).Name)
				}
			}
			return validLabels
		}(pull.Labels)

		fmt.Printf("Pull %d: \n\t state: %s \n\t title: %s \n\t body: %+v\n\t labels: %+v\n\t mergedat:%s\n\t IssueURL:%s\n\t HTMLUrl:%s \n\t\n",
			*pull.ID,
			*pull.State,
			*pull.Title,
			*pull.Body,
			labels,
			mergedAt,
			*pull.IssueURL,
			*pull.HTMLURL)
	}

}

func cut(cmd *cobra.Command, args []string) {
	fmt.Printf("release notes:%+v", args)
}

func releaseStatus(cmd *cobra.Command, args []string) {
	fmt.Printf("Release isn't done yet. Please try again")

	pulls, _, err := ghClient.PullRequests.List(context.Background(), ghOrganization, ghRepository, prOptions)
	if err != nil {
		fmt.Printf("Failed to list pulls: %+v", err)
	}

	for _, pull := range pulls {
		mergedAt := ""
		if pull.MergedAt != nil {
			mergedAt = pull.MergedAt.Format("2006 01/02-15:04:05.000")
		}

		labels := func(labels []*github.Label) []string {
			validLabels := make([]string, 0)
			for _, label := range labels {
				if label != nil {
					validLabels = append(validLabels, *(*label).Name)
				}
			}
			return validLabels
		}(pull.Labels)

		ghClient := github.NewClient(nil)
		prOptions := &github.PullRequestListOptions{
			State: "all",
			Base:  branch,
		}

		complete := 0
		remaining := 0
		needsDocumentation := 0
		for label := range labels {
			if pull.state == "open" {
				fmt.Printf("Pull %d: \n\t state: %s \n\t title: %s \n\t body: %+v\n\t labels: %+v\n\t mergedat:%s\n\t IssueURL:%s\n\t HTMLUrl:%s \n\t\n",
					*pull.ID,
					*pull.State,
					*pull.Title,
					*pull.Body,
					labels,
					mergedAt,
					*pull.IssueURL,
					*pull.HTMLURL)
				remaining++
			}

			complete++
			continue
		}
	}

}

func cutRelease(cmd *cobra.Command, args []string) {
	fmt.Printf("Cutting release")
}

func main() {
	rootCmd.PersistentFlags().StringVar(&ghRepository, "ghrepository", "istio", "GitHub repository")
	rootCmd.PersistentFlags().StringVar(&ghOrganization, "ghorg", "istio", "GitHub organization")
	rootCmd.PersistentFlags().StringVar(&branch, "branch", "", "Repository branch")
	versionCmd.Flags().StringVar(&releaseNoteTemplate, "template", "", "Release notes template.")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(releaseNotesCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(cutReleaseCmd)
	rootCmd.Execute()
}
