package fabric

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type installer struct {
	URL     string `json:"url"`
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
}

var base = "https://meta.fabricmc.net/v2/versions/installer"

func GetInstaller() (string, error) {
	res, err := http.Get(base)
	if err != nil {
		return "", fmt.Errorf("failed to get json: %s", err)
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("got non 200 status code: %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %s", err)
	}

	var installers []installer
	err = json.Unmarshal(body, &installers)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal body: %s", err)
	}

	for _, installer := range installers {
		if installer.Stable {
			log.Printf("Found fabric installer version %s", installer.Version)
			return installer.URL, nil
		}
	}

	return "", fmt.Errorf("no stable version found")
}
