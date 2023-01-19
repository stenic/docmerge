package docmerge

import (
	"context"
	"io"
	"os"
	"path"
	"regexp"

	"github.com/google/go-github/v49/github"
	"github.com/sirupsen/logrus"
)

type Github struct {
	client  *github.Client
	context context.Context
	log     *logrus.Entry

	owner string
}

func (g Github) getAllRepos(opts *github.RepositoryListByOrgOptions) ([]*github.Repository, error) {
	if opts == nil {
		opts = &github.RepositoryListByOrgOptions{}
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := g.client.Repositories.ListByOrg(g.context, g.owner, opts)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allRepos, nil
}

func (g Github) downloadDocs(repo *github.Repository, targetDir string) error {
	tree, _, err := g.client.Git.GetTree(g.context, g.owner, *repo.Name, *repo.DefaultBranch, true)
	if err != nil {
		return err
	}
	re := regexp.MustCompile(`^docs\/`)
	os.RemoveAll(targetDir)
	for _, entry := range tree.Entries {
		repoPath := entry.GetPath()
		if re.Match([]byte(repoPath)) {
			localPath := path.Join(targetDir, entry.GetPath()[4:])
			g.log.Debugf("Fetching %s into %s", repoPath, localPath)
			if err := g.downloadFile(repo, repoPath, localPath); err != nil {
				g.log.Error(err)
			}
		}
	}

	return nil
}

func (g Github) downloadFile(repo *github.Repository, repoPath, localPath string) error {
	iorc, _, err := g.client.Repositories.DownloadContents(g.context, g.owner, repo.GetName(), repoPath, &github.RepositoryContentGetOptions{})
	if err != nil {
		return err
	}
	defer iorc.Close()

	if err := os.MkdirAll(path.Dir(localPath), 0755); err != nil {
		return err
	}

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
