package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	gistsApiEndpoint = "https://api.github.com/gists"
	fileName         = fmt.Sprintf("stdgist-%d", time.Now().Unix())
)

const usage = `Standard Gist

Usage: stdgist [options]

stdgist reads from standard input and posts the given
input as a gist.

Options:
  --name, -n  Name for the file.
`

type gistReqBody struct {
	Public bool                         `json:"public"`
	Files  map[string]map[string]string `json:"files"`
}

type gistResBody struct {
	HtmlUrl string `json:"html_url"`
}

func init() {
	flag.StringVar(&fileName, "name", fileName, "")
	flag.StringVar(&fileName, "n", fileName, "")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usage)
	}
}

func main() {
	flag.Parse()

	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Check if stdin is from a terminal (i.e. no input to read).
	if stat.Mode()&os.ModeCharDevice != 0 {
		fmt.Fprintln(os.Stderr, "Content needs to be given to stdin to use")
		os.Exit(1)
	}

	var content bytes.Buffer
	_, err = io.Copy(&content, os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	body := gistReqBody{
		Public: true,
		Files: map[string]map[string]string{
			fileName: map[string]string{
				"content": content.String(),
			},
		},
	}

	content.Reset()
	encoder := json.NewEncoder(&content)
	err = encoder.Encode(body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	res, err := http.Post(gistsApiEndpoint, "application/json", &content)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer res.Body.Close()

	var resBody gistResBody
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&resBody)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if resBody.HtmlUrl != "" {
		fmt.Println(resBody.HtmlUrl)
	}
}
