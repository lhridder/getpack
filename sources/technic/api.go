package technic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type TechnicPack struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	DisplayName   string `json:"displayName"`
	PlatformURL   string `json:"platformUrl"`
	Description   string `json:"description"`
	Version       string `json:"version"`
	ServerPackURL string `json:"serverPackUrl"`
	Logo          struct {
		URL string `json:"url"`
	}
}

const (
	base   = "https://api.technicpack.net/modpack/"
	suffix = "?build=<build>"
)

func Get(packname string) (*TechnicPack, error) {
	res, err := http.Get(base + packname + suffix)
	if err != nil {
		return nil, fmt.Errorf("failed to get json: %s", err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("got non 200 status code: %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %s", err)
	}

	var pack *TechnicPack
	err = json.Unmarshal(body, &pack)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal body: %s", err)
	}

	parts := strings.Split(pack.Version, " ")
	for _, part := range parts {
		if regexp.MustCompile(`\d`).MatchString(part) {
			pack.Version = part
		}
	}

	return pack, nil
}
