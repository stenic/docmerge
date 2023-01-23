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
	Use:         "mergedocs endpoint",
	Args: cobra.ExactArgs(1),
	Run: func(c *cobra.Command, args []string) {
		dm := docmerge.DocMerge{}
		if err := dm.Run(cfg, args[0]); err != nil {
			logrus.Error(err)
		}
	},
}

func init() {
	rootCommand.Flags().StringVar(&cfg.OutputDir, "output-dir", "./output", "Directory to output all docs")

	rootCommand.Flags().StringVar(&cfg.Owner, "owner", "", "Owner")
	rootCommand.Flags().StringVar(&cfg.Token, "token", os.Getenv("DM_TOKEN"), "Token, please use env DM_TOKEN")
	rootCommand.Flags().StringVar(&cfg.TopicFilter, "topic-filter", "", "Topic filter")
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
