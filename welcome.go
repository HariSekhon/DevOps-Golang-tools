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

const DESCRIPTION = `
Prints a slick welcome message with last login time

Tested on Mac OS X and Linux
`

var prog = path.Base(os.Args[0])

// not compatible with logrus
//var stderr = log.New(os.Stderr, "", 0)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\nusage: %s [options]\n\n", DESCRIPTION, prog)
		flag.PrintDefaults()
		os.Exit(3)
	}
	var debug = flag.Bool("debug", false, "Debug mode")
	var quick = flag.Bool("quick", false, "Print instantly without fancy scrolling effect, saves 2-3 seconds (you can also Control-C to make output complete instantly)")
	flag.Parse()
	if *debug || os.Getenv("DEBUG") != "" {
		log.SetLevel(log.DebugLevel)
		log.Debug("debug logging enabled")
	}
	msg := construct_msg()
	KeyboardInterruptHandler(msg)
	// if we're being run in buffered 'go run', just print quickly without spinner
	matched, _ := regexp.MatchString("/go-build\\d+/[^/]+/exe/[^/]+$", os.Args[0])
	if *quick || os.Getenv("QUICK") != "" || matched {
		fmt.Println(msg)
	} else {
		print_with_spinner(msg)
	}
}

func titlecase_user(user string) string {
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

func construct_msg() string {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	var username string
	username = user.Username
	username = titlecase_user(username)
	msg := fmt.Sprintf("Welcome %s - ", username)
	msg_no_last_login_info := "no last login information available!"
	/*
		last, err := os.Executable("last")
		if err != nil {
			msg += msg_no_last_login_info
			return msg
		}
	*/
	skip_regex := regexp.MustCompile("^(?:reboot|wtmp)|^\\s*$")
	var stdout_buf, stderr_buf bytes.Buffer
	cmd := exec.Command("last", "-100")
	cmd.Stdout = &stdout_buf
	cmd.Stderr = &stderr_buf
	err = cmd.Run()
	if err != nil {
		msg += msg_no_last_login_info
		msg += fmt.Sprintf(" ('last' command failed to execute: %s)", err)
		return msg
	}
	stdout, stderr := string(stdout_buf.Bytes()), string(stderr_buf.Bytes())
	if strings.TrimSpace(stderr) != "" {
		msg += msg_no_last_login_info
		msg += fmt.Sprintf(" ('last' stderr: %s)", stderr)
		return msg
	}
	lines := strings.Split(stdout, "\n")
	last_line := ""
	for _, line := range lines {
		if skip_regex.MatchString(line) {
			continue
		}
		last_line = line
		break
	}
	if last_line != "" {
		msg += "last login was "
		last_user_regex := regexp.MustCompile("\\s+.*$")
		last_user := last_user_regex.ReplaceAllString(last_line, "")
		if last_user == "root" {
			last_user = "ROOT"
		}
		date_regex := regexp.MustCompile(".*(\\w{3}\\s+\\w{3}\\s+\\d+)")
		last_line = date_regex.ReplaceAllString(last_line, "$1")
		if last_user == "ROOT" {
			msg += "ROOT"
		} else if strings.ToLower(last_user) == strings.ToLower(username) {
			msg += "by you"
		} else {
			msg += fmt.Sprintf("by %s", last_user)
		}
		msg += fmt.Sprintf(" => %s", last_line)
	} else {
		msg += "no last login information available!"
	}
	return msg
}

func print_with_spinner(msg string) {
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
	sleep_duration, err := time.ParseDuration("0.0085s")
	if err != nil {
		log.Fatal(err)
	}
	for _, char := range msg {
		fmt.Printf(" ")
		j := 0
		for {
			var random_char rune
			if j > 3 {
				random_char = char
			} else {
				random_index := rand.Intn(len(charlist))
				random_char = charlist[random_index]
			}
			fmt.Printf("\b%s", string(char))
			stdout.Flush()
			if char == random_char {
				break
			}
			j += 1
			time.Sleep(sleep_duration)
		}
	}
	fmt.Println()
}

func KeyboardInterruptHandler(msg string) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\r")
		fmt.Println(msg)
		os.Exit(0)
	}()
}
