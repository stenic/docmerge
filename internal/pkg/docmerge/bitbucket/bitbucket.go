package bitbucket

import (
	"docmerge/internal/pkg/docmerge/adapter"
	"os"
	"strings"

	bbapi "github.com/ktrysmt/go-bitbucket"
	"github.com/sirupsen/logrus"
)

func New(token string) adapter.Adapter {
	parts := strings.Split(token, ":")
	return bitbucket{
		client: bbapi.NewBasicAuth(parts[0], parts[1]),
	}
}

type bitbucket struct {
	client *bbapi.Client
}

func (a bitbucket) GetRepositories(owner string) chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		res, err := a.client.Repositories.ListForAccount(&bbapi.RepositoriesOptions{
			Owner: owner,
		})
		if err != nil {
			logrus.Fatal(err)
		}

		for _, repo := range res.Items {
			ch <- repo.Slug
		}
	}()

	return ch
}
func (a bitbucket) GetFiles(owner, repo string) chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		files, err := a.client.Repositories.Repository.ListFiles(&bbapi.RepositoryFilesOptions{
			Owner:    owner,
			RepoSlug: repo,
			Ref:      "HEAD",
			Path:     "docs",
		})
		if err != nil {
			if e, ok := err.(*bbapi.UnexpectedResponseStatusError); ok && strings.Contains(e.Status, "404") {
				return
			}
		}

		for _, file := range files {
			ch <- file.Path
		}
	}()

	return ch
}

func (a bitbucket) DownloadFile(owner, repo, file, localPath string) error {
	b, err := a.client.Repositories.Repository.GetFileContent(&bbapi.RepositoryFilesOptions{
		Owner:    owner,
		RepoSlug: repo,
		Ref:      "HEAD",
		Path:     file,
	})

	if err != nil {
		return err
	}
	f, err := os.Create(localPath)
	if err != nil {
		return err
	}
	f.Write(b)
	f.Sync()
	return f.Close()
}
