package purpur

import (
	"fmt"
	"getpack/config"
	"getpack/utils"
	"log"
	"os"
)

func Get() error {
	version := config.Global.Purpur.Version
	filename := fmt.Sprintf("%s.jar", version)

	url, err := getURL()
	if err != nil {
		return fmt.Errorf("failed to get url: %s", err)
	}

	err = os.Mkdir("purpurinstaller", os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create folder: %s", err)
	}

	err = os.Chdir("purpurinstaller")
	if err != nil {
		return fmt.Errorf("failed to go to folder: %s", err)
	}

	err = utils.Download(url, filename)
	if err != nil {
		return fmt.Errorf("failed to download jar: %s", err)
	}

	folder := config.Global.Target + "purpur/"
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

	log.Println("Finished installing purpur")

	return nil
}
