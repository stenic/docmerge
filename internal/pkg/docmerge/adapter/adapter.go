package adapter

type Adapter interface {
	GetRepositories(owner string) chan string
	GetFiles(owner, repo string) chan string
	DownloadFile(owner, repo, file, localPath string) error
}
