package docmerge

import (
	"log"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

func (d DocMerge) runGitlab(owner string) error {
	glClient, err := gitlab.NewClient(d.cfg.GitlabToken)
	if err != nil {
		log.Fatal(err)
	}

	// List all projects
	logrus.Infof("Fetching project for %s", owner)
	gl := Gitlab{
		client: glClient,
	}
	// projects, err := gl.getProjects(owner)
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	for project := range gl.getProjects(owner) {
		logrus.Info(project.NameWithNamespace)
		files, err := gl.getFiles(project.PathWithNamespace, "docs")
		if err != nil {
			logrus.Error(err)
			continue
		}
		targetDir := path.Join(d.cfg.OutputDir, project.PathWithNamespace)
		os.RemoveAll(targetDir)

		for _, file := range files {

			b, _, err := glClient.RepositoryFiles.GetRawFile(project.PathWithNamespace, file, &gitlab.GetRawFileOptions{})
			if err != nil {
				logrus.Error(err)
				continue
			}
			localPath := path.Join(targetDir, file[4:])
			if err := os.MkdirAll(path.Dir(localPath), 0755); err != nil {
				return err
			}
			f, err := os.Create(localPath)
			if err != nil {
				logrus.Error(err)
				continue
			}
			f.Write(b)
			f.Sync()
			f.Close()
		}
	}

	return nil
}

type Gitlab struct {
	client *gitlab.Client
}

func (g Gitlab) getProjects(owner string) chan *gitlab.Project {
	ch := make(chan *gitlab.Project)

	opts := &gitlab.ListGroupProjectsOptions{
		Simple:           gitlab.Bool(true),
		IncludeSubGroups: gitlab.Bool(true),
	}

	go func() {
		for {
			logrus.Debug("Fetching more repostories")
			projects, resp, err := g.client.Groups.ListGroupProjects(owner, opts)
			if err != nil {
				logrus.Error(err)
			}

			for _, project := range projects {
				ch <- project
			}

			// Exit the loop when we've seen all pages.
			if resp.NextPage == 0 {
				close(ch)
				return
			}

			// Update the page number to get the next page.
			opts.Page = resp.NextPage
		}
	}()

	return ch
}

func (g Gitlab) getFiles(pid interface{}, prefix string) ([]string, error) {
	var allFiles []string
	topts := &gitlab.ListTreeOptions{
		Path:        gitlab.String(prefix),
		Recursive:   gitlab.Bool(true),
		ListOptions: gitlab.ListOptions{},
	}
	for {
		treeNodes, tresp, err := g.client.Repositories.ListTree(pid, topts)
		if err != nil {
			return allFiles, err
		}
		for _, node := range treeNodes {
			if node.Type == "blob" {
				allFiles = append(allFiles, node.Path)
			}
		}

		// Exit the loop when we've seen all pages.
		if tresp.NextPage == 0 {
			break
		}

		// Update the page number to get the next page.
		topts.ListOptions.Page = tresp.NextPage
	}

	return allFiles, nil
}
