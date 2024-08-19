///bin/sh -c true; exec /usr/bin/env go run "$0" "$@"
//  vim:ts=4:sts=4:sw=4:noet
//
//  Author: Hari Sekhon
//  Date: 2024-08-19 08:51:51 +0200 (Mon, 19 Aug 2024)
//  Original Date: 2012-12-29 10:53:23 +0000 (Sat, 29 Dec 2012)
//
//  https///github.com/HariSekhon/DevOps-Golang-tools
//
//  License: see accompanying Hari Sekhon LICENSE file
//
//  If you're using my code you're welcome to connect with me on LinkedIn and optionally send me feedback to help steer this or other code I publish
//
//  https://www.linkedin.com/in/HariSekhon
//

/*
* Filter program to print net line additions/removals from diff / patch files or stdin"
*
* This is a rewrite of a perl version I used for over a decade from
*
*	https://github.com/HariSekhon/DevOps-Perl-tools
*
* which was itself was a rewrite of a shell version I used for years before that in my extensive
* and borderline ridiculously over developed but immensely cool bashrc
* (nearly 4500 lines at the time I've split this off)
* (6500 with aliases files + an additional 21,000 lines of supporting scripts)
* that eventually evolved into this huge repo:
*
*	https://github.com/HariSekhon/DevOps-Bash-tools
*
*/

// TODO: use counters so that I don't discount 2 removals for 1 addition etc

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

// Command-line options
var (
	additionsOnly     = flag.Bool("a", false, "Show only additions")
	removalsOnly      = flag.Bool("r", false, "Show only removals")
	blocks            = flag.Bool("b", false, "Show changes in blocks of additions first and then removals")
	ignoreCase        = flag.Bool("i", false, "Ignore case in comparisons")
	ignoreWhitespace  = flag.Bool("w", false, "Ignore whitespace in comparisons")
	addPrefix         string
	removePrefix      string
	ignoreWhitespaceRe = regexp.MustCompile(`\s+`)
)

// Utility functions for case and whitespace handling
func transformations(str string) string {
	if *ignoreCase {
		str = strings.ToLower(str)
	}
	if *ignoreWhitespace {
		str = ignoreWhitespaceRe.ReplaceAllString(str, "")
	}
	return str
}

// Process a diff file and collect additions and removals
func diffnet(scanner *bufio.Scanner) {
	additions := make(map[int]string)
	removals := make(map[int]string)

	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		if addPrefix == "" && strings.HasPrefix(line, "+") {
			addPrefix = "+"
		}
		if removePrefix == "" && strings.HasPrefix(line, "-") {
			removePrefix = "-"
		}

		if strings.HasPrefix(line, "+") {
			additions[lineNum] = line[1:]
		} else if strings.HasPrefix(line, "-") {
			removals[lineNum] = line[1:]
		}
	}

	// Output based on flags
	if *blocks || *additionsOnly || *removalsOnly {
		if !*removalsOnly {
			printChanges(additions, removals, addPrefix)
		}
		if !*additionsOnly {
			printChanges(removals, additions, removePrefix)
		}
	} else {
		printCombinedChanges(additions, removals)
	}
}

// Print additions or removals based on flag settings
func printChanges(main, opposite map[int]string, prefix string) {
	keys := sortedKeys(main)
	for _, i := range keys {
		// Print if not found in the opposite map (removals for additions, etc.)
		if !containsTransform(opposite, main[i]) {
			fmt.Println(prefix + main[i])
		}
	}
}

// Print combined changes when both additions and removals are processed together
func printCombinedChanges(additions, removals map[int]string) {
	allKeys := append(sortedKeys(additions), sortedKeys(removals)...)
	sort.Ints(allKeys)
	seen := make(map[int]bool)
	for _, i := range allKeys {
		if seen[i] {
			continue
		}
		seen[i] = true
		if add, ok := additions[i]; ok {
			if !containsTransform(removals, add) {
				fmt.Println(addPrefix + add)
			}
		} else if rem, ok := removals[i]; ok {
			if !containsTransform(additions, rem) {
				fmt.Println(removePrefix + rem)
			}
		}
	}
}

// Check if a transformed string exists in the map
func containsTransform(data map[int]string, value string) bool {
	transformedValue := transformations(value)
	for _, v := range data {
		if transformations(v) == transformedValue {
			return true
		}
	}
	return false
}

// Get sorted keys of a map (line numbers)
func sortedKeys(data map[int]string) []int {
	keys := make([]int, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func main() {
	flag.Parse()

	if *additionsOnly && *removalsOnly {
		log.Fatal("Error: --additions-only and --removals-only are mutually exclusive!")
	}

	if flag.NArg() == 0 {
		// Read from stdin
		scanner := bufio.NewScanner(os.Stdin)
		diffnet(scanner)
	} else {
		// Read from provided files
		for _, filename := range flag.Args() {
			file, err := os.Open(filename)
			if err != nil {
				log.Fatalf("Error opening file %s: %v", filename, err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			diffnet(scanner)
		}
	}
}
