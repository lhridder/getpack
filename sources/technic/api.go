package technic

import (
	"encoding/json"
	"fmt"
	"getpack/util"
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
	body, err := util.Fetch(base+packname+suffix, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get json: %s", err)
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
