package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Target  string `json:"target"`
	Enabled struct {
		Curse   bool `json:"curse"`
		Paper   bool `json:"paper"`
		Purpur  bool `json:"purpur"`
		Forge   bool `json:"forge"`
		Fabric  bool `json:"fabric"`
		Bedrock bool `json:"bedrock"`
	}
	Curse struct {
		APIkey   string `yaml:"apikey"`
		Modpacks []int  `yaml:"modpacks"`
	}
	Paper struct {
		Version string `yaml:"version"`
	}
	Purpur struct {
		Version string `yaml:"version"`
	}
	Forge struct {
		Version string `yaml:"version"`
	}
	Fabric struct {
		Versions []string `yaml:"versions"`
	}
}

var Global Config

func Load() error {
	file, err := ioutil.ReadFile("config.yml")
	if err != nil {
		return err
	}
	return yaml.Unmarshal(file, &Global)
}
