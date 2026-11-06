#!/bin/sh
# Removes four root certificates from the trust store of the image whose
# root directory is ROOT, and fails if any is left:
#   sh ca-filter.sh ROOT LAYOUT IMAGE EXPECTED
#
# LAYOUT "bundle": /etc/ssl/certs must hold only ca-certificates.crt, and
# ROOT may have no /usr/share/ca-certificates or
# /usr/local/share/ca-certificates directory. A ROOT that holds only a copy of
# /etc always passes that part.
# LAYOUT "any": Alpine's ca-certificates package may be installed too.
# Its source files for the four, and their lines in /etc/ca-certificates.conf,
# are then removed as well, so a later update-ca-certificates run does not
# bring them back. So are the links that led to them; each is logged. LAYOUT
# "any" follows links correctly only when ROOT is /: see resolve below.
#
# IMAGE names the image in the build log.
# EXPECTED: how many of the four the store should hold, printed beside the
# number removed. "=N" makes the step fail unless it removes exactly N.
#
# The step first tries its own checks on a known certificate.
# It fails on a file whose BEGIN lines, END lines and split-out certificates
# do not give the same count.
# Last, it checks every file Go reads roots from (crypto/x509/root_linux.go).
# It matches by fingerprint, and by each name as single-byte or UTF-16 text.
# A name written with different case or spacing, or as UTF-32 text, is matched
# by fingerprint only.
set -eu
nl='
'
# The four roots, by SHA-256 fingerprint and subject name.
roots='86a1ecba089c4a8d3bbe2734c612ba341d813e043cf9e8a862cd5c57a36bbe6b emSign ECC Root CA - G3
125609aa301da0a249b97a8239cb6a34216f44dcac9f3954b14292f2e8c8608f emSign Root CA - C1
bc4d809b15189d78db3e1d8cf4f9726a795da1643ca5f1358e1ddb0edc0d7eb3 emSign ECC Root CA - C3
d7a7a0fb5d7e2731d771e9484ebcdef71d5f0c3e0a2948782bc83ee0ea699ef4 AAA Certificate Services'
# The same names as hexadecimal bytes: single-byte text, then UTF-16.
hexnames=' 65 6d 53 69 67 6e 20 45 43 43 20 52 6f 6f 74 20 43 41 20 2d 20 47 33
 00 65 00 6d 00 53 00 69 00 67 00 6e 00 20 00 45 00 43 00 43 00 20 00 52 00 6f 00 6f 00 74 00 20 00 43 00 41 00 20 00 2d 00 20 00 47 00 33
 65 6d 53 69 67 6e 20 52 6f 6f 74 20 43 41 20 2d 20 43 31
 00 65 00 6d 00 53 00 69 00 67 00 6e 00 20 00 52 00 6f 00 6f 00 74 00 20 00 43 00 41 00 20 00 2d 00 20 00 43 00 31
 65 6d 53 69 67 6e 20 45 43 43 20 52 6f 6f 74 20 43 41 20 2d 20 43 33
 00 65 00 6d 00 53 00 69 00 67 00 6e 00 20 00 45 00 43 00 43 00 20 00 52 00 6f 00 6f 00 74 00 20 00 43 00 41 00 20 00 2d 00 20 00 43 00 33
 41 41 41 20 43 65 72 74 69 66 69 63 61 74 65 20 53 65 72 76 69 63 65 73
 00 41 00 41 00 41 00 20 00 43 00 65 00 72 00 74 00 69 00 66 00 69 00 63 00 61 00 74 00 65 00 20 00 53 00 65 00 72 00 76 00 69 00 63 00 65 00 73'
fail() {
  printf 'ca-filter: %s\n' "$*" >&2
  echo "ca-filter: this step removes these root certificates:" >&2
  printf '%s\n' "$roots" | while read -r h n; do echo "ca-filter:   $n" >&2; done
  echo "ca-filter: the image was not built." >&2
  [ -z "${work:-}" ] || : >"$work/failed"
  exit 1
}
for c in awk base64 cat find grep head mktemp mv od readlink rm sed sha256sum sort tr wc; do
  command -v "$c" >/dev/null || fail "the image this step runs in has no $c command, which this step needs"
done
root=$(printf %s "$1" | sed 's,/*$,,') || fail "cannot read ROOT"
layout=$2
image=$3
expect=$4
case $layout in bundle | any) ;; *) fail "LAYOUT must be bundle or any, not '$layout'" ;; esac
case $expect in =*) want=${expect#=} ;; *) want=$expect ;; esac
case $want in '' | *[!0-9]*) fail "EXPECTED must be a number, or = and a number, not '$expect'" ;; esac
ssl=$root/etc/ssl
certs=$ssl/certs
bundle=$certs/ca-certificates.crt
conf=$root/etc/ca-certificates.conf
share=$root/usr/share/ca-certificates
localsrc=$root/usr/local/share/ca-certificates
work=$(mktemp -d) || fail "cannot make a temporary directory"
# When a command fails outside fail(), still say that the image was not built.
trap 'rc=$?
if [ "$rc" != 0 ] && [ ! -e "$work/failed" ]; then
  echo "ca-filter: a command above failed." >&2
  echo "ca-filter: the image was not built." >&2
fi
rm -rf "$work" || :
exit "$rc"' EXIT
trap 'exit 1' INT TERM
printf '%s\n' "$hexnames" >"$work/names"
# The same names, one file per root, in the order of the list.
k=0
while IFS= read -r hx; do
  k=$((k + 1))
  printf '%s\n' "$hx" >>"$work/n$(((k + 1) / 2))"
done <<EOF
$hexnames
EOF
fps=" "
while read -r h rest; do fps="$fps$h "; done <<EOF
$roots
EOF

# split_pem FILE writes each certificate in FILE to its own file in a new
# directory, and prints the directory.
split_pem() {
  in=$1
  d=$(mktemp -d "$work/s.XXXXXX") || exit 1
  awk -v dir="$d" '
    /\r/ { exit 1 }
    /^-----BEGIN / { if (inb || $0 != "-----BEGIN CERTIFICATE-----") exit 1; inb = 1; k++; f = sprintf("%s/b%05d.pem", dir, k); print > f; next }
    /^-----END / { if (!inb || $0 != "-----END CERTIFICATE-----") exit 1; print > f; close(f); inb = 0; next }
    inb { print > f; next }
    /^[ \t]*$/ { next }
    { exit 1 }
    END { if (inb) exit 1 }
  ' <"$in" || fail "$in holds something other than PEM CERTIFICATE blocks, such as comment text, another PEM type or CRLF line endings. Remove that content, convert the file to LF line endings, or use a base image whose trust store holds only PEM certificates"
  b=$(grep -c '^-----BEGIN CERTIFICATE-----$' "$in" || :)
  e=$(grep -c '^-----END CERTIFICATE-----$' "$in" || :)
  set -- "$d"/b*.pem
  k=$#
  [ -e "$1" ] || k=0
  [ "$b" = "$k" ] && [ "$e" = "$k" ] || fail "$in has $b BEGIN lines and $e END lines, but $k certificates were read"
  printf '%s\n' "$d"
}
# fingerprint BLOCK FILE decodes BLOCK to BLOCK.der and prints its SHA-256. FILE
# names the file in messages only. It runs in a command substitution, where
# some shells ignore set -e, so each command is checked.
fingerprint() {
  sed '1d;$d' "$1" >"$1.body" || fail "cannot read a certificate in $2"
  tr -d ' \t\n' <"$1.body" >"$1.b64" || fail "cannot read a certificate in $2"
  base64 -d <"$1.b64" >"$1.der" || fail "a certificate in $2 is not valid base64"
  sha256sum <"$1.der" >"$1.sum" || fail "cannot hash a certificate in $2"
  read -r h rest <"$1.sum" || fail "cannot hash a certificate in $2"
  case $h in *[!0-9a-f]*) fail "cannot hash a certificate in $2" ;; esac
  [ "${#h}" = 64 ] || fail "cannot hash a certificate in $2"
  echo "$h"
}
# isroot HASH succeeds if HASH is the fingerprint of one of the four. It uses
# only a shell case, so nothing in it can fail.
isroot() {
  case $fps in *" $1 "*) return 0 ;; esac
  return 1
}
# rootname HASH prints the name of the root with that fingerprint. It uses
# only shell builtins, so nothing in it can fail.
rootname() {
  while read -r rh rname; do
    if [ "$rh" = "$1" ]; then
      printf '%s\n' "$rname"
      return 0
    fi
  done <<EOF
$roots
EOF
}
# describe LINE prints, for a message, the certificate that LINE names; LINE
# is one "HASH NAME" line of found's output.
describe() {
  if isroot "${1%% *}"; then
    printf '%s\n' "${1#* }"
  else
    printf 'a certificate with another fingerprint that contains the name %s\n' "${1#* }"
  fi
}
# named BLOCK succeeds if BLOCK.der contains one of the four names, and leaves
# BLOCK.hex. od writes one space before each byte; tr joins its lines into one,
# so a match starts on a byte. Callers use it as a condition, where set -e
# does not apply, so each command is checked.
named() {
  od -An -v -tx1 <"$1.der" >"$1.od" || fail "cannot read a certificate's name"
  tr -s ' \n' '  ' <"$1.od" >"$1.hex" || fail "cannot read a certificate's name"
  nrc=0
  grep -qF -f "$work/names" "$1.hex" || nrc=$?
  [ "$nrc" -le 1 ] || fail "cannot read a certificate's name"
  return "$nrc"
}
# whichname BLOCK prints which of the four names BLOCK.hex holds.
whichname() {
  i=0
  while read -r rh rname; do
    i=$((i + 1))
    wrc=0
    grep -qF -f "$work/n$i" "$1.hex" || wrc=$?
    [ "$wrc" -le 1 ] || fail "cannot read a certificate's name"
    if [ "$wrc" = 0 ]; then
      printf '%s\n' "$rname"
      return 0
    fi
  done <<EOF
$roots
EOF
}
# found FILE prints "HASH NAME" for each certificate in FILE that is one of
# the four, or carries one of their names.
found() {
  s=$(split_pem "$1") || exit 1
  for c in "$s"/b*.pem; do
    [ -e "$c" ] || continue
    h=$(fingerprint "$c" "$1") || exit 1
    if isroot "$h"; then
      printf '%s %s\n' "$h" "$(rootname "$h")"
    elif named "$c"; then
      printf '%s %s\n' "$h" "$(whichname "$c")"
    fi
  done
}
# nobreak PATH DIR stops the step if PATH, a name in DIR, holds a line break:
# paths pass through files one per line.
nobreak() {
  case $1 in *"$nl"*) fail "a file name in $2 holds a line break, which this step cannot check. Rename or remove it" ;; esac
}

# Try the fingerprint, the name match and the list on emSign ECC Root CA - C3.
cat >"$work/known.pem" <<'PEM'
-----BEGIN CERTIFICATE-----
MIICKzCCAbGgAwIBAgIKe3G2gla4EnycqDAKBggqhkjOPQQDAzBaMQswCQYDVQQG
EwJVUzETMBEGA1UECxMKZW1TaWduIFBLSTEUMBIGA1UEChMLZU11ZGhyYSBJbmMx
IDAeBgNVBAMTF2VtU2lnbiBFQ0MgUm9vdCBDQSAtIEMzMB4XDTE4MDIxODE4MzAw
MFoXDTQzMDIxODE4MzAwMFowWjELMAkGA1UEBhMCVVMxEzARBgNVBAsTCmVtU2ln
biBQS0kxFDASBgNVBAoTC2VNdWRocmEgSW5jMSAwHgYDVQQDExdlbVNpZ24gRUND
IFJvb3QgQ0EgLSBDMzB2MBAGByqGSM49AgEGBSuBBAAiA2IABP2lYa57JhAd6bci
MK4G9IGzsUJxlTm801Ljr6/58pc1kjZGDoeVjbk5Wum739D+yAdBPLtVb4Ojavti
sIGJAnB9SMVK4+kiVCJNk7tCDK93nCOmfddhEc5lx/h//vXyqaNCMEAwHQYDVR0O
BBYEFPtaSNCAIEDyqOkAB2kZd6fmw/TPMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMB
Af8EBTADAQH/MAoGCCqGSM49BAMDA2gAMGUCMQC02C8Cif22TGK6Q04ThHK1rt0c
3ta13FaPWEBaLd4gTCKDypOofu4SQMfWh0/434UCMBwUZOR8loMRnLDRWmFLpg9J
0wD8ofzkpf9/rdcw0Md3f76BB1UwUCAU9Vc4CqgxUQ==
-----END CERTIFICATE-----
PEM
s=$(split_pem "$work/known.pem") || exit 1
[ "$(fingerprint "$s/b00001.pem" "$work/known.pem")" = bc4d809b15189d78db3e1d8cf4f9726a795da1643ca5f1358e1ddb0edc0d7eb3 ] ||
  fail "the image this step runs in computes a wrong fingerprint for a known certificate"
named "$s/b00001.pem" && [ "$(whichname "$s/b00001.pem")" = "emSign ECC Root CA - C3" ] ||
  fail "this step's name check does not work in the image it runs in"
isroot bc4d809b15189d78db3e1d8cf4f9726a795da1643ca5f1358e1ddb0edc0d7eb3 &&
  [ "$(rootname bc4d809b15189d78db3e1d8cf4f9726a795da1643ca5f1358e1ddb0edc0d7eb3)" = "emSign ECC Root CA - C3" ] ||
  fail "this step does not read its own list of root certificates correctly"

[ -f "$bundle" ] && [ ! -L "$bundle" ] || fail "$bundle is not a regular file"
# Go's other root paths on Linux must not exist.
absent() { [ ! -e "$1" ] && [ ! -L "$1" ] || fail "$1 exists. Go reads roots from it, and this step cannot check that path"; }
absent "$root/etc/pki/tls/certs/ca-bundle.crt"
absent "$ssl/ca-bundle.pem"
absent "$root/etc/pki/tls/cacert.pem"
absent "$root/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem"
absent "$root/etc/pki/tls/certs"
# These trust-store directories may not be links: a link among them would
# resolve outside ROOT.
for d in "$root/etc" "$ssl" "$certs" "$root/etc/pki" "$root/usr" "$root/usr/share" "$share" "$root/usr/local" "$root/usr/local/share" "$localsrc"; do
  [ ! -L "$d" ] || fail "$d is a link; this step needs the trust-store directories to be real directories"
done
if [ -e "$root/etc/pki" ]; then
  find "$root/etc/pki" -type l >"$work/pkilinks"
  [ ! -s "$work/pkilinks" ] || fail "$(head -n 1 "$work/pkilinks") is a link under $root/etc/pki, which this step cannot check"
fi
# normpath PATH prints PATH from /, with "." and ".." parts folded in, so
# two spellings of one path compare equal. It is used only to compare paths,
# never to open one. It runs in a command substitution, so set -f and IFS=/
# stay inside it.
normpath() {
  set -f
  np=
  IFS=/
  for part in $1; do
    case $part in
      '' | .) ;;
      ..) np=${np%/*} ;;
      *) np=$np/$part ;;
    esac
  done
  printf '%s\n' "${np:-/}"
}
# resolve PATH prints where PATH leads, reading each absolute link target
# inside ROOT. It stops the step after 16 links. The system resolves each
# directory the chain passes through that is not in the list above. So with a
# ROOT other than /, a link through such a directory can resolve outside ROOT.
resolve() {
  p=$1
  n=0
  while [ -L "$p" ]; do
    n=$((n + 1))
    [ "$n" -le 16 ] || fail "$1 leads through more than 16 links"
    # A target may end in line breaks, which $(...) would drop.
    t=$(readlink "$p" && echo .) || exit 1
    t=${t%.}
    t=${t%"$nl"}
    case $t in *"$nl"*) fail "$1 leads to a path that holds a line break, which this step cannot check" ;; esac
    case $t in
      /*) p=$root$t ;;
      *) p=${p%/*}/$t ;;
    esac
  done
  printf '%s\n' "$p"
}
# /etc/ssl/cert.pem may only be a link to the bundle.
if [ -e "$ssl/cert.pem" ] || [ -L "$ssl/cert.pem" ]; then
  ct=$(readlink "$ssl/cert.pem" && echo .) || ct=.
  ct=${ct%.}
  case ${ct%"$nl"} in
    certs/ca-certificates.crt | /etc/ssl/certs/ca-certificates.crt) ;;
    *) fail "$ssl/cert.pem is not a link to $bundle" ;;
  esac
fi
if [ "$layout" = bundle ]; then
  for d in "$share" "$localsrc"; do
    [ ! -e "$d" ] || fail "$d exists, but this image's store should hold only $bundle"
  done
  for f in "$certs"/* "$certs"/.[!.]* "$certs"/..?*; do
    [ -L "$f" ] || [ -e "$f" ] || continue
    nobreak "$f" "$certs"
    [ "$f" = "$bundle" ] || fail "$f exists, but this image's store should hold only $bundle"
  done
fi

for f in "$certs"/* "$certs"/.[!.]* "$certs"/..?*; do
  nobreak "$f" "$certs"
done

# Remove the four's source files, and their lines in the configuration.
: >"$work/sources"
for d in "$share" "$localsrc"; do
  [ -e "$d" ] || continue
  find "$d" -name "*$nl*" >"$work/odd"
  [ ! -s "$work/odd" ] || fail "a file name in $d holds a line break, which this step cannot check. Rename or remove it"
  find "$d" ! -type d ! -type f >"$work/odd"
  [ ! -s "$work/odd" ] || fail "$(head -n 1 "$work/odd") is not a regular file"
  find "$d" -type f >>"$work/sources"
done
: >"$work/gone"
: >"$work/deleted"
while IFS= read -r f; do
  h=$(found "$f") || exit 1
  [ -n "$h" ] || continue
  one=$(split_pem "$f") || exit 1
  set -- "$one"/b*.pem
  [ "$#" = 1 ] && [ -e "$1" ] || fail "$f holds other certificates beside one of the four; split it so each root has its own file"
  isroot "${h%% *}" || fail "$f holds $(describe "$h")"
  printf '%s\n' "${h%% *}" >>"$work/gone"
  case $f in
    "$share"/*)
      line=${f#"$share"/}
      [ -f "$conf" ] || fail "$conf is missing"
      nline=$(grep -cxF -e "$line" "$conf" || :)
      [ "$nline" = 1 ] || fail "$conf names $line $nline times, but this step needs exactly one line naming it. Edit $conf so one line names it, or remove $f before this step"
      rc=0
      grep -vxF -e "$line" "$conf" >"$work/conf" || rc=$?
      [ "$rc" -le 1 ] || fail "cannot read $conf"
      cat "$work/conf" >"$conf"
      ;;
  esac
  rm -f "$f"
  nf=$(normpath "$f") || exit 1
  printf '%s\n' "$nf" >>"$work/deleted"
done <"$work/sources"

# Remove each link in /etc/ssl/certs whose target, with . and .. folded as
# text, is a file or link this step deleted; then links to those. A link that
# reaches a deleted file through a linked directory is kept, and can no
# longer be read.
pass=0
while :; do
  again=0
  for f in "$certs"/* "$certs"/.[!.]* "$certs"/..?*; do
    [ -L "$f" ] || continue
    nobreak "$f" "$certs"
    r=$(resolve "$f") || exit 1
    nr=$(normpath "$r") || exit 1
    drc=0
    grep -qxF -e "$nr" "$work/deleted" || drc=$?
    [ "$drc" -le 1 ] || fail "cannot read this step's list of deleted files"
    if [ "$drc" = 0 ]; then
      rm -f "$f"
      nf=$(normpath "$f") || exit 1
      printf '%s\n' "$nf" >>"$work/deleted"
      printf 'ca-filter: removed link %s\n' "$f"
      again=1
    fi
  done
  [ "$again" = 1 ] || break
  pass=$((pass + 1))
  [ "$pass" -lt 5 ] || fail "links in $certs still lead to deleted files after $pass passes"
done

# Rewrite the bundle without the four.
s=$(split_pem "$bundle") || exit 1
total=0
: >"$work/out.pem"
: >"$work/seen"
for c in "$s"/b*.pem; do
  [ -e "$c" ] || continue
  total=$((total + 1))
  h=$(fingerprint "$c" "$bundle") || exit 1
  if isroot "$h"; then
    echo "$h" >>"$work/seen"
  elif named "$c"; then
    fail "$bundle holds $(describe "$h $(whichname "$c")")"
  else
    cat "$c" >>"$work/out.pem"
  fi
done
kept=$(grep -c '^-----BEGIN CERTIFICATE-----$' "$work/out.pem" || :)
sort -u "$work/seen" >"$work/seen.u"
removed=$(wc -l <"$work/seen.u") || exit 1
removed=$((removed + 0))
[ $((total - kept)) = "$removed" ] || fail "$bundle holds one of the four more than once"
[ "$kept" -gt 0 ] || fail "$bundle would hold no certificate"
case $expect in
  =*) [ "$removed" = "$want" ] || fail "$image removed $removed of the four, where exactly $want were required. If you build on another base image, change the count on this step's RUN line" ;;
esac
if [ -s "$work/gone" ]; then
  sort -u "$work/gone" >"$work/gone.u"
  [ "$(cat "$work/gone.u")" = "$(cat "$work/seen.u")" ] ||
    fail "$bundle and the source files under /usr/share and /usr/local/share disagree about which of the four they hold"
fi
rm -f "$bundle.new" || fail "cannot replace $bundle"
cat "$work/out.pem" >"$bundle.new"
mv -f "$bundle.new" "$bundle"

# Check every source file, and every file Go reads in /etc/ssl/certs. A link
# that points at nothing is skipped: Go and OpenSSL cannot read it either.
: >"$work/checked"
# check FILE LABEL fails if FILE holds one of the four. LABEL names it.
check() {
  h=$(found "$1") || exit 1
  [ -z "$h" ] || fail "$2 still holds $(describe "$(printf '%s\n' "$h" | head -n 1)")"
  printf '%s\n' "$1" >>"$work/checked"
}
: >"$work/files"
for d in "$share" "$localsrc"; do
  if [ -e "$d" ]; then find "$d" -type f >>"$work/files"; fi
done
while IFS= read -r f; do check "$f" "$f"; done <"$work/files"
for f in "$certs"/* "$certs"/.[!.]* "$certs"/..?*; do
  [ -L "$f" ] || [ -e "$f" ] || continue
  r=$(resolve "$f") || exit 1
  [ -e "$r" ] || continue
  [ ! -d "$r" ] || fail "$f is a directory; this step cannot check what it holds"
  crc=0
  grep -qxF -e "$r" "$work/checked" || crc=$?
  [ "$crc" -le 1 ] || fail "cannot read this step's list of checked files"
  [ "$crc" = 1 ] || continue
  if [ "$r" = "$f" ]; then check "$r" "$f"; else check "$r" "$f (which leads to $r)"; fi
done
# A failure in a subshell that this shell did not stop on leaves this mark.
[ ! -e "$work/failed" ] || exit 1
printf '%s\n' "$roots" | while read -r h name; do
  if grep -qx "$h" "$work/seen"; then echo "ca-filter: removed $name"; fi
done
echo "ca-filter: $bundle keeps $kept of $total certificates"
if [ "$removed" = "$want" ]; then
  printf 'ca-filter: %s removed %s of the four (expected %s)\n' "$image" "$removed" "$want"
else
  printf 'ca-filter: %s removed %s of the four (expected %s). The build goes on. To fail the build when the count differs, write it as =%s on the RUN line\n' "$image" "$removed" "$want" "$want"
fi
