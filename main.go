package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/google/go-github/v32/github"	
)

var (
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

func getCommitForTagFromBody(targetTagName string, body []byte) (map[string]interface{}, error) {
	var dataDump []map[string]interface{}
	if err := json.Unmarshal(body, &dataDump); err != nil {
		return nil, err
	}

	for _, tag := range dataDump {
		tagName := tag["name"]
		if targetTagName == tagName {
			return tag, nil
		}
	}

	return nil, fmt.Errorf("not found\n")
}

func getCommitFromTag(tag map[string]interface{}) (map[string]interface{}, error) {
	commit, success := tag["commit"].(map[string]interface{})
	if !success {
		return nil, fmt.Errorf("Can't parse commit\n")

	}
	return commit, nil
}

func genReleaseNotes(cmd *cobra.Command, args []string) {
	fmt.Printf("Generating release notes:%+v\n\n", args)

	targetTag := "1.6.0"
	body, err := httpGet(fmt.Sprintf("https://api.github.com/repos/istio/istio/commits/%s", targetTag))
	if err != nil {
		fmt.Printf("Failed to get commit: %+v", err)
		return
	}

	var commitBody map[string]interface{}
	err = json.Unmarshal(body, &commitBody)
	if err != nil {
		fmt.Printf("Failed to parse commit body: %+v", err)
		return
	}

	fmt.Printf("commiturl: %s", commitBody)
	commitCommit := commitBody["commit"].(map[string]interface{})
	commitCommitter, success := commitCommit["committer"].(map[string]interface{})
	if !success {
		fmt.Printf("Failed to parse committer: %+v", err)
		return
	}

	fmt.Printf("Commit for %s:%s published: %+v\n", targetTag, commitCommitter["date"])

}

func cut(cmd *cobra.Command, args []string) {
	fmt.Printf("release notes:%+v", args)
}

func releaseStatus(cmd *cobra.Command, args []string) {
	fmt.Printf("Release isn't done yet. Please try again")
}

func cutRelease(cmd *cobra.Command, args []string) {
	fmt.Printf("Cutting release")
}

func main() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(releaseNotesCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(cutReleaseCmd)
	rootCmd.Execute()
}
