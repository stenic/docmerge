package docmerge

import (
	"context"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"

	"github.com/google/go-github/v49/github"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"
)

func (d DocMerge) runGithub(owner string) error {
	tempDir := "/tmp/docmerge-cache"
	tc := &http.Client{
		Transport: &oauth2.Transport{
			Base: httpcache.NewTransport(diskcache.New(tempDir)),
			Source: oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: d.cfg.GithubToken},
			),
		},
	}

	client := github.NewClient(tc)
	gh := Github{
		client:  client,
		context: d.ctx,
		owner:   owner,
		log:     logrus.NewEntry(logrus.New()),
	}

	repos, err := gh.getAllRepos(nil)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		gh.log = logrus.WithField("repo", repo.GetName())
		if d.cfg.GithubTopicFilter != "" && !slices.Contains(repo.Topics, d.cfg.GithubTopicFilter) {
			gh.log.Debug("Skipping: does not have label 'internal-docs'")
			continue
		}

		gh.log.Infof("Fetching %s", repo.GetName())
		if err := gh.downloadDocs(repo, path.Join(d.cfg.OutputDir, *repo.Name)); err != nil {
			gh.log.Error(err)
		}
	}

	return nil
}

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
