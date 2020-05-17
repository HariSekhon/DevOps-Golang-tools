//  vim:ts=4:sts=4:sw=4:noet
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

// Returns the first HTTP(s) server argument to respond and serve its default page without error
//
// See also much more mature version find_active_server.py in DevOps Python tools - https://github.com/harisekhon/DevOps-Python-tools

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"net/http"
	"path"
	"regexp"
)

func main() {
	urls := os.Args[1:]
	if len(urls) < 1 {
		fmt.Printf("usage: %s <url> [<url> <url> ...]\n", path.Base(os.Args[0]))
		os.Exit(3)
	}

	results := make(chan string)

	http_prefix_regex := regexp.MustCompile("(?i)^https?://")

	for _, url := range urls {
		if !http_prefix_regex.MatchString(url) {
			url = "http://" + url
		}
		go get_url(url, results)
	}
	// print first result
	fmt.Println(<-results)
}

func get_url(url string, results chan string) {
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
