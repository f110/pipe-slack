package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	usage = `Usage: pipe-slack source incoming_webhook_url 
`
)

type slackMessage struct {
	Text string `json:"text"`
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}
	if !strings.HasPrefix(os.Args[2], "https") {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	var source io.Reader
	if os.Args[1] == "-" {
		source = os.Stdin
	} else {
		f, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not open %s\n", os.Args[1])
			os.Exit(1)
		}
		source = f
	}

	s := bufio.NewScanner(source)
	toSlack := &slackMessage{}
	slackUrl := os.Args[2]
	for s.Scan() {
		msg := s.Text()
		toSlack.Text = msg

		payload, err := json.Marshal(toSlack)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fail marshal payload: %s\n", err)
			os.Exit(1)
		}
		res, err := http.PostForm(slackUrl, url.Values{"payload": {string(payload)}})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fail post message: %s\n", err)
		}
		res.Body.Close()
	}
}
