#!/bin/sh
# Copies, into the directory named by $3, the license and notice files of the
# ant binary named by $1 and of every module compiled into it. $2 is ant's
# release version.
#
# The binary records each module's version and content hash. Each module is
# downloaded at that version, checked against that hash, and its license
# files copied, so they are the licenses of the code inside this binary and
# not of whatever the repository requires today. ant's own module is recorded
# without a version, so its license comes from the release's tag. A module
# with no license file, or one that does not match its hash, fails the build.
set -eu

ant=$1
version=$2
out=$3

# The loop below calls this as a condition, where `set -e` does not apply,
# so a failed copy exits by hand.
copy_licenses() {
	found=1
	for file in $(cd "$1" && find . -type f ! -name '*.go' \( -iname 'licen[cs]e*' -o -iname 'copying*' -o -iname 'unlicense*' -o -iname 'notice*' -o -iname 'patents*' \)); do
		case "$(basename "$file" | tr '[:upper:]' '[:lower:]')" in
		notice* | patents*) ;;
		*) [ "$(dirname "$file")" != . ] || found=0 ;;
		esac
		mkdir -p "$2/$(dirname "$file")" || exit 1
		cp -f "$1/$file" "$2/$file" || exit 1
	done
	return "$found"
}

modules=$(mktemp)
trap 'rm -f "$modules"' EXIT
go version -m "$ant" | awk -F '\t' '$2 == "mod" || $2 == "dep" { print $3, $4, $5 }' >"$modules"
[ -s "$modules" ] || {
	echo "$ant records no modules" >&2
	exit 1
}

while read -r module recorded hash; do
	dest=$out/third_party/$module@$recorded
	if [ "$recorded" = "(devel)" ]; then
		recorded=v$version
		dest=$out
	fi
	download=$(go mod download -json "$module@$recorded")
	dir=$(printf '%s\n' "$download" | sed -n 's/^[[:space:]]*"Dir": "\(.*\)",\{0,1\}$/\1/p')
	sum=$(printf '%s\n' "$download" | sed -n 's/^[[:space:]]*"Sum": "\(.*\)",\{0,1\}$/\1/p')
	if [ -n "$hash" ] && [ "$sum" != "$hash" ]; then
		echo "$module@$recorded: downloaded $sum, the binary records $hash" >&2
		exit 1
	fi
	copy_licenses "$dir" "$dest" || {
		echo "$module@$recorded: no license file at the module's root" >&2
		exit 1
	}
done <"$modules"
