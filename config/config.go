package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Debug   bool
	Target  string `yaml:"target"`
	Packs   map[int][]string
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

type Packs struct {
	Packs map[int][]string `yaml:"packs"`
}

var Global Config

func Load() error {
	configfile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		return fmt.Errorf("failed to read config.yml: %s", err)
	}

	err = yaml.Unmarshal(configfile, &Global)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config.yml: %s", err)
	}

	packsfile, err := ioutil.ReadFile("packs.yml")
	if err != nil {
		return fmt.Errorf("failed to read packs.yml: %s", err)
	}

	var packs Packs
	err = yaml.Unmarshal(packsfile, &packs)
	if err != nil {
		return fmt.Errorf("failed to unmarshal packs.yml: %s", err)
	}

	Global.Packs = packs.Packs
	return nil
}
