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
	Base        = "https://maven.minecraftforge.net/net/minecraftforge/forge/"
)

func GetURL(mcversion string, forgeversion string) (string, string, error) {
	if forgeversion != "" {
		version := fmt.Sprintf("%s-%s", mcversion, forgeversion)
		return fmt.Sprintf("%s%s/forge-%s-installer.jar", Base, version, version), version, nil
	}

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

	fversion := promos.Promos[fmt.Sprintf("%s-latest", mcversion)]
	version := fmt.Sprintf("%s-%s", mcversion, fversion)

	log.Printf("Found forge version %s", version)

	return fmt.Sprintf("%s%s/forge-%s-installer.jar", Base, version, version), version, nil
}
