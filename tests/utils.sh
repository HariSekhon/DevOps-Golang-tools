#!/usr/bin/env bash
#  vim:ts=4:sts=4:sw=4:et
#
#  Author: Hari Sekhon
#  Date: 2015-05-25 01:38:24 +0100 (Mon, 25 May 2015)
#
#  https://github.com/harisekhon/devops-golang-tools
#
#  License: see accompanying Hari Sekhon LICENSE file
#
#  If you're using my code you're welcome to connect with me on LinkedIn
#  and optionally send me feedback to help improve or steer this or other code I publish
#
#  https://www.linkedin.com/in/harisekhon
#

set -eu
[ -n "${DEBUG:-}" ] && set -x
srcdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# shellcheck disable=SC1090
. "$srcdir/excluded.sh"

# shellcheck disable=SC1090
. "$srcdir/../bash-tools/lib/utils.sh"

export COMPOSE_PROJECT_NAME="go-tools"

# shellcheck disable=SC1090
. "$srcdir/excluded.sh"

#export GOBIN="$srcdir/../bin"
bin="bin"

build(){
    local target="$1"
    if [ -f "$bin/$target" ]; then
        echo "$bin/$target detected"
        echo
    else
        echo "$bin/$target not detected, building now"
        echo
        make
        echo
    fi
}
