package main

import (
	"fmt"
	"getpack/config"
	"getpack/sources/bedrock"
	"getpack/sources/fabric"
	"getpack/sources/forge"
	"getpack/sources/paper"
	"getpack/sources/purpur"
	"getpack/sources/spigot"
	"log"
	"os"
)

var (
	cfg config.Config
	dir string
)

func main() {
	err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %s", err)
		return
	}
	cfg = config.Global
	dir, _ = os.Getwd()

	err = getVersions()
	if err != nil {
		log.Println(err)
	}
}

func getVersions() error {
	var err error

	if cfg.Enabled.Paper {
		err = paper.Get()
		if err != nil {
			return fmt.Errorf("failed to get paper: %s", err)
		}

		err = os.Chdir(dir)
		if err != nil {
			return fmt.Errorf("failed to change directories to home: %s", err)
		}

		err = os.RemoveAll("paperinstaller")
		if err != nil {
			return fmt.Errorf("failed to delete paperinstaller folder: %s", err)
		}

	}

	if cfg.Enabled.Purpur {
		err = purpur.Get()
		if err != nil {
			return fmt.Errorf("failed to get purpur: %s", err)
		}

		err = os.Chdir(dir)
		if err != nil {
			return fmt.Errorf("failed to change directories to home: %s", err)
		}

		err = os.RemoveAll("purpurinstaller")
		if err != nil {
			return fmt.Errorf("failed to delete purpurinstaller folder: %s", err)
		}
	}

	if cfg.Enabled.Spigot {
		err = spigot.Get()
		if err != nil {
			log.Printf("failed to get spigot: %s", err)
		}

		err = os.Chdir(dir)
		if err != nil {
			return fmt.Errorf("failed to change directories to home: %s", err)
		}

		err = os.RemoveAll("spigotinstaller")
		if err != nil {
			return fmt.Errorf("failed to delete spigotinstaller folder: %s", err)
		}
	}

	if cfg.Enabled.Forge {
		err = forge.Get()
		if err != nil {
			return fmt.Errorf("failed to get forge: %s", err)
		}

		err = os.Chdir(dir)
		if err != nil {
			return fmt.Errorf("failed to change directories to home: %s", err)
		}

		err = os.RemoveAll("forgeinstaller")
		if err != nil {
			return fmt.Errorf("failed to delete forgeinstaller folder: %s", err)
		}
	}

	if cfg.Enabled.Fabric {
		url, err := fabric.GetInstaller()
		if err == nil {
			err = fabric.Install(url)
			if err != nil {
				return fmt.Errorf("failed to install fabric: %s", err)
			}
		} else {
			return fmt.Errorf("failed to fabric installer: %s", err)
		}

		err = os.Chdir(dir)
		if err != nil {
			return fmt.Errorf("failed to change directories to home: %s", err)
		}

		err = os.RemoveAll("fabricinstaller")
		if err != nil {
			return fmt.Errorf("failed to delete fabricinstaller folder: %s", err)
		}
	}

	if cfg.Enabled.Bedrock {
		err = bedrock.Get()
		if err != nil {
			return fmt.Errorf("failed to get bedrock: %s", err)
		}

		err = os.Chdir(dir)
		if err != nil {
			return fmt.Errorf("failed to change directories to home: %s", err)
		}

		err = os.RemoveAll("bedrockinstaller")
		if err != nil {
			return fmt.Errorf("failed to delete bedrockinstaller folder: %s", err)
		}
	}

	return nil
}
