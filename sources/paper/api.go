package paper

import (
	"encoding/json"
	"fmt"
	"getpack/config"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type paperVersion struct {
	Builds []int `json:"builds"`
}

const base = "https://api.papermc.io/v2/projects/paper/versions/"

func getURL() (string, error) {
	version := config.Global.Paper.Version
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

	var paper paperVersion
	err = json.Unmarshal(body, &paper)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal body: %s", err)
	}

	build := strconv.Itoa(paper.Builds[len(paper.Builds)-1])

	log.Printf("Found paper version %s-%s", version, build)
	url := fmt.Sprintf("%s%s/builds/%s/downloads/paper-%s-%s.jar", base, version, build, version, build)

	return url, nil
}
