package ftb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type ftbtarget struct {
	Version string `json:"version"`
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Updated int    `json:"updated"`
}

type ftbversion struct {
	Id      int         `json:"id"`
	Name    string      `json:"name"`
	Type    string      `json:"type"`
	Updated int         `json:"updated"`
	Targets []ftbtarget `json:"targets"`
}

type ftbart struct {
	URL    string `json:"url"`
	Type   string `json:"type"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type FTBpack struct {
	Id            int          `json:"id"`
	Name          string       `json:"name"`
	Versions      []ftbversion `json:"versions"`
	Art           []ftbart     `json:"art"`
	Description   string       `json:"description"`
	Version       ftbversion
	ServerPackURL string
}

const base = "https://api.modpacks.ch/public/modpack/"

func Get(packid int) (*FTBpack, error) {
	res, err := http.Get(base + strconv.Itoa(packid))
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

	var pack *FTBpack
	err = json.Unmarshal(body, &pack)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal body: %s", err)
	}

	pack.Version = pack.Versions[len(pack.Versions)-1]
	parts := strings.Split(pack.Version.Name, " ")
	for _, part := range parts {
		if regexp.MustCompile(`\d`).MatchString(part) {
			pack.Version.Name = part
		}
	}

	pack.ServerPackURL = fmt.Sprintf("%s%d/%d/server/linux", base, packid, pack.Version.Id)

	return pack, nil
}
