package docmerge

type DMConfig struct {
	OutputDir string

	GithubToken       string
	GithubOwner       string
	GithubTopicFilter string

	GitlabToken string
	GitlabOwner string
}
