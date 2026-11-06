#!/bin/sh
# Writes, into the directory named by $1, what the image owes the authors of
# the Alpine packages inside it: THIRD_PARTY_NOTICES, which lists every
# installed package with its license and says where the source of the GPL and
# LGPL ones is, and licenses/, the text of every license the list names.
#
# The list comes from apk's own database, so it is the image's packages and
# not a copy kept by hand. The texts come from Alpine's spdx-licenses-text
# package, which the Dockerfile installs for this script and removes after
# it. A license with no text there fails the build: the image never names a
# license it does not carry.
set -eu

out=$1
installed=${APK_INSTALLED:-/lib/apk/db/installed}
texts=${SPDX_TEXTS:-/usr/share/spdx/text}

mkdir -p "$out/licenses"

# One line per package: name, version, license, upstream, origin, commit.
packages=$(mktemp)
trap 'rm -f "$packages"' EXIT
awk 'BEGIN { RS = ""; FS = "\n"; OFS = "\t" }
{
	split("", field)
	for (i = 1; i <= NF; i++) field[substr($i, 1, 1)] = substr($i, 3)
	if (field["P"] != "spdx-licenses-text")
		print field["P"], field["V"], field["L"], field["U"], field["o"], field["c"]
}' "$installed" | sort >"$packages"
[ -s "$packages" ] || {
	echo "no packages in $installed" >&2
	exit 1
}

{
	cat <<'TEXT'
The cma-sandbox-fetch and cma-sandbox-worker scripts in this image are part
of https://github.com/anthropics/anthropic-cli (the sandboxes/ directory)
and are licensed under the Apache License 2.0, licenses/Apache-2.0.txt. The
ant binary is the Claude CLI from that repository. ant/ holds its license
and, under third_party/, those of Go and of every module compiled into it.

The Alpine Linux packages in this image are listed below, exactly as Alpine
built them. The text of each license named is in licenses/, under the
identifier used here.

PACKAGE  VERSION  LICENSE  UPSTREAM
TEXT
	awk -F '\t' '{ print $1 "  " $2 "  " $3 "  " $4 }' "$packages"
	cat <<'TEXT'

SOURCE CODE

The packages below are under the GNU General Public License or the GNU
Lesser General Public License. Their complete source is published with
every release of this image, in the same registry and under the same tag,
as us-docker.pkg.dev/anthropic-public-registry/sandboxes/cma-sandbox-source:
for each package, Alpine's build recipe and patches and the upstream
archives the recipe names. Each line names the commit of Alpine's package
tree, https://gitlab.alpinelinux.org/alpine/aports, that the package was
built from; the recipe is the directory named after the package's origin.

PACKAGE  VERSION  ORIGIN  APORTS COMMIT  UPSTREAM
TEXT
	awk -F '\t' '$3 ~ /GPL/ { print $1 "  " $2 "  " $5 "  " $6 "  " $4 }' "$packages"
} >"$out/THIRD_PARTY_NOTICES"

{
	echo Apache-2.0
	cut -f3 "$packages" | tr ' ()' '\n\n\n' | grep -v -x -e AND -e OR -e WITH -e ''
} | sort -u | while read -r license; do
	if [ ! -f "$texts/$license.txt" ]; then
		echo "no text for the license $license in $texts" >&2
		exit 1
	fi
	cp "$texts/$license.txt" "$out/licenses/"
done
