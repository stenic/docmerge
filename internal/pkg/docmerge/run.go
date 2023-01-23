package docmerge

import (
	"context"
	"docmerge/internal/pkg/docmerge/adapter"
	"docmerge/internal/pkg/docmerge/github"
	"docmerge/internal/pkg/docmerge/gitlab"
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/sirupsen/logrus"
)

type DocMerge struct {
	cfg DMConfig
	ctx context.Context
}

func (d DocMerge) Run(cfg DMConfig, endpoint string) error {
	d.ctx = context.Background()
	d.cfg = cfg
	d.cfg.DocsDir = "docs"

	logrus.SetLevel(logrus.DebugLevel)

	var a adapter.Adapter
	switch endpoint {
	case "github":
		a = github.New(cfg.Token)
	case "gitlab":
		a = gitlab.New(cfg.Token)
	default:
		logrus.Fatal("Endpoint not supported")
	}

	re := regexp.MustCompile(fmt.Sprintf(`^%s\/`, d.cfg.DocsDir))
	for repo := range a.GetRepositories(cfg.Owner) {
		log := logrus.WithField("repo", repo)
		log.Infof("Fetching %s", repo)
		targetDir := path.Join(d.cfg.OutputDir, repo)
		os.RemoveAll(targetDir)
		for file := range a.GetFiles(cfg.Owner, repo) {
			if !re.MatchString(file) {
				log.Tracef("Skipping %s", file)
				continue
			}
			log.Debugf("Downloading %s", file)
			localPath := path.Join(targetDir, file[len(d.cfg.DocsDir):])
			if err := os.MkdirAll(path.Dir(localPath), 0755); err != nil {
				return err
			}
			if err := a.DownloadFile(cfg.Owner, repo, file, localPath); err != nil {
				logrus.Error(err)
			}
		}
	}
	return nil
}
