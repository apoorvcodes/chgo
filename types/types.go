package types

type Config struct {
	AccessToken string `yaml:"accessToken"`
	UserIdent   string `yaml:"userIdent"`
	Locale      string `yaml:"locale"`
}

type Course struct {
	ID       string
	Title    string
	Language string
	Duration string
	URL      string
	Lessons  string
	Source   string
}
