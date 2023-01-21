package docmerge

import (
	"log"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/exp/slices"
)

func (d DocMerge) runGitlab(owner string) error {
	glClient, err := gitlab.NewClient(d.cfg.GitlabToken)
	if err != nil {
		log.Fatal(err)
	}

	gl := Gitlab{
		client: glClient,
		log:    logrus.NewEntry(logrus.New()),
	}

	logrus.Infof("Fetching project for %s", owner)
	for project := range gl.getProjects(owner) {
		if d.cfg.GitlabTopicFilter != "" && !slices.Contains(project.Topics, d.cfg.GitlabTopicFilter) {
			gl.log.Debugf("Skipping: does not have label '%s'", d.cfg.GitlabTopicFilter)
			continue
		}

		gl.log.Infof("Fetching %s", project.PathWithNamespace)
		files, err := gl.getFiles(project.PathWithNamespace, d.cfg.DocsDir)
		if err != nil {
			gl.log.Error(err)
			continue
		}
		targetDir := path.Join(d.cfg.OutputDir, project.PathWithNamespace)
		if err := os.RemoveAll(targetDir); err != nil {
			gl.log.Error(err)
		}

		for _, file := range files {
			localPath := path.Join(targetDir, file[len(d.cfg.DocsDir):])
			if err := gl.downloadFile(*project, file, localPath); err != nil {
				gl.log.Error(err)
			}
		}
	}
	logrus.Infof("Completed project for %s", owner)

	return nil
}

type Gitlab struct {
	client *gitlab.Client
	log    *logrus.Entry
}

func (g Gitlab) getProjects(owner string) chan *gitlab.Project {
	ch := make(chan *gitlab.Project)

	opts := &gitlab.ListGroupProjectsOptions{
		Simple:           gitlab.Bool(true),
		IncludeSubGroups: gitlab.Bool(true),
	}

	go func() {
		for {
			g.log.Debug("Fetching more repostories")
			projects, resp, err := g.client.Groups.ListGroupProjects(owner, opts)
			if err != nil {
				g.log.Error(err)
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

func (g Gitlab) downloadFile(project gitlab.Project, file, localPath string) error {
	b, _, err := g.client.RepositoryFiles.GetRawFile(project.PathWithNamespace, file, &gitlab.GetRawFileOptions{})
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
