package http_util

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// 	url := ""
// ""
func WWWFromPost(url, data string) (error, []byte) {
	payload := strings.NewReader(data)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err, nil
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return err, nil
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err, nil
	}
	return nil, body
}
