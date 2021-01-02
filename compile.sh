#!/usr/bin/env bash
#  vim:ts=4:sts=4:sw=4:et
#
#  Author: Hari Sekhon
#  Date: 2021-01-01 22:07:48 +0000 (Fri, 01 Jan 2021)
#
#  https://github.com/HariSekhon/go-tools
#
#  License: see accompanying Hari Sekhon LICENSE file
#
#  If you're using my code you're welcome to connect with me on LinkedIn and optionally send me feedback to help steer this or other code I publish
#
#  https://www.linkedin.com/in/HariSekhon
#

set -euo pipefail
[ -n "${DEBUG:-}" ] && set -x
srcdir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# shellcheck disable=SC1090
. "$srcdir/bash-tools/lib/utils.sh"

# shellcheck disable=SC2034,SC2154
usage_description="
Compiles Go programs in the adjacent directory

Handles different versions of Golang to work with older versions

- go get's dependencies from go.mod manually if the Golang version doesn't support 'go mod'
"

# used by usage() in lib/utils.sh
# shellcheck disable=SC2034
usage_args=""

help_usage "$@"

cd "$srcdir"

# must not be a directory with go.mod
export GOPATH=~/go
export GOBIN="${GOBIN:-$PWD/bin}"

echo
echo "go env:"
echo
go env
echo
echo "GOPATH = ${GOPATH:-}"
echo "GOBIN  = ${GOBIN:-}"
echo
which go
if is_mac; then
    readlink(){
        greadlink "$@"
    }
fi
ls -l "$(readlink -f "$(which go)")"
echo
go version
echo
echo

# for older versions of Go that don't support 'go mod'
if ! go help mod &>/dev/null; then
    awk '/require/{gsub("v", "", $3); print $2}' go.mod |
    xargs -L 1 go get
    echo
else
    # golang 1.11.13 seems to fail, try to download deps explicitly
    echo "go mod download"
    go mod download
    echo
fi

opts=()
# race detector doesn't work with musl on Alpine
if ! grep -qi Alpine /etc/*release &>/dev/null; then
    opts+=(-race)
fi

for x in *.go; do
    #echo "go build -race -o bin/ $x";
    #go build -race -o bin/ "$x" ||
    echo "go install ${opts[*]} $x";
    go install "${opts[@]}" "$x"
    echo
done
echo "Golang compile succeeded"
