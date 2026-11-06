#!/bin/sh
# Collects, into the directory named by $2, the complete source of every
# package in the apk database named by $1 that is under a GPL or LGPL
# license: Alpine's build recipe, patches and install scripts as they stood
# at the commit the package was built from, and beside them the upstream
# archives the recipe names. SOURCES in that directory lists what is there.
#
# The recipes come from Alpine's own GitHub mirror of its package tree, which
# carries the same commits as gitlab.alpinelinux.org. abuild, Alpine's build
# tool, downloads the archives and checks them against the recipe's
# checksums, the way Alpine does before it builds. A package whose recipe or
# archives cannot be had fails the build: the source image never lacks the
# source of something the sandbox image ships.
set -eu

installed=$1
out=$2
aports=${APORTS_REPOSITORY:-https://github.com/alpinelinux/aports.git}

# A download fails now and then; each is tried up to four times, the wait
# doubling from 10 s.
retry() {
	try=1
	wait=10
	until "$@"; do
		if [ "$try" -ge 4 ]; then
			echo "giving up on: $*" >&2
			return 1
		fi
		echo "failed, trying again in ${wait}s: $*" >&2
		sleep "$wait"
		try=$((try + 1))
		wait=$((wait * 2))
	done
}

mkdir -p "$out"
awk 'BEGIN { RS = ""; FS = "\n"; OFS = "\t" }
{
	split("", field)
	for (i = 1; i <= NF; i++) field[substr($i, 1, 1)] = substr($i, 3)
	if (field["L"] ~ /GPL/) print field["P"], field["V"], field["L"], field["o"], field["c"]
}' "$installed" | sort >"$out/SOURCES.tsv"
[ -s "$out/SOURCES.tsv" ] || {
	echo "no GPL or LGPL package in $installed" >&2
	exit 1
}

# abuild asks this mirror for an archive before it asks the upstream host.
# Alpine keeps there every archive a release was built from, which outlives
# an upstream that moves or refuses a download. The release is the one the
# packages belong to, read from the same database.
release=$(awk '/^P:alpine-release$/ { found = 1 } found && /^V:/ { split(substr($0, 3), v, "."); print v[1] "." v[2]; exit }' "$installed")
[ -n "$release" ] || {
	echo "no alpine-release package in $installed" >&2
	exit 1
}
export DISTFILES_MIRROR="https://distfiles.alpinelinux.org/distfiles/v$release"

cut -f4,5 "$out/SOURCES.tsv" | sort -u >"$out/origins.tsv"
while IFS='	' read -r origin commit; do
	tree=$(mktemp -d)
	git -C "$tree" init -q
	git -C "$tree" remote add origin "$aports"
	git -C "$tree" sparse-checkout set --no-cone "/main/$origin/" "/community/$origin/"
	retry git -C "$tree" fetch -q --depth 1 --filter=blob:none origin "$commit"
	retry git -C "$tree" checkout -q FETCH_HEAD
	recipe=$tree/main/$origin
	[ -d "$recipe" ] || recipe=$tree/community/$origin
	[ -f "$recipe/APKBUILD" ] || {
		echo "$origin: no APKBUILD in commit $commit" >&2
		exit 1
	}
	cp -R "$recipe" "$out/$origin"
	cd "$out/$origin"
	SRCDEST=$out/$origin retry abuild -F fetch
	SRCDEST=$out/$origin abuild -F verify
	rm -rf src
	cd "$out"
	rm -rf "$tree"
done <"$out/origins.tsv"
rm "$out/origins.tsv"

{
	cat <<'TEXT'
The source of the GPL and LGPL packages in the cma-sandbox image published
under the same tag as this one. Each directory holds one Alpine source
package: its APKBUILD recipe, patches and install scripts at the commit of
https://gitlab.alpinelinux.org/alpine/aports it was built from, and the
upstream archives the recipe names. Alpine's abuild builds the package from
the directory as it is.

PACKAGE  VERSION  LICENSE  DIRECTORY  APORTS COMMIT
TEXT
	awk -F '\t' '{ print $1 "  " $2 "  " $3 "  " $4 "  " $5 }' "$out/SOURCES.tsv"
} >"$out/SOURCES"
rm "$out/SOURCES.tsv"
