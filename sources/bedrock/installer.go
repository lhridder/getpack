package bedrock

import (
	"fmt"
	"getpack/config"
	"getpack/util"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func Get() error {
	req, err := http.NewRequest("GET", "https://www.minecraft.net/en-us/download/server/bedrock", nil)
	if err != nil {
		return fmt.Errorf("failed to format request: %s", err)
	}
	req.Header.Set("Authority", "www.minecraft.net")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.107 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept-Language", "nl-NL,nl;q=0.9,en-US;q=0.8,en;q=0.7")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get: %s", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %s", err)
	}

	var link string
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		if strings.Contains(line, "azureedge") && strings.Contains(line, "linux") {
			for _, item := range strings.Split(line, string('"')) {
				match, _ := regexp.MatchString("(http|https)://[a-zA-Z0-9./?=_-]*", item)
				if match {
					link = item
				}
			}
		}
	}

	parts := strings.Split(link, "/")

	log.Printf("Found bedrock version %s", parts[len(parts)-1])

	err = os.Mkdir("bedrockinstaller", os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create folder: %s", err)
	}

	err = os.Chdir("bedrockinstaller")
	if err != nil {
		return fmt.Errorf("failed to go to folder: %s", err)
	}

	filename := "bedrock-server.zip"

	err = util.Download(link, filename)
	if err != nil {
		return fmt.Errorf("failed to download bedrock zip: %s", err)
	}

	err = exec.Command("unzip", filename).Run()
	if err != nil {
		return fmt.Errorf("failed to unzip bedrock server: %s", err)
	}

	folder := config.Global.Target + "bedrock/"
	_, err = os.Stat(folder)
	if os.IsNotExist(err) {
		err = os.Mkdir(folder, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create target folder: %s", err)
		}
	}

	err = os.Rename(filename, folder+filename)
	if err != nil {
		return fmt.Errorf("failed to move zip to target folder: %s", err)
	}

	err = os.Rename("bedrock_server", folder+"bedrock_server")
	if err != nil {
		return fmt.Errorf("failed to move executable to target folder: %s", err)
	}

	log.Println("Finished installing bedrock")

	return nil
}
