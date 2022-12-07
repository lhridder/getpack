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
)

func Get() error {
	mcversion := config.Global.Forge.Version

	err := os.Mkdir("forgeinstaller", os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create folder: %s", err)
	}

	err = os.Chdir("forgeinstaller")
	if err != nil {
		return fmt.Errorf("failed to go to folder: %s", err)
	}

	err = Install(mcversion)
	if err != nil {
		return err
	}

	zipname := fmt.Sprintf("%s.zip", mcversion)
	err = exec.Command("zip", "-qr", zipname, "./").Run()
	if err != nil {
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

	log.Println("Finished installing forge")

	return nil
}

func Install(mcversion string) error {
	url, version, err := GetURL(mcversion)
	if err != nil {
		return fmt.Errorf("failed to get url: %s", err)
	}

	err = util.Download(url, "forge-installer.jar")
	if err != nil {
		return fmt.Errorf("failed to download installer: %s", err)
	}

	output, err := exec.Command("java", "-jar", "forge-installer.jar", "--installServer").Output()
	if err != nil {
		log.Println(string(output))
		return fmt.Errorf("failed to run installer: %s", err)
	}

	err = os.Remove("forge-installer.jar")
	if err != nil {
		return fmt.Errorf("failed to remove installer: %s", err)
	}

	err = os.Remove("forge-installer.jar.log")
	if err != nil {
		return fmt.Errorf("failed to remove installer log: %s", err)
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
