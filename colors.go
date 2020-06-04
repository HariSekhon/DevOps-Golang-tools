///bin/sh -c true; exec /usr/bin/env go run "$0" "$@"
//  vim:ts=4:sts=4:sw=4:noet
//
//  Author: Hari Sekhon
//  Date: 2020-06-04 15:11:25 +0100 (Thu, 04 Jun 2020)
//  Original Date: 2014-06-07 22:17:09 +0100 (Sat, 07 Jun 2014)
//  Ported from Perl version in DevOps Perl tools repo
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
	//"log"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

const description = `
Program to show all the ASCII terminal code Foreground/Background color combinations in a terminal to make it easy to pick for writing fancy programs

Ported from an original Perl version from 2014 in the DevOps Perl tools repo: https://github.com/harisekhon/devops-perl-tools

Tested on Mac OS X and Linux
`

var prog = path.Base(os.Args[0])

// effects 4 = underline, 5 = blink, look ugly - added only in verbose mode
var effects = []int{0, 1}

const text = "hari"
const length = len(text) + 2

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\nusage: %s [options]\n\n", description, prog)
		flag.PrintDefaults()
		os.Exit(3)
	}
	var debug = flag.Bool("D", false, "Debug mode")
	var verbose = flag.Bool("v", false, "Verbose mode (print underlines and blinking effects too)")
	flag.Parse()
	if *debug || os.Getenv("DEBUG") != "" {
		log.SetLevel(log.DebugLevel)
		log.Debug("debug logging enabled")
	}

	if *verbose {
		effects = []int{0, 1, 4, 5}
	}
	fmt.Println()
	fmt.Println(`ASCII Terminal Codes Color Key:

EF  = Effect [1 = bold, 4 = underline, 5 = blink (only shown in verbose mode)]
TXT = Foreground text color
BG  = Background solid color
`)

	headerFormatString := fmt.Sprintf("%5s BG %%-%ds", "", length)
	fmt.Printf(headerFormatString, "none")
	for bg := 40; bg <= 47; bg++ {
		formatString := fmt.Sprintf(" %%%ddm  ", length-1)
		fmt.Printf(formatString, bg)
	}
	fmt.Printf("\n%5s\n", "EF;TXT")
	printLine(0)
	printLine(1)
	for txtcode := 30; txtcode <= 37; txtcode++ {
		printLine(txtcode)
	}
	fmt.Println()
}

func printLine(txtcode int) {
	var effectTxt string
	for _, effect := range effects {
		if effect == 0 {
			effectTxt = fmt.Sprintf("%d", txtcode)
		} else {
			effectTxt = fmt.Sprintf("%d;%d", effect, txtcode)
		}
		fmt.Printf(" %4sm ", effectTxt)
		fmt.Printf("\033[0m\033[%dm  %s  \033[0m  ", txtcode, text)
		for bg := 40; bg <= 47; bg++ {
			fmt.Printf("\033[%d;%dm\033[%dm  %s  \033[0m ", effect, txtcode, bg, text)
		}
		fmt.Println()
		if txtcode == 0 || txtcode == 1 {
			break
		}
	}
}
