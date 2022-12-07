package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Target  string           `yaml:"target"`
	Packs   map[int][]string `yaml:"packs"`
	Enabled struct {
		Curse   bool `yaml:"curse"`
		Technic bool `yaml:"technic"`
		FTB     bool `yaml:"ftb"`
		Paper   bool `yaml:"paper"`
		Purpur  bool `yaml:"purpur"`
		Spigot  bool `yaml:"spigot"`
		Forge   bool `yaml:"forge"`
		Fabric  bool `yaml:"fabric"`
		Bedrock bool `yaml:"bedrock"`
	}
	Curse struct {
		APIkey   string `yaml:"apikey"`
		Modpacks []int  `yaml:"modpacks"`
	}
	Technic struct {
		Modpacks []string `yaml:"modpacks"`
	}
	FTB struct {
		Modpacks []int `yaml:"modpacks"`
	}
	Paper struct {
		Version string `yaml:"version"`
	}
	Purpur struct {
		Version string `yaml:"version"`
	}
	Spigot struct {
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
