#!/bin/sh
# Copies, into the directory named by $1, the license and notice files of
# every module compiled into the Go package named by $2. The image ships them
# because the binary carries that code.
#
# A package's license is the nearest one at or above its directory, so code a
# module took from elsewhere keeps the license that sits beside it. A package
# with no license file anywhere up to its module's root fails the build: the
# image never ships code without its license.
#
# Run from the module's root, with the modules already downloaded, and with
# GOOS and GOARCH set as for the build, so the packages are the binary's.
set -eu

out=$1
package=$2

# Copies $1's license and notice files to $2. Succeeds only if a license was
# among them; a NOTICE or PATENTS file alone does not license anything. The
# loops below call it as a condition, where `set -e` does not apply, so a
# failed copy exits by hand. -f because the module cache's files are
# read-only and several packages share one license.
copy_notices() {
	found=1
	for file in "$1"/*; do
		[ -f "$file" ] || continue
		case "$(basename "$file" | tr '[:upper:]' '[:lower:]')" in
		*.go) continue ;;
		licen[cs]e* | copying* | unlicense*) found=0 ;;
		notice* | patents*) ;;
		*) continue ;;
		esac
		mkdir -p "$2" || exit 1
		cp -f "$file" "$2/" || exit 1
	done
	return "$found"
}

packages=$(mktemp)
trap 'rm -f "$packages"' EXIT
go list -deps -f '{{with .Module}}{{if not .Main}}{{.Path}}@{{.Version}}	{{.Dir}}	{{$.Dir}}{{end}}{{end}}' "$package" >"$packages"

while IFS='	' read -r module moduledir packagedir; do
	dir=$packagedir
	until copy_notices "$dir" "$out/$module${dir#"$moduledir"}"; do
		if [ "$dir" = "$moduledir" ]; then
			echo "$module: no license file at or above $packagedir" >&2
			exit 1
		fi
		dir=$(dirname "$dir")
	done
	copy_notices "$moduledir" "$out/$module" || true
done <"$packages"
