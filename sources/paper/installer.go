package paper

import (
	"fmt"
	"getpack/config"
	"getpack/util"
	"log"
	"os"
	"time"
)

func Get() error {
	start := time.Now()

	url, err := getURL()
	if err != nil {
		return fmt.Errorf("failed to get url: %s", err)
	}

	filename := fmt.Sprintf("%s.jar", config.Global.Paper.Version)

	err = os.Mkdir("paperinstaller", os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create folder: %s", err)
	}

	err = os.Chdir("paperinstaller")
	if err != nil {
		return fmt.Errorf("failed to go to folder: %s", err)
	}

	err = util.Download(url, filename)
	if err != nil {
		return fmt.Errorf("failed to download jar: %s", err)
	}

	folder := config.Global.Target + "paper/"
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

	log.Printf("Finished installing paper in %.2fs", time.Now().Sub(start).Seconds())

	return nil
}
