package purpur

import (
	"encoding/json"
	"fmt"
	"getpack/config"
	"io/ioutil"
	"log"
	"net/http"
)

type purpurVersion struct {
	Builds struct {
		Latest string `json:"latest"`
	}
}

const base = "https://api.purpurmc.org/v2/purpur/"

func getURL() (string, error) {
	version := config.Global.Purpur.Version

	res, err := http.Get(base + version)
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

	var purpur purpurVersion
	err = json.Unmarshal(body, &purpur)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal body: %s", err)
	}

	build := purpur.Builds.Latest
	log.Printf("Found purpur version %s-%s", version, build)

	return fmt.Sprintf("%s%s/%s/download", base, version, build), nil
}
