package docmerge

type DMConfig struct {
	OutputDir string
	DocsDir   string

	Endpoint    string
	Token       string
	Owner       string
	TopicFilter string
}
