package model

type ProjectsMapping struct {
	Projects []Project `yaml:"projects"`
}
type Project struct {
	Name    string `yaml:"name"`
	Channel string `yaml:"channel"`
	ID      int    `yaml:"id"`
}