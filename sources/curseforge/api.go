package curseforge

import (
	"encoding/json"
	"fmt"
	"getpack/config"
	"io/ioutil"
	"net/http"
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
	request, err := http.NewRequest("GET", fmt.Sprintf("%smods/%s", base, strconv.Itoa(packid)), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to format request: %s", err)
	}
	request.Header.Set("x-api-key", config.Global.Curse.APIkey)

	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to request pack json: %s", err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("received non 200 status code: %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	var pack *Cursepack
	err = json.Unmarshal(body, &pack)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %s", err)
	}

	latestfile := pack.Data.LatestFiles[len(pack.Data.LatestFiles)-1]
	if latestfile.ServerPackFileID == 0 {
		request2, err := http.NewRequest("GET", fmt.Sprintf("%smods/%s/files", base, strconv.Itoa(packid)), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to format request: %s", err)
		}
		request2.Header.Set("x-api-key", config.Global.Curse.APIkey)

		res2, err := client.Do(request2)
		if err != nil {
			return nil, fmt.Errorf("failed to request pack files json: %s", err)
		}

		if res2.StatusCode != 200 {
			return nil, fmt.Errorf("received non 200 status code: %s", res2.Status)
		}

		body2, err := ioutil.ReadAll(res2.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %s", err)
		}

		var files cursefiles
		err = json.Unmarshal(body2, &files)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal json: %s", err)
		}

		for _, latestfile = range files.Data {
			if latestfile.ServerPackFileID != 0 {
				break
			}
		}
	}

	url, err := getServerPackURL(packid, latestfile.ServerPackFileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get server pack url: %s", err)
	}

	pack.URL = url
	parts := strings.Split(strings.ReplaceAll(latestfile.DisplayName, ".zip", ""), " ")
	pack.Version = parts[len(parts)-1]

	return pack, nil
}

func getServerPackURL(packid int, serverpackid int) (string, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%smods/%s/files/%s/download-url", base, strconv.Itoa(packid), strconv.Itoa(serverpackid)), nil)
	if err != nil {
		return "", fmt.Errorf("failed to format request: %s", err)
	}
	request.Header.Set("x-api-key", config.Global.Curse.APIkey)

	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to request: %s", err)
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("received non 200 status code: %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %s", err)
	}

	var serverpack serverpack
	err = json.Unmarshal(body, &serverpack)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal json: %s", err)
	}

	return serverpack.Data, nil
}
