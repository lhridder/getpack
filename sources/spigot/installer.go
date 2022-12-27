package spigot

import (
	"fmt"
	"getpack/config"
	"getpack/util"
	"log"
	"os"
	"os/exec"
	"time"
)

const buildtools = "https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar"

func Get() error {
	start := time.Now()

	err := os.Mkdir("spigotinstaller", os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create folder: %s", err)
	}

	err = os.Chdir("spigotinstaller")
	if err != nil {
		return fmt.Errorf("failed to go to folder: %s", err)
	}

	err = util.Download(buildtools, "buildtools.jar")
	if err != nil {
		return fmt.Errorf("failed to download buildtools: %s", err)
	}

	if config.Global.Debug {
		log.Printf("Downloading buildtools took %.2fs", time.Now().Sub(start).Seconds())
	}

	version := config.Global.Spigot.Version
	startbuildtools := time.Now()
	output, err := exec.Command("java", "-jar", "buildtools.jar", "--rev", version).Output()
	if err != nil {
		log.Println(string(output))
		return fmt.Errorf("failed to run buildtools: %s", err)
	}

	if config.Global.Debug {
		log.Printf("Running buildtools took %.2fs", time.Now().Sub(startbuildtools).Seconds())
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

	log.Printf("Finished installing spigot in %.2fs", time.Now().Sub(start).Seconds())

	return nil
}
