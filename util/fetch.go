package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const useragent = "github.com/lhridder/getpack"

func Fetch(url string, key map[string]string) ([]byte, error) {
	res, err := Get(url, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get: %s", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %s", err)
	}

	err = res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close body: %s", err)
	}

	return body, nil
}

func Get(url string, key map[string]string) (*http.Response, error) {
	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to format request: %s", err)
	}

	req.Header.Set("User-Agent", useragent)
	if key != nil {
		for k, v := range key {
			req.Header.Set(k, v)
		}
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get json: %s", err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("got non 200 status code: %s", res.Status)
	}

	return res, nil
}
