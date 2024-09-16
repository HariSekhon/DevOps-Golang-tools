///bin/sh -c true; exec /usr/bin/env go run "$0" "$@"
//  vim:ts=4:sts=4:sw=4:noet
//
//  Author: Hari Sekhon
//  Date: 2024-09-16 02:22:33 +0200 (Mon, 16 Sep 2024)
//
//  https///github.com/HariSekhon/DevOps-Golang-tools
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
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

const description = `
On Linux parses the /proc/<pid>/maps to list all dyanmic so libraries that a program is using

The runtime equivalent of the classic Linux ldd command

Because sometimes the system pldd command gives results like this:

	pldd: cannot attach to process 32781: Operation not permitted
`

var prog = path.Base(os.Args[0])

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\nusage: %s <pid>\n\n", description, prog)
		flag.PrintDefaults()
		os.Exit(3)
	}
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		flag.Usage()
	}

	pid := args[0]

	if !isNumeric(pid) {
		fmt.Printf("Error: PID '%s' is not a valid number.\n", pid)
		os.Exit(1)
	}

	// Get the maximum possible PID from the system
	maxPid, err := getMaxPid()
	if err != nil {
		fmt.Printf("Error reading max PID: %v\n", err)
		os.Exit(1)
	}

	// Convert PID to integer and validate range
	pidNum, err := strconv.Atoi(pid)
	if err != nil || pidNum <= 0 || pidNum > maxPid {
		fmt.Printf("Error: PID '%s' is not in the valid range of 1 to %d.\n", pid, maxPid)
		os.Exit(1)
	}

	procMapsPath := fmt.Sprintf("/proc/%s/maps", pid)

	file, err := os.Open(procMapsPath)
	if err != nil {
		fmt.Printf("Failed to open %s: %v\n", procMapsPath, err)
		return
	}
	defer file.Close()

	// use a map to dedupe .so libraries
	soFiles := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		// last field contains the file path
		libPath := fields[len(fields)-1]

		// if the path contains in '.so' its a shared library
		if strings.HasSuffix(libPath, ".so") || strings.Contains(libPath, ".so.") {
			// get the absolute path of the shared library
			absPath, err := filepath.EvalSymlinks(libPath)
			if err != nil {
				absPath = libPath
			}

			// dedupe .so against map
			if !soFiles[absPath] {
				soFiles[absPath] = true
				fmt.Println(absPath)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading %s: %v\n", procMapsPath, err)
	}
}

func isNumeric(arg string) bool {
	_, err := strconv.Atoi(arg)
	return err == nil
}

// get the maximum possible PID from /proc/sys/kernel/pid_max
func getMaxPid() (int, error) {
	data, err := ioutil.ReadFile("/proc/sys/kernel/pid_max")
	if err != nil {
		return 0, err
	}
	maxPidStr := strings.TrimSpace(string(data))
	return strconv.Atoi(maxPidStr)
}
