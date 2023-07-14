package forge

import (
	"fmt"
	"getpack/config"
	"getpack/util"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func Get() error {
	start := time.Now()

	err := os.Mkdir("forgeinstaller", os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create folder: %s", err)
	}

	err = os.Chdir("forgeinstaller")
	if err != nil {
		return fmt.Errorf("failed to go to folder: %s", err)
	}

	for _, mcversion := range config.Global.Forge.Versions {
		log.Println("")

		err := os.Mkdir(mcversion, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create mcversion folder: %s", err)
		}

		err = os.Chdir(mcversion)
		if err != nil {
			return fmt.Errorf("failed to go to mcversion folder: %s", err)
		}

		err = Install(mcversion)
		if err != nil {
			return err
		}

		zipname := fmt.Sprintf("%s.zip", mcversion)
		output, err := exec.Command("zip", "-qr", zipname, "./").Output()
		if err != nil {
			log.Println(string(output))
			return fmt.Errorf("failed to zip installer for %s: %s", mcversion, err)
		}

		folder := config.Global.Target + "forge/"
		_, err = os.Stat(folder)
		if os.IsNotExist(err) {
			err = os.Mkdir(folder, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create target folder: %s", err)
			}
		}

		err = os.Rename(zipname, folder+zipname)
		if err != nil {
			return fmt.Errorf("failed to move zip to target folder: %s", err)
		}

		err = os.Chdir("..")
		if err != nil {
			return fmt.Errorf("failed to go back a directory: %s", err)
		}

		log.Printf("Finished installing forge %s", mcversion)
	}

	log.Printf("Finished installing forge versions in %.2fs", time.Now().Sub(start).Seconds())

	return nil
}

func Install(mcversion string) error {
	start := time.Now()

	url, version, err := GetURL(mcversion)
	if err != nil {
		return fmt.Errorf("failed to get url: %s", err)
	}

	err = util.Download(url, "forge-installer.jar")
	if err != nil {
		return fmt.Errorf("failed to download installer: %s", err)
	}

	if config.Global.Debug {
		log.Printf("Downloading forge installer took %.2fs", time.Now().Sub(start).Seconds())
	}

	install := time.Now()
	output, err := exec.Command("java", "-jar", "forge-installer.jar", "--installServer").Output()
	if err != nil {
		log.Println(string(output))
		return fmt.Errorf("failed to run installer: %s", err)
	}

	if config.Global.Debug {
		log.Printf("Running forge installer took %.2fs", time.Now().Sub(install).Seconds())
	}

	err = os.Remove("forge-installer.jar")
	if err != nil {
		return fmt.Errorf("failed to remove installer: %s", err)
	}

	_, err = os.Stat("forge-installer.jar.log")
	if os.IsNotExist(err) {
		err = os.Remove("installer.log")
		if err != nil {
			return fmt.Errorf("failed to remove installer log: %s", err)
		}
	} else {
		err = os.Remove("forge-installer.jar.log")
		if err != nil {
			return fmt.Errorf("failed to remove installer log: %s", err)
		}
	}

	major, err := strconv.ParseInt(strings.Split(mcversion, ".")[1], 10, 0)
	if err != nil {
		return fmt.Errorf("failed to parse major version: %s", err)
	}

	if major < 17 {
		oldname := fmt.Sprintf("forge-%s.jar", version)
		err = os.Rename(oldname, "server.jar")
		if err != nil {
			return fmt.Errorf("failed to rename server jar: %s", err)
		}
	} else {
		err = os.Remove("run.bat")
		if err != nil {
			return fmt.Errorf("failed to remove bat file: %s", err)
		}
	}

	return nil
}
