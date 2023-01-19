package docmerge

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/sirupsen/logrus"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v49/github"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"golang.org/x/exp/slices"
)

func Run(cfg DMConfig) error {
	if cfg.GithubOwner == "" || cfg.GithubToken == "" {
		return fmt.Errorf("github-owner and github-token are required")
	}

	ctx := context.Background()
	tempDir := "/tmp/docmerge-cache"
	tc := &http.Client{
		Transport: &oauth2.Transport{
			Base: httpcache.NewTransport(diskcache.New(tempDir)),
			Source: oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: cfg.GithubToken},
			),
		},
	}

	client := github.NewClient(tc)
	gh := Github{
		client:  client,
		context: ctx,
		owner:   cfg.GithubOwner,
		log:     logrus.NewEntry(logrus.New()),
	}

	repos, err := gh.getAllRepos(nil)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		gh.log = logrus.WithField("repo", repo.GetName())
		if cfg.GithubTopicFilter != "" && !slices.Contains(repo.Topics, cfg.GithubTopicFilter) {
			gh.log.Debug("Skipping: does not have label 'internal-docs'")
			continue
		}

		gh.log.Infof("Fetching %s", repo.GetName())
		if err := gh.downloadDocs(repo, path.Join(cfg.OutputDir, *repo.Name)); err != nil {
			gh.log.Error(err)
		}
	}

	return nil
}
