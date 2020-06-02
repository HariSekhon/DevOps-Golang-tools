#!/usr/bin/env bash
#  vim:ts=4:sts=4:sw=4:et
#
#  Author: Hari Sekhon
#  Date: 2020-06-02 12:37:27 +0100 (Tue, 02 Jun 2020)
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

cd "$srcdir/.."

# shellcheck disable=SC1091
. ./tests/utils.sh

run ./uniq2 < tests/data/uniq.txt

run ./uniq2 tests/data/uniq.txt

run ./uniq2 -ignore-case tests/data/uniq.txt

run ./uniq2 -ignore-whitespace tests/data/uniq.txt

run ./uniq2 -ignore-whitespace -ignore-case tests/data/uniq.txt
