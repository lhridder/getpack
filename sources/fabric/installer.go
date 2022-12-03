package fabric

import (
	"fmt"
	"getpack/config"
	"getpack/util"
	"log"
	"os"
	"os/exec"
)

func Install(url string) error {
	err := os.Mkdir("fabricinstaller", os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create folder: %s", err)
	}

	err = os.Chdir("fabricinstaller")
	if err != nil {
		return fmt.Errorf("failed to go to folder: %s", err)
	}

	err = util.Download(url, "installer.jar")
	if err != nil {
		return fmt.Errorf("failed to download installer: %s", err)
	}

	for _, mcversion := range config.Global.Fabric.Versions {
		err := os.Mkdir(mcversion, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create mcversion folder: %s", err)
		}

		err = os.Chdir(mcversion)
		if err != nil {
			return fmt.Errorf("failed to go to mcversion folder: %s", err)
		}

		output, err := exec.Command("java", "-jar", "../installer.jar", "server", "-mcversion", mcversion, "-downloadMinecraft").Output()
		if err != nil {
			log.Println(output)
			return fmt.Errorf("failed to run installer for %s: %s", mcversion, err)
		}

		err = os.Rename("server.jar", "vanilla.jar")
		if err != nil {
			return fmt.Errorf("failed to rename server jar: %s", err)
		}

		err = os.Rename("fabric-server-launch.jar", "server.jar")
		if err != nil {
			return fmt.Errorf("failed to rename server jar: %s", err)
		}

		fabricfile, err := os.Create("fabric-server-launcher.properties")
		if err != nil {
			return fmt.Errorf("failed to create fabricfile: %s", err)
		}

		_, err = fmt.Fprintln(fabricfile, "serverJar=vanilla.jar")
		if err != nil {
			return fmt.Errorf("failed to write text to fabricfile: %s", err)
		}

		err = fabricfile.Close()
		if err != nil {
			return fmt.Errorf("failed to close fabricfile: %s", err)
		}

		zipname := fmt.Sprintf("%s.zip", mcversion)
		err = exec.Command("zip", "-qr", zipname, "./").Run()
		if err != nil {
			return fmt.Errorf("failed to zip installer for %s: %s", mcversion, err)
		}

		folder := config.Global.Target + "fabric/"
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

		log.Printf("Finished installing fabric %s", mcversion)
	}

	log.Println("Finished installing fabric versions")

	return nil
}
