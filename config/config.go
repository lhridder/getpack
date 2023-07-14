package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Debug   bool             `yaml:"debug"`
	Target  string           `yaml:"target"`
	Deploy  bool             `yaml:"deploy"`
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
	Discord struct {
		Enabled bool   `yaml:"enabled"`
		Token   string `yaml:"token"`
		Channel string `yaml:"channel"`
	}
	Curse struct {
		APIkey   string         `yaml:"apikey"`
		Modpacks map[string]int `yaml:"modpacks"`
	}
	Technic struct {
		Modpacks []string `yaml:"modpacks"`
	}
	FTB struct {
		Modpacks map[string]int `yaml:"modpacks"`
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
		Versions []string `yaml:"versions"`
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
