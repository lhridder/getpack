package spigot

import (
	"fmt"
	"getpack/config"
	"getpack/utils"
	"log"
	"os"
	"os/exec"
)

const buildtools = "https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar"

func Get() error {
	err := os.Mkdir("spigotinstaller", os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create folder: %s", err)
	}

	err = os.Chdir("spigotinstaller")
	if err != nil {
		return fmt.Errorf("failed to go to folder: %s", err)
	}

	err = utils.Download(buildtools, "buildtools.jar")
	if err != nil {
		return fmt.Errorf("failed to download buildtools: %s", err)
	}

	version := config.Global.Spigot.Version
	output, err := exec.Command("java", "-jar", "buildtools.jar", "--rev", version).Output()
	if err != nil {
		log.Println(string(output))
		return fmt.Errorf("failed to run buildtools: %s", err)
	}

	filename := fmt.Sprintf("spigot-%s.jar", version)

	folder := config.Global.Target + "spigot/"
	_, err = os.Stat(folder)
	if os.IsNotExist(err) {
		err = os.Mkdir(folder, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create target folder: %s", err)
		}
	}

	err = os.Rename(filename, folder+filename)
	if err != nil {
		return fmt.Errorf("failed to move jar to target folder: %s", err)
	}

	log.Println("Finished installing spigot")

	return nil
}
