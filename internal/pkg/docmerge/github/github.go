package github

import (
	"context"
	"docmerge/internal/pkg/docmerge/adapter"
	"io"
	"net/http"
	"os"

	ghapi "github.com/google/go-github/v49/github"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func New(token string) adapter.Adapter {
	tempDir := "/tmp/docmerge-cache"
	tc := &http.Client{
		Transport: &oauth2.Transport{
			Base: httpcache.NewTransport(diskcache.New(tempDir)),
			Source: oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			),
		},
	}

	client := ghapi.NewClient(tc)

	return github{
		client:  client,
		context: context.Background(),
	}
}

type github struct {
	client  *ghapi.Client
	context context.Context
}

func (a github) GetRepositories(owner string) chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		opts := &ghapi.RepositoryListByOrgOptions{}

		for {
			repos, resp, err := a.client.Repositories.ListByOrg(a.context, owner, opts)
			if err != nil {
				logrus.Error(err)
				continue
			}
			for _, repo := range repos {
				ch <- *repo.Name
			}
			if resp.NextPage == 0 {
				return
			}
			opts.Page = resp.NextPage
		}
	}()

	return ch
}
func (a github) GetFiles(owner, repo string) chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		tree, _, err := a.client.Git.GetTree(a.context, owner, repo, "HEAD", true)
		if err != nil {
			logrus.Error(err)
			return
		}

		for _, entry := range tree.Entries {
			if entry.GetType() != "blob" {
				continue
			}
			ch <- entry.GetPath()
		}
	}()

	return ch
}

func (a github) DownloadFile(owner, repo, file, localPath string) error {
	iorc, _, err := a.client.Repositories.DownloadContents(a.context, owner, repo, file, &ghapi.RepositoryContentGetOptions{})
	if err != nil {
		return err
	}
	defer iorc.Close()

	outFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, iorc); err != nil {
		return err
	}

	return err
}
