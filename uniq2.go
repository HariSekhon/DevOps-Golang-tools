///bin/sh -c true; exec /usr/bin/env go run "$0" "$@"
//  vim:ts=4:sts=4:sw=4:noet
//
//  Author: Hari Sekhon
//  Date: 2020-06-01 20:11:40 +0100 (Mon, 01 Jun 2020)
//  Original Date: 2015-02-07 16:06:33 +0000 (Sat, 07 Feb 2015)
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
	"bufio"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"strings"
)

const description = `
Filter program to print only non-repeated lines in input - unlike the unix command 'uniq' lines do not have to be adjacent, this is order preserving compared to 'sort | uniq'. I rustled this up quickly after needing to parse unique missing modules for building but maintaining order as some modules depend on others being built first

Works as a standard unix filter program taking either standard input or files supplied as arguments

Since this must maintain unique lines in memory for comparison, do not use this on very large files/inputs

Port of Perl program uniq_order_preserved.pl written early 2015 which can be found in the adjacent DevOps Perl tools repo

Tested on Mac OS X and Linux
`

var prog = path.Base(os.Args[0])

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\nusage: %s [file1] [file2] ...\n\n", description, prog)
		flag.PrintDefaults()
		os.Exit(3)
	}
	var debug = flag.Bool("D", false, "Debug mode")
	var ignoreCase = flag.Bool("c", false, "Ignore case in comparisons")
	var ignoreWhitespace = flag.Bool("w", false, "Ignore whitespace in comparisons")
	flag.Parse()
	if *debug || os.Getenv("DEBUG") != "" {
		log.SetLevel(log.DebugLevel)
		log.Debug("debug logging enabled")
	}

	if len(flag.Args()) < 1 {
		printUniq("-", *ignoreCase, *ignoreWhitespace)
	} else {
		for _, filename := range flag.Args() {
			printUniq(filename, *ignoreCase, *ignoreWhitespace)
		}
	}
}

var uniqMap = make(map[string]bool)

// returns True if line is uniq, False if it's been seen before
func uniq(line string, ignoreCase bool, ignoreWhitespace bool) bool {
	if ignoreWhitespace {
		line = strings.TrimSpace(line)
	}
	if ignoreCase {
		line = strings.ToLower(line)
	}
	_, exists := uniqMap[line]
	if exists {
		return false
	}
	uniqMap[line] = true
	return true
}

func printUniq(filename string, ignoreCase bool, ignoreWhitespace bool) {
	if filename == "-" {
		stdin := bufio.NewReader(os.Stdin)
		processLines(stdin, ignoreCase, ignoreWhitespace)
		return
	}

	filehandle, err := os.Open(filename)
	if err != nil {
		log.Error(err)
		return
	}
	defer filehandle.Close()
	processLines(filehandle, ignoreCase, ignoreWhitespace)
}

func processLines(reader io.Reader, ignoreCase bool, ignoreWhitespace bool) {
	var line string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line = scanner.Text()
		if uniq(line, ignoreCase, ignoreWhitespace) {
			fmt.Println(line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
