package ftb

import (
	"fmt"
	"getpack/config"
	"getpack/util"
	"getpack/utils"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func (pack *FTBpack) Install(packname string) error {
	start := time.Now()

	instructions := config.Global.Packs[pack.Id]
	if len(instructions) == 0 {
		return fmt.Errorf("instructions for pack are too short/not available")
	}

	err := os.Mkdir(packname, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create folder: %s", err)
	}

	err = os.Chdir(packname)
	if err != nil {
		return fmt.Errorf("failed to go to folder: %s", err)
	}

	err = util.Download(pack.ServerPackURL, "serverinstall")
	if err != nil {
		return fmt.Errorf("failed to download: %s", err)
	}

	err = os.Chmod("./serverinstall", 0755)
	if err != nil {
		return fmt.Errorf("failed to make server installer executable")
	}

	output, err := exec.Command("./serverinstall", strconv.Itoa(pack.Id), strconv.Itoa(pack.Version.Id), "--auto", "--nojava").Output()
	if err != nil {
		log.Println(string(output))
		return fmt.Errorf("failed to run installer: %s", err)
	}

	err = os.Remove("serverinstall")
	if err != nil {
		return fmt.Errorf("failed to delete installer: %s", err)
	}

	err = utils.Instructions(instructions)
	if err != nil {
		return fmt.Errorf("failed to run instructions: %s", err)
	}

	zipname := fmt.Sprintf("%s.zip", pack.Version.Name)
	output, err = exec.Command("zip", "-qr", zipname, "./").Output()
	if err != nil {
		log.Println(string(output))
		return fmt.Errorf("failed to zip pack: %s", err)
	}

	folder := fmt.Sprintf("%smodpacks/ftb/%s/", config.Global.Target, packname)
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

	log.Printf("Finished installing %s version %s in %.2fs", pack.Name, pack.Version.Name, time.Now().Sub(start).Seconds())

	return nil
}
