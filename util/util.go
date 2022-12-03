package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func Download(url string, filename string) error {
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get jar: %s", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create jar: %s", err)
	}

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return fmt.Errorf("failed to write jar: %s", err)
	}

	err = res.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to close body: %s", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("failed to close jar: %s", err)
	}

	return nil
}
