package docmerge

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type DocMerge struct {
	cfg DMConfig
	ctx context.Context
}

func (d DocMerge) Run(cfg DMConfig) error {
	d.ctx = context.Background()
	d.cfg = cfg
	d.cfg.DocsDir = "docs"

	logrus.SetLevel(logrus.DebugLevel)

	if cfg.GithubOwner != "" {
		if cfg.GithubToken == "" {
			return fmt.Errorf("github-owner and github-token are required")
		}

		logrus.Infof("Processing github %s", cfg.GithubOwner)
		if err := d.runGithub(cfg.GithubOwner); err != nil {
			logrus.Error(err)
		}
	}

	if cfg.GitlabOwner != "" {
		if cfg.GitlabToken == "" {
			return fmt.Errorf("gitlab-owner and gitlab-token are required")
		}

		logrus.Infof("Processing gitlab %s", cfg.GitlabOwner)
		if err := d.runGitlab(cfg.GitlabOwner); err != nil {
			logrus.Error(err)
		}
	}

	return nil
}
