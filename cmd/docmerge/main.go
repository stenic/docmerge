package main

import (
	"os"

	"docmerge/internal/pkg/docmerge"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cfg = docmerge.DMConfig{}
)

var rootCommand = &cobra.Command{
	Use: "mergedocs",
	Run: func(c *cobra.Command, args []string) {
		dm := docmerge.DocMerge{}
		if err := dm.Run(cfg); err != nil {
			logrus.Error(err)
		}
	},
}

func init() {
	rootCommand.Flags().StringVar(&cfg.OutputDir, "output-dir", "./output", "Directory to output all docs")

	rootCommand.Flags().StringVar(&cfg.GithubOwner, "github-owner", "", "Github owner")
	rootCommand.Flags().StringVar(&cfg.GithubToken, "github-token", os.Getenv("DM_GITHUB_TOKEN"), "Github token, please use env DM_GITHUB_TOKEN")
	rootCommand.Flags().StringVar(&cfg.GithubTopicFilter, "github-topic-filter", "", "Github topic filter")

	rootCommand.Flags().StringVar(&cfg.GitlabOwner, "gitlab-owner", "", "Gitlab owner")
	rootCommand.Flags().StringVar(&cfg.GitlabToken, "gitlab-token", os.Getenv("DM_GITLAB_TOKEN"), "Gitlab token, please use env DM_GITLAB_TOKEN")

}

func main() {
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
