package forge

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type forgepromotions struct {
	Promos map[string]string `json:"promos"`
}

const (
	versionlist = "https://files.minecraftforge.net/net/minecraftforge/forge/promotions_slim.json"
	base        = "https://maven.minecraftforge.net/net/minecraftforge/forge/"
)

func GetURL(mcversion string) (string, string, error) {
	res, err := http.Get(versionlist)
	if err != nil {
		return "", "", fmt.Errorf("failed to get json: %s", err)
	}

	if res.StatusCode != 200 {
		return "", "", fmt.Errorf("got non 200 status code: %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read body: %s", err)
	}

	var promos forgepromotions
	err = json.Unmarshal(body, &promos)
	if err != nil {
		return "", "", fmt.Errorf("failed to unmarshal body: %s", err)
	}

	forgeversion := promos.Promos[fmt.Sprintf("%s-recommended", mcversion)]

	version := fmt.Sprintf("%s-%s", mcversion, forgeversion)

	log.Printf("Found forge version %s", version)

	return fmt.Sprintf("%s%s/forge-%s-installer.jar", base, version, version), version, nil
}
