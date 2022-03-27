package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func CurlGet(baseURL string, path string) []byte {
	baseURL = strings.Trim(baseURL, "/")
	req, err := http.NewRequest("GET", strings.TrimRight(baseURL, "/")+"/"+strings.TrimLeft(path, "/"), nil)

	if err != nil {
		log.Println(err.Error())
	}

	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println("Unable to retrieve news:", err.Error())
	}

	defer resp.Body.Close()

	if resp != nil && resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
		}

		//bodyString := string(bodyBytes)
		//log.Println(bodyString)

		return bodyBytes
	}

	return []byte{}
}
