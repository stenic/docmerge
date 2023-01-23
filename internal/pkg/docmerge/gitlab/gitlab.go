package gitlab

import (
	"context"
	"docmerge/internal/pkg/docmerge/adapter"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	glapi "github.com/xanzy/go-gitlab"
)

func New(token string) adapter.Adapter {
	client, err := glapi.NewClient(token)
	if err != nil {
		logrus.Fatal(err)
	}

	return gitlab{
		client:  client,
		context: context.Background(),
	}
}

type gitlab struct {
	client  *glapi.Client
	context context.Context
}

func (a gitlab) GetRepositories(owner string) chan string {
	ch := make(chan string)

	opts := &glapi.ListGroupProjectsOptions{
		Simple:           glapi.Bool(true),
		IncludeSubGroups: glapi.Bool(true),
	}

	go func() {
		defer close(ch)
		for {
			projects, resp, err := a.client.Groups.ListGroupProjects(owner, opts)
			if err != nil {
				logrus.Error(err)
				continue
			}

			for _, project := range projects {
				ch <- project.PathWithNamespace
			}

			// Exit the loop when we've seen all pages.
			if resp.NextPage == 0 {
				return
			}

			// Update the page number to get the next page.
			opts.Page = resp.NextPage
		}
	}()

	return ch
}
func (a gitlab) GetFiles(owner, repo string) chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		topts := &glapi.ListTreeOptions{
			// Path:        glapi.String(prefix),
			Recursive:   glapi.Bool(true),
			ListOptions: glapi.ListOptions{},
		}
		for {
			treeNodes, tresp, err := a.client.Repositories.ListTree(repo, topts)
			if err != nil {
				continue
			}
			for _, node := range treeNodes {
				if node.Type == "blob" {
					ch <- node.Path
				}
			}

			// Exit the loop when we've seen all pages.
			if tresp.NextPage == 0 {
				break
			}

			// Update the page number to get the next page.
			topts.ListOptions.Page = tresp.NextPage
		}
	}()

	return ch
}
func (a gitlab) DownloadFile(owner, repo, file, localPath string) error {
	b, _, err := a.client.RepositoryFiles.GetRawFile(repo, file, &glapi.GetRawFileOptions{})
	if err != nil {
		return err
	}
	if err := os.MkdirAll(path.Dir(localPath), 0755); err != nil {
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
