package model

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Conf struct {
	Files      FileConf      `yaml:"files"`
	Statements StatementConf `yaml:"statements"`
}

type FileConf struct {
	Include FilePattern `yaml:"include"`
	Exclude FilePattern `yaml:"exclude"`
}

type FilePattern struct {
	Dirs     []string `yaml:"dirs"`
	Patterns []string `yaml:"patterns"`
}

type StatementConf struct {
	Exclude StatementPattern `yaml:"exclude"`
}

type StatementPattern struct {
	Patterns   []string `yaml:"patterns"`
	Annotation string   `yaml:"annotation"`
}

func LoadConfFromYAML(path string) (Conf, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Conf{}, err
	}

	var conf Conf
	if unmarshalErr := yaml.Unmarshal(data, &conf); unmarshalErr != nil {
		return Conf{}, unmarshalErr
	}
	return conf, nil
}
