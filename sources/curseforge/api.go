package curseforge

import (
	"encoding/json"
	"fmt"
	"getpack/config"
	"getpack/util"
	"regexp"
	"strconv"
	"strings"
)

type file struct {
	DisplayName      string   `json:"displayName"`
	FileName         string   `json:"fileName"`
	ServerPackFileID int      `json:"serverPackFileId"`
	IsServerPack     bool     `json:"isServerPack"`
	GameVersions     []string `json:"gameVersions"`
}

type serverpack struct {
	Data string `json:"data"`
}

type Cursepack struct {
	Data struct {
		ID     int    `json:"id"`
		GameID int    `json:"gameId"`
		Name   string `json:"name"`
		Slug   string `json:"slug"`
		Links  struct {
			WebsiteURL string `json:"websiteUrl"`
		}
		Summary       string `json:"summary"`
		Status        int    `json:"status"`
		DownloadCount int    `json:"downloadCount"`
		Logo          struct {
			ThumnbailURL string `json:"thumnbailUrl"`
			URL          string `json:"url"`
		}
		LatestFiles []file `json:"latestFiles"`
	}
	URL     string
	Version string
}

type cursefiles struct {
	Data []file `json:"data"`
}

const base = "https://api.curseforge.com/v1/"

func Get(packid int) (*Cursepack, error) {
	url := fmt.Sprintf("%smods/%s", base, strconv.Itoa(packid))

	header := map[string]string{"x-api-key": config.Global.Curse.APIkey}
	body, err := util.Fetch(url, header)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %s", err)
	}

	var pack *Cursepack
	err = json.Unmarshal(body, &pack)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %s", err)
	}

	latestfile := pack.Data.LatestFiles[len(pack.Data.LatestFiles)-1]
	if latestfile.ServerPackFileID == 0 {
		url = fmt.Sprintf("%smods/%s/files", base, strconv.Itoa(packid))

		body, err = util.Fetch(url, header)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch: %s", err)
		}

		var files cursefiles
		err = json.Unmarshal(body, &files)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal json: %s", err)
		}

		for _, latestfile = range files.Data {
			if latestfile.ServerPackFileID != 0 {
				break
			}
		}
	}

	serverpackurl, err := getServerPackURL(packid, latestfile.ServerPackFileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get server pack url: %s", err)
	}

	pack.URL = serverpackurl
	parts := strings.Split(strings.ReplaceAll(latestfile.DisplayName, ".zip", ""), " ")
	for _, part := range parts {
		if regexp.MustCompile(`\d`).MatchString(part) {
			//TODO integrate into regex
			if !strings.HasPrefix(part, "[") {
				pack.Version = part
			}
		}
	}

	return pack, nil
}

func getServerPackURL(packid int, serverpackid int) (string, error) {
	url := fmt.Sprintf("%smods/%s/files/%s/download-url", base, strconv.Itoa(packid), strconv.Itoa(serverpackid))

	header := map[string]string{"x-api-key": config.Global.Curse.APIkey}
	body, err := util.Fetch(url, header)
	if err != nil {
		return "", fmt.Errorf("failed to fetch: %s", err)
	}

	var serverpack serverpack
	err = json.Unmarshal(body, &serverpack)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal json: %s", err)
	}

	return serverpack.Data, nil
}
