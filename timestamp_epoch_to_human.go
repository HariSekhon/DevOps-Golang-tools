///bin/sh -c true; exec /usr/bin/env go run "$0" "$@"
//  vim:ts=4:sts=4:sw=4:noet
//
//  Author: Hari Sekhon
//  Date: 2023-06-14 22:44:00 +0100 (Wed, 14 Jun 2023)
//
//  https://github.com/HariSekhon/DevOps-Golang-tools
//
//  License: see accompanying Hari Sekhon LICENSE file
//
//  If you're using my code you're welcome to connect with me on LinkedIn and optionally send me feedback to help steer this or other code I publish
//
//  https://www.linkedin.com/in/HariSekhon
//

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const description = `
Converts epoch timestamp from logs such as External Secrets pod logs into a human readable format

Works like a standard unix filter program - takes either a filename or standard input and prints to standard output replacing the epoch with human readable time
`

var prog = path.Base(os.Args[0])

func readline() string {
	in := bufio.NewReader(os.Stdin)
	line, err := in.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	line = strings.TrimSpace(line)
	return line
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\nusage: %s [<file>]\n\n", description, prog)
		flag.PrintDefaults()
		os.Exit(3)
	}
	//var debug = flag.Bool("D", false, "Debug mode")
	flag.Parse()
	//if *debug || os.Getenv("DEBUG") != "" {
	//    log.SetLevel(log.DebugLevel)
	//    log.Debug("debug logging enabled")
	//}

	var scanner *bufio.Scanner
	if len(os.Args) > 1 {
		filename := os.Args[1]
		filehandle, err := os.Open(filename)
		if err != nil {
			//fmt.Fprintf(os.Stderr, "error: %s\n", err)
			//os.Exit(1)
			log.Fatal(err)
		}
		defer filehandle.Close()
		scanner = bufio.NewScanner(filehandle)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	re := regexp.MustCompile(`\b\d{10}(?:\.\d{1,7})?\b`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindAllString(line, -1)
		for _, match := range matches {
			epoch, err := strconv.ParseFloat(match, 64)
			if err == nil {
				seconds := int64(epoch)
				milliseconds := int64((epoch - float64(seconds)) * 1000)

				t := time.Unix(seconds, milliseconds*int64(time.Millisecond))
				//convertedTime := t.Format("2006-01-02 15:04:05.000")
				convertedTime := t.Format(time.RFC3339)

				line = strings.Replace(line, match, convertedTime, 1)
			}
		}
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
