package http_util

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func JsonRPC(method, url, data string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
