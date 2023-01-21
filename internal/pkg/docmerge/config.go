package docmerge

type DMConfig struct {
	OutputDir string
	DocsDir   string

	GithubToken       string
	GithubOwner       string
	GithubTopicFilter string

	GitlabToken       string
	GitlabOwner       string
	GitlabTopicFilter string
}
