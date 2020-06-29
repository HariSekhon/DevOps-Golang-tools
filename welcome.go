///bin/sh -c true; exec /usr/bin/env go run "$0" "$@"
//  vim:ts=4:sts=4:sw=4:noet
//
//  Author: Hari Sekhon
//  Date: 2020-04-24 14:14:44 +0100 (Fri, 24 Apr 2020)
//
//  https://github.com/harisekhon/devops-golang-tools
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
	"bytes"
	"flag"
	"fmt"
	//"log"
	// drop in replacement, with more levels and .SetLevel()
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path"
	"regexp"
	"strings"
	"syscall"
	"time"
)

const description = `
Prints a slick welcome message with last login time

Tested on Mac OS X and Linux
`

var prog = path.Base(os.Args[0])

// not compatible with logrus nor necessary, use Fprintf(os.Stderr, ...) instead
//var stderr = log.New(os.Stderr, "", 0)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\nusage: %s [options]\n\n", description, prog)
		flag.PrintDefaults()
		os.Exit(3)
	}
	var debug = flag.Bool("D", false, "Debug mode")
	var quick = flag.Bool("q", false, "Quick - print instantly without fancy scrolling effect, saves 2-3 seconds (you can also Control-C to make output complete instantly)")
	flag.Parse()
	if *debug || os.Getenv("DEBUG") != "" {
		log.SetLevel(log.DebugLevel)
		log.Debug("debug logging enabled")
	}
	msg := constructMsg()
	keyboardInterruptHandler(msg)
	// if we're being run in buffered 'go run', just print quickly without spinner
	matched, _ := regexp.MatchString("/go-build\\d+/[^/]+/exe/[^/]+$", os.Args[0])
	if *quick || os.Getenv("QUICK") != "" || matched {
		fmt.Println(msg)
	} else {
		printWithSpinner(msg)
	}
}

func titlecaseUser(user string) string {
	if user == "root" {
		user = strings.ToUpper(user)
	} else {
		matched, _ := regexp.MatchString("\\d$", user)
		if len(user) < 4 && matched {
			// probably not a real name
			// pass
		} else {
			user = strings.Title(user)
		}
	}
	return user
}

func constructMsg() string {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	var username string
	username = user.Username
	username = titlecaseUser(username)
	msg := fmt.Sprintf("Welcome %s - ", username)
	msgNoLastLoginInfo := "no last login information available!"
	/*
		last, err := os.Executable("last")
		if err != nil {
			msg += msgNoLastLoginInfo
			return msg
		}
	*/
	regexSkip := regexp.MustCompile("^(?:reboot|wtmp)|^\\s*$")
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command("last", "-100")
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err = cmd.Run()
	if err != nil {
		msg += msgNoLastLoginInfo
		msg += fmt.Sprintf(" ('last' command failed to execute: %s)", err)
		return msg
	}
	stdout, stderr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	if strings.TrimSpace(stderr) != "" {
		msg += msgNoLastLoginInfo
		msg += fmt.Sprintf(" ('last' stderr: %s)", stderr)
		return msg
	}
	lines := strings.Split(stdout, "\n")
	lastLine := ""
	for _, line := range lines {
		if regexSkip.MatchString(line) {
			continue
		}
		lastLine = line
		break
	}
	if lastLine != "" {
		msg += "last login was "
		regexLastUser := regexp.MustCompile("\\s+.*$")
		lastUser := regexLastUser.ReplaceAllString(lastLine, "")
		if lastUser == "root" {
			lastUser = "ROOT"
		}
		regexDate := regexp.MustCompile(".*(\\w{3}\\s+\\w{3}\\s+\\d+)")
		lastLine = regexDate.ReplaceAllString(lastLine, "$1")
		if lastUser == "ROOT" {
			msg += "ROOT"
		} else if strings.ToLower(lastUser) == strings.ToLower(username) {
			msg += "by you"
		} else {
			msg += fmt.Sprintf("by %s", lastUser)
		}
		msg += fmt.Sprintf(" => %s", lastLine)
	} else {
		msg += "no last login information available!"
	}
	return msg
}

func printWithSpinner(msg string) {
	if strings.TrimSpace(os.Getenv("QUICK")) != "" {
		fmt.Println(msg)
		return
	}
	stdout := bufio.NewWriter(os.Stdout)
	// many non-ASCII character sets we don't care about
	// unicode.Lower.R16 - 16-bit code-points for lowercase chars
	// unicode..
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz012345689@#$%^&*()"
	charlist := []rune(chars)
	sleepDuration, err := time.ParseDuration("0.0085s")
	if err != nil {
		log.Fatal(err)
	}
	for _, char := range msg {
		fmt.Printf(" ")
		j := 0
		for {
			var randomChar rune
			if j > 3 {
				randomChar = char
			} else {
				randomIndex := rand.Intn(len(charlist))
				randomChar = charlist[randomIndex]
			}
			fmt.Printf("\b%s", string(char))
			stdout.Flush()
			if char == randomChar {
				break
			}
			j++
			time.Sleep(sleepDuration)
		}
	}
	fmt.Println()
}

func keyboardInterruptHandler(msg string) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\r")
		fmt.Println(msg)
		os.Exit(0)
	}()
}
