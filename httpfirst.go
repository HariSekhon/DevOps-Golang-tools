///bin/sh -c true; exec /usr/bin/env go run "$0" "$@"
//  vim:ts=4:sts=4:sw=4:noet
//  args: duckduckgo.com google.com
//
//  Author: Hari Sekhon
//  Date: 2020-05-17 16:07:00 +0100 (Sun, 17 May 2020)
//
//  https://github.com/harisekhon/go-tools
//
//  License: see accompanying Hari Sekhon LICENSE file
//
//  If you're using my code you're welcome to connect with me on LinkedIn and optionally send me feedback to help steer this or other code I publish
//
//  https://www.linkedin.com/in/harisekhon
//

package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
)

const description = `
Returns the first HTTP(s) server argument to respond and serve its default page without error

See also much more mature version find_active_server.py in DevOps Python tools - https://github.com/harisekhon/DevOps-Python-tools
`

var prog = path.Base(os.Args[0])

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\nusage: %s <url> [<url> <url> ...]\n\n", description, prog)
		flag.PrintDefaults()
		os.Exit(3)
	}
	var debug = flag.Bool("D", false, "Debug mode")
	flag.Parse()
	if *debug || os.Getenv("DEBUG") != "" {
		log.SetLevel(log.DebugLevel)
		log.Debug("debug logging enabled")
	}

	//urls := os.Args[1:]
	urls := flag.Args()
	if len(urls) < 1 {
		//fmt.Printf("usage: %s <url> [<url> <url> ...]\n", path.Base(os.Args[0]))
		//os.Exit(3)
		flag.Usage()
	}
	for _, url := range urls {
		matched, _ := regexp.MatchString("^-", url)
		if matched {
			flag.Usage()
		}
	}

	results := make(chan string)

	httpPrefixRegex := regexp.MustCompile("(?i)^https?://")

	for _, url := range urls {
		if !httpPrefixRegex.MatchString(url) {
			url = "http://" + url
		}
		go getURL(url, results)
	}
	// print first result
	// will hang if none succeed - add timeout and more professional handlings like my find_active_server.py program
	fmt.Println(<-results)
}

func getURL(url string, results chan string) {
	res, err := http.Get(url)
	if err != nil {
		// ignore
		//panic(err)
		return
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		// ignore
		//panic(err)
		return
	}
	results <- url
}
