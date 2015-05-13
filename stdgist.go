package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	gistsApiEndpoint = "https://api.github.com/gists"
)

type gistReqBody struct {
	Public bool                         `json:"public"`
	Files  map[string]map[string]string `json:"files"`
}

type gistResBody struct {
	HtmlUrl string `json:"html_url"`
}

func main() {
	var content string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}

	var data bytes.Buffer
	fileName := fmt.Sprintf("stdgist-%s", strconv.FormatInt(time.Now().Unix(), 10))
	body := gistReqBody{
		Public: true,
		Files: map[string]map[string]string{
			fileName: map[string]string{
				"content": content,
			},
		},
	}

	encoder := json.NewEncoder(&data)
	err := encoder.Encode(body)
	if err != nil {
		panic(err)
	}

	res, err := http.Post(gistsApiEndpoint, "application/json", &data)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var resBody gistResBody
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&resBody)
	if err != nil {
		panic(err)
	}

	if resBody.HtmlUrl != "" {
		fmt.Println(resBody.HtmlUrl)
	}
}
