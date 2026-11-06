//go:build unix

package image

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"maps"
	"math/big"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Both images remove four root certificates from their trust store with
// ca-filter.sh. Each Dockerfile mounts it from its build context and runs it.
// testdata/ca-filter/roots.pem holds the four certificates it removes, and
// kept.pem one that it leaves.
//
// TestCAFilterInImages checks the Dockerfiles. TestCAFilter runs the script on
// trust stores it builds.

// caFilterScript is the script's path, relative to sandboxes/.
const caFilterScript = "image/ca-filter.sh"

// caFilterImage is one Dockerfile that runs the script.
type caFilterImage struct {
	dockerfile string // relative to sandboxes/
	source     string // the script's path in the build context, as its RUN line mounts it
	args       string // ROOT LAYOUT IMAGE EXPECTED, as its RUN line passes them
	// context is the directory the build's context is made from, relative to
	// sandboxes/.
	context string
}

// readRepoFile returns the text of a file under sandboxes/.
func readRepoFile(t *testing.T, parts ...string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(append([]string{".."}, parts...)...))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// caFilterText returns the text of image/ca-filter.sh.
func caFilterText(t *testing.T) string {
	t.Helper()
	return readRepoFile(t, strings.Split(caFilterScript, "/")...)
}

// checkCAFilterImages fails unless each Dockerfile runs the script once from a
// file that matches image/ca-filter.sh, leaves the filtered store unchanged after
// that step, and sets neither SSL_CERT_FILE nor SSL_CERT_DIR.
func checkCAFilterImages(t *testing.T, images []caFilterImage) {
	t.Helper()
	text := caFilterText(t)
	alpine := argDefault(t, readRepoFile(t, "image", "Dockerfile"), "BASE_IMAGE")
	for _, img := range images {
		t.Run(img.dockerfile, func(t *testing.T) {
			file := readRepoFile(t, strings.Split(img.dockerfile, "/")...)
			if n := strings.Count(file, "/tmp/ca-filter.sh"); n != 2 {
				t.Fatalf("%s names /tmp/ca-filter.sh %d times, want 2: once where its RUN line mounts the script and once where it runs it", img.dockerfile, n)
			}
			for _, pasted := range []string{"<<'CAFILTER'", "split_pem"} {
				if strings.Contains(file, pasted) {
					t.Errorf("%s holds %q, a sign of a pasted copy of the script. Mount %s instead", img.dockerfile, pasted, img.source)
				}
			}
			run := "\nRUN --mount=type=bind,source=" + img.source + ",target=/tmp/ca-filter.sh sh /tmp/ca-filter.sh " + img.args + "\n"
			start := strings.Index(file, run)
			if start < 0 {
				t.Fatalf("%s has no line %q", img.dockerfile, strings.TrimSpace(run))
			}
			script := path.Join(img.context, img.source)
			if got := readRepoFile(t, strings.Split(script, "/")...); got != text {
				t.Errorf("%s runs %s, which differs from %s. Copy %s over it", img.dockerfile, script, caFilterScript, caFilterScript)
			}
			// A .dockerignore that keeps everything out by default must let the
			// script in.
			if ignore, err := os.ReadFile(filepath.Join("..", filepath.FromSlash(img.context), ".dockerignore")); err == nil {
				lines := strings.Split(string(ignore), "\n")
				if slices.Contains(lines, "*") && !slices.Contains(lines, "!"+img.source) {
					t.Errorf("%s/.dockerignore has no line !%s, so the build of %s cannot read the script", img.context, img.source, img.dockerfile)
				}
			} else if !errors.Is(err, os.ErrNotExist) {
				t.Fatal(err)
			}
			// The stage that runs the script: the lines from the last FROM
			// before the script's RUN line.
			var stage []string
			for _, line := range strings.Split(file[:start], "\n") {
				l := strings.TrimSpace(line)
				if l == "" || strings.HasPrefix(l, "#") {
					continue
				}
				if instruction(l) == "FROM" {
					stage = nil
				}
				stage = append(stage, l)
			}
			if len(stage) == 0 || instruction(stage[0]) != "FROM" {
				t.Fatalf("%s runs the script before any FROM line", img.dockerfile)
			}
			_, stepStage := fromLine(stage[0])
			// Every line after the script's RUN line that is not a comment
			// counts, continuation lines and heredoc bodies included.
			var after []string
			for _, line := range strings.Split(file[start+len(run):], "\n") {
				if l := strings.TrimSpace(line); l != "" && !strings.HasPrefix(l, "#") {
					after = append(after, l)
				}
			}
			// Every line of the file that is not a comment.
			var lines []string
			for _, line := range strings.Split(file, "\n") {
				if l := strings.TrimSpace(line); l != "" && !strings.HasPrefix(l, "#") {
					lines = append(lines, l)
				}
			}
			// An ENV line's continuation lines count as ENV too.
			inEnv := false
			for _, l := range lines {
				env := inEnv || instruction(l) == "ENV"
				if env && strings.Contains(l, "SSL_CERT_") {
					t.Errorf("%s sets SSL_CERT_FILE or SSL_CERT_DIR, which could bring back the four removed root certificates: %s", img.dockerfile, l)
				}
				inEnv = env && strings.HasSuffix(l, "\\")
			}
			switch {
			case strings.HasPrefix(img.args, "/ "):
				filtered := map[string]bool{stepStage: true}
				for _, l := range after {
					if instruction(l) == "FROM" {
						base, name := fromLine(l)
						if !filtered[base] {
							t.Errorf("%s starts a stage after the script that does not build on the filtered stage: %s", img.dockerfile, l)
						}
						filtered[name] = true
						continue
					}
					if strings.Contains(l, "apk ") || changesTrustStore(l) {
						t.Errorf("%s changes the trust store after the script, which could bring back the four removed root certificates: %s. Move it before the step", img.dockerfile, l)
					}
				}
			case strings.HasPrefix(img.args, "/base "):
				wantFrom := "--platform=$BUILDPLATFORM ${FILTER_IMAGE} AS trust-store"
				if strings.Join(strings.Fields(stage[0])[1:], " ") != wantFrom {
					t.Errorf("%s runs the script in a stage that does not start with the line: FROM %s", img.dockerfile, wantFrom)
				}
				named := map[string]int{}
				for i, l := range lines {
					if instruction(l) != "FROM" {
						continue
					}
					_, name := fromLine(l)
					named[strings.ToLower(name)]++
					if strings.EqualFold(name, "base") && i+1 < len(lines) && instruction(lines[i+1]) != "FROM" {
						t.Errorf("%s's base stage holds more than its FROM line: %s", img.dockerfile, lines[i+1])
					}
				}
				for _, name := range []string{"base", "trust-store"} {
					if named[name] != 1 {
						t.Errorf("%s names %d stages %q, want 1", img.dockerfile, named[name], name)
					}
				}
				if !slices.Contains(stage, "COPY --from=base /etc /base/etc") {
					t.Errorf("%s's filtering stage does not copy the base's /etc with: COPY --from=base /etc /base/etc", img.dockerfile)
				}
				if len(after) == 0 || after[0] != "FROM base" {
					t.Fatalf("%s's next stage after the script is not: FROM base", img.dockerfile)
				}
				want := "COPY --from=trust-store /base/etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt"
				copied := false
				for _, l := range after[1:] {
					switch {
					case l == want:
						copied = true
					case instruction(l) == "FROM":
						t.Errorf("%s starts another stage after FROM base: %s", img.dockerfile, l)
					case changesTrustStore(l):
						t.Errorf("%s's last stage changes the trust store besides copying the filtered bundle, which could bring back the four removed root certificates: %s", img.dockerfile, l)
					}
				}
				if !copied {
					t.Errorf("%s's last stage does not copy the filtered bundle: %s", img.dockerfile, want)
				}
				if got := argDefault(t, file, "FILTER_IMAGE"); got != alpine {
					t.Errorf("%s's FILTER_IMAGE defaults to %q, and image/Dockerfile's BASE_IMAGE to %q. Bump them together", img.dockerfile, got, alpine)
				}
			default:
				t.Fatalf("%s: a ROOT of %q is neither / nor /base", img.dockerfile, strings.Fields(img.args)[0])
			}
		})
	}
}

// caFilterTrustDirs are the directories Go or Alpine's package keeps roots in.
var caFilterTrustDirs = []string{"/etc/ssl", "/etc/pki", "/usr/share/ca-certificates", "/usr/local/share/ca-certificates"}

// changesTrustStore reports whether a Dockerfile line can change a trust store:
// it names one, sets SSL_CERT_FILE or SSL_CERT_DIR, or is a COPY or ADD whose
// destination could hold one. A relative destination counts.
func changesTrustStore(line string) bool {
	if strings.Contains(line, "update-ca-certificates") {
		return true
	}
	if instruction(line) == "ENV" && strings.Contains(line, "SSL_CERT_") {
		return true
	}
	for _, dir := range caFilterTrustDirs {
		if strings.Contains(line, dir) {
			return true
		}
	}
	fields := strings.Fields(line)
	if len(fields) < 3 || (instruction(line) != "COPY" && instruction(line) != "ADD") {
		return false
	}
	dest := strings.TrimSuffix(fields[len(fields)-1], "/") + "/"
	if !strings.HasPrefix(dest, "/") {
		return true
	}
	for _, dir := range caFilterTrustDirs {
		if strings.HasPrefix(dir+"/", dest) {
			return true
		}
	}
	return false
}

// instruction returns a Dockerfile line's first word in upper case, since
// instructions are not case sensitive.
func instruction(line string) string {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return ""
	}
	return strings.ToUpper(fields[0])
}

// fromLine returns the image or stage a FROM line builds on, and the name it
// gives the new stage ("" when it gives none). It skips flags such as
// --platform.
func fromLine(line string) (base, name string) {
	fields := strings.Fields(line)[1:]
	for len(fields) > 0 && strings.HasPrefix(fields[0], "--") {
		fields = fields[1:]
	}
	if len(fields) > 0 {
		base = fields[0]
	}
	if len(fields) >= 3 && strings.EqualFold(fields[1], "AS") {
		name = fields[2]
	}
	return base, name
}

// argDefault returns the default of the ARG named name in a Dockerfile.
func argDefault(t *testing.T, dockerfile, name string) string {
	t.Helper()
	for _, line := range strings.Split(dockerfile, "\n") {
		if v, ok := strings.CutPrefix(line, "ARG "+name+"="); ok {
			return v
		}
	}
	t.Fatalf("no ARG %s= line", name)
	return ""
}

// caFilterRoot is one of the four certificates the script removes.
type caFilterRoot struct {
	name string
	pem  []byte
	fp   string // SHA-256 of the certificate, in lowercase hexadecimal
}

// caFilterRoots reads the four from testdata, in the order the script lists
// them.
func caFilterRoots(t *testing.T) []caFilterRoot {
	t.Helper()
	rest := []byte(readRepoFile(t, "image", "testdata", "ca-filter", "roots.pem"))
	var roots []caFilterRoot
	for {
		var b *pem.Block
		b, rest = pem.Decode(rest)
		if b == nil {
			break
		}
		c, err := x509.ParseCertificate(b.Bytes)
		if err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256(b.Bytes)
		roots = append(roots, caFilterRoot{name: c.Subject.CommonName, pem: pem.EncodeToMemory(b), fp: hex.EncodeToString(sum[:])})
	}
	var names []string
	for _, r := range roots {
		names = append(names, r.name)
	}
	want := []string{"emSign ECC Root CA - G3", "emSign Root CA - C1", "emSign ECC Root CA - C3", "AAA Certificate Services"}
	if !slices.Equal(names, want) {
		t.Fatalf("roots.pem holds %q, want %q", names, want)
	}
	return roots
}

// caFilterCert returns a new self-signed CA certificate named cn, as PEM.
func caFilterCert(t *testing.T, cn string) []byte {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: cn},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
}

// caFilterBMPCert returns a new self-signed CA certificate whose name is cn
// written as an ASN.1 BMPString, which holds UTF-16 text.
func caFilterBMPCert(t *testing.T, cn string) []byte {
	t.Helper()
	var utf16 []byte
	for _, r := range cn {
		utf16 = append(utf16, byte(r>>8), byte(r))
	}
	name, err := asn1.Marshal(pkix.RDNSequence{{{
		Type:  asn1.ObjectIdentifier{2, 5, 4, 3}, // commonName
		Value: asn1.RawValue{Tag: asn1.TagBMPString, Bytes: utf16},
	}}})
	if err != nil {
		t.Fatal(err)
	}
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		RawSubject:            name,
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
}

// caFilterStore is a trust store the test builds under a scratch root.
type caFilterStore struct {
	root  string
	certs map[string][]byte // file name, without .crt, to PEM
}

// caFilterLayout is how newCAFilterStore lays a store out.
type caFilterLayout int

const (
	bundleOnly         caFilterLayout = iota // only /etc/ssl/certs/ca-certificates.crt
	alpineFull                               // Alpine's package, its links absolute as Alpine makes them
	alpineFullRelative                       // Alpine's package, its links relative
)

// bundlePath is the store's ca-certificates.crt.
func (s caFilterStore) bundlePath() string {
	return filepath.Join(s.root, "etc", "ssl", "certs", "ca-certificates.crt")
}

// newCAFilterStore writes certs into a store of the given layout under a new
// root, whose name holds a comma, an apostrophe and a space.
func newCAFilterStore(t *testing.T, certs map[string][]byte, layout caFilterLayout) caFilterStore {
	t.Helper()
	s := caFilterStore{root: filepath.Join(t.TempDir(), "a,b'c d"), certs: certs}
	certsDir := filepath.Join(s.root, "etc", "ssl", "certs")
	caFilterMkdir(t, certsDir)
	var bundle []byte
	names := slices.Sorted(maps.Keys(certs))
	for _, name := range names {
		bundle = append(bundle, certs[name]...)
		bundle = append(bundle, '\n')
	}
	caFilterWrite(t, s.bundlePath(), bundle)
	caFilterLink(t, "certs/ca-certificates.crt", filepath.Join(s.root, "etc", "ssl", "cert.pem"))
	if layout == bundleOnly {
		return s
	}
	mozilla := filepath.Join(s.root, "usr", "share", "ca-certificates", "mozilla")
	caFilterMkdir(t, mozilla)
	caFilterMkdir(t, filepath.Join(s.root, "usr", "local", "share", "ca-certificates"))
	conf := "# Automatically generated\n#\n"
	for i, name := range names {
		caFilterWrite(t, filepath.Join(mozilla, name+".crt"), certs[name])
		conf += "mozilla/" + name + ".crt\n"
		link := "ca-cert-" + name + ".pem"
		target := "/usr/share/ca-certificates/mozilla/" + name + ".crt"
		if layout == alpineFullRelative {
			target = "../../.." + target
		}
		caFilterLink(t, target, filepath.Join(certsDir, link))
		caFilterLink(t, link, filepath.Join(certsDir, fmt.Sprintf("%08x.0", i)))
	}
	caFilterWrite(t, filepath.Join(s.root, "etc", "ca-certificates.conf"), []byte(conf))
	return s
}

// caFilterShells returns the shells to run the script under: sh always, and
// dash, bash and busybox where this host has them. Under busybox, every command
// the script's tool check names is busybox's own.
func caFilterShells(t *testing.T) map[string][]string {
	t.Helper()
	shells := map[string][]string{"sh": {"sh"}}
	// dash and bash run too when sh is another shell. bash, run as bash, does
	// not carry set -e into command substitutions, which the script must not
	// rely on.
	sh, errSh := exec.LookPath("sh")
	for _, name := range []string{"dash", "bash"} {
		if p, err := exec.LookPath(name); errSh == nil && err == nil && !sameFile(sh, p) {
			shells[name] = []string{name}
		}
	}
	if bb, err := exec.LookPath("busybox"); err == nil {
		shells["busybox"] = []string{bb, "sh"}
	} else {
		t.Log("busybox is not installed here, so the script does not run under it")
	}
	return shells
}

// caFilterTools returns the commands the script's tool check names, from its
// "for c in … ; do" line.
func caFilterTools(t *testing.T, script string) []string {
	t.Helper()
	for _, line := range strings.Split(script, "\n") {
		if rest, ok := strings.CutPrefix(line, "for c in "); ok {
			if list, ok := strings.CutSuffix(rest, "; do"); ok {
				return strings.Fields(list)
			}
		}
	}
	t.Fatal("the script has no tool check line")
	return nil
}

// requireCAFilterTools skips a test on a host without a command the script needs.
func requireCAFilterTools(t *testing.T, script string) {
	t.Helper()
	for _, tool := range caFilterTools(t, script) {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s is not installed", tool)
		}
	}
}

// sameFile reports whether a and b lead to the same file.
func sameFile(a, b string) bool {
	ai, errA := os.Stat(a)
	bi, errB := os.Stat(b)
	return errA == nil && errB == nil && os.SameFile(ai, bi)
}

// runCAFilter runs script on the store with args LAYOUT IMAGE EXPECTED, and returns
// its output and exit code. The shell reads the script from a file, as each
// Dockerfile's RUN line runs it.
func runCAFilter(t *testing.T, shell []string, script string, s caFilterStore, args ...string) (string, int) {
	t.Helper()
	file := filepath.Join(t.TempDir(), "ca-filter.sh")
	caFilterWrite(t, file, []byte(script))
	env := os.Environ()
	if filepath.Base(shell[0]) == "busybox" {
		// The busybox applets come first on PATH, so the script uses them.
		applets := t.TempDir()
		for _, c := range caFilterTools(t, script) {
			caFilterLink(t, shell[0], filepath.Join(applets, c))
		}
		env = append(env, "PATH="+applets+string(os.PathListSeparator)+os.Getenv("PATH"))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, shell[0], append(append(shell[1:], file, s.root), args...)...)
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	var exit *exec.ExitError
	switch {
	case err == nil:
		return string(out), 0
	case errors.As(err, &exit):
		return string(out), exit.ExitCode()
	default:
		t.Fatalf("run %s: %v\n%s", strings.Join(shell, " "), err, out)
		return "", 0
	}
}

// bundleFingerprints returns the fingerprint of each certificate in the
// store's bundle, read with Go's PEM decoder, which is what the server uses.
func (s caFilterStore) bundleFingerprints(t *testing.T) []string {
	t.Helper()
	rest, err := os.ReadFile(s.bundlePath())
	if err != nil {
		t.Fatal(err)
	}
	var fps []string
	for {
		var b *pem.Block
		b, rest = pem.Decode(rest)
		if b == nil {
			break
		}
		sum := sha256.Sum256(b.Bytes)
		fps = append(fps, hex.EncodeToString(sum[:]))
	}
	return fps
}

func caFilterMkdir(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
}

func caFilterWrite(t *testing.T, path string, b []byte) {
	t.Helper()
	caFilterMkdir(t, filepath.Dir(path))
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatal(err)
	}
}

func caFilterLink(t *testing.T, target, path string) {
	t.Helper()
	if err := os.Symlink(target, path); err != nil {
		t.Fatal(err)
	}
}

// TestCAFilter runs the script under each shell this host has, on stores laid
// out as the images' are, and on stores it must refuse. Each refusal names
// what it found and says the image was not built.
func TestCAFilter(t *testing.T) {
	script := caFilterText(t)
	requireCAFilterTools(t, script)
	roots := caFilterRoots(t)
	others := map[string][]byte{}
	for i := range 20 {
		others[fmt.Sprintf("Test_Root_%02d", i)] = caFilterCert(t, fmt.Sprintf("Test Root %02d", i))
	}
	// with returns others plus the given roots, each under its own name.
	with := func(rs ...caFilterRoot) map[string][]byte {
		m := map[string][]byte{}
		for k, v := range others {
			m[k] = v
		}
		for _, r := range rs {
			m[strings.ReplaceAll(r.name, " ", "_")] = r.pem
		}
		return m
	}
	three := roots[:3]
	otherName := caFilterCert(t, "emSign Root CA - G9")
	sameName := caFilterCert(t, "emSign ECC Root CA - G3")
	sameNameBMP := caFilterBMPCert(t, "emSign ECC Root CA - G3")
	tabbed := bytes.Replace(roots[0].pem, []byte("\n"), []byte("\n\t"), 1)

	type tc struct {
		name  string
		store func(t *testing.T) caFilterStore
		args  []string
		fail  string // a part of the refusal; empty when the run succeeds
		keep  int    // certificates left in the bundle on success
		// late is set when the refusal comes from the last check, after the
		// bundle is rewritten. Every other refusal leaves the bundle as it was.
		late bool
		// after, when set, checks the store after a successful run.
		after func(t *testing.T, s caFilterStore)
	}
	cases := []tc{
		{name: "the four in a bundle-only store", store: func(t *testing.T) caFilterStore { return newCAFilterStore(t, with(roots...), bundleOnly) }, args: []string{"bundle", "test", "=4"}, keep: 20},
		{name: "three of the four, any count", store: func(t *testing.T) caFilterStore { return newCAFilterStore(t, with(three...), bundleOnly) }, args: []string{"any", "test", "3"}, keep: 20},
		{name: "none of the four, any count", store: func(t *testing.T) caFilterStore { return newCAFilterStore(t, with(), bundleOnly) }, args: []string{"any", "test", "3"}, keep: 20},
		{name: "four where three are expected, any count", store: func(t *testing.T) caFilterStore { return newCAFilterStore(t, with(roots...), bundleOnly) }, args: []string{"any", "test", "3"}, keep: 20},
		{name: "Alpine's package with absolute links, three of the four", store: func(t *testing.T) caFilterStore { return newCAFilterStore(t, with(three...), alpineFull) }, args: []string{"any", "test", "3"}, keep: 20},
		{name: "emSign Root CA - G1 stays in a bundle-only store", store: func(t *testing.T) caFilterStore {
			m := with(roots...)
			m["emSign_Root_CA_-_G1"] = []byte(readRepoFile(t, "image", "testdata", "ca-filter", "kept.pem"))
			return newCAFilterStore(t, m, bundleOnly)
		}, args: []string{"bundle", "test", "=4"}, keep: 21},
		{name: "emSign Root CA - G1 stays in Alpine's package", store: func(t *testing.T) caFilterStore {
			m := with(roots...)
			m["emSign_Root_CA_-_G1"] = []byte(readRepoFile(t, "image", "testdata", "ca-filter", "kept.pem"))
			return newCAFilterStore(t, m, alpineFull)
		}, args: []string{"any", "test", "=4"}, keep: 21},
		{name: "another emSign name is kept", store: func(t *testing.T) caFilterStore {
			m := with(three...)
			m["emSign_Root_CA_-_G9"] = otherName
			return newCAFilterStore(t, m, bundleOnly)
		}, args: []string{"any", "test", "3"}, keep: 21},
		{name: "three where four are expected", store: func(t *testing.T) caFilterStore { return newCAFilterStore(t, with(three...), bundleOnly) }, args: []string{"bundle", "test", "=4"}, fail: "test removed 3 of the four, where exactly 4 were required"},
		{name: "a same-named certificate with another key", store: func(t *testing.T) caFilterStore {
			m := with(three...)
			m["zz"] = sameName
			return newCAFilterStore(t, m, bundleOnly)
		}, args: []string{"any", "test", "3"}, fail: "a certificate with another fingerprint that contains the name emSign ECC Root CA - G3"},
		{name: "a same-named source file", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), alpineFull)
			caFilterWrite(t, filepath.Join(s.root, "usr", "local", "share", "ca-certificates", "zz.crt"), sameName)
			return s
		}, args: []string{"any", "test", "3"}, fail: "a certificate with another fingerprint that contains the name emSign ECC Root CA - G3"},
		{name: "a second copy with a tab in it", store: func(t *testing.T) caFilterStore {
			m := with(three...)
			m["zz"] = tabbed
			return newCAFilterStore(t, m, bundleOnly)
		}, args: []string{"any", "test", "3"}, fail: "holds one of the four more than once"},
		{name: "carriage returns", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), bundleOnly)
			b, _ := os.ReadFile(s.bundlePath())
			caFilterWrite(t, s.bundlePath(), bytes.ReplaceAll(b, []byte("\n"), []byte("\r\n")))
			return s
		}, args: []string{"any", "test", "3"}, fail: "holds something other than PEM CERTIFICATE blocks, such as comment text, another PEM type or CRLF line endings"},
		{name: "a configuration line for one of the four twice", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), alpineFull)
			conf := filepath.Join(s.root, "etc", "ca-certificates.conf")
			b, _ := os.ReadFile(conf)
			caFilterWrite(t, conf, append(b, "mozilla/"+strings.ReplaceAll(three[0].name, " ", "_")+".crt\n"...))
			return s
		}, args: []string{"any", "test", "3"}, fail: "ca-certificates.conf names mozilla/emSign_ECC_Root_CA_-_G3.crt 2 times, but this step needs exactly one line naming it"},
		{name: "no configuration line for one of the four", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), alpineFull)
			conf := filepath.Join(s.root, "etc", "ca-certificates.conf")
			b, _ := os.ReadFile(conf)
			caFilterWrite(t, conf, bytes.Replace(b, []byte("mozilla/"+strings.ReplaceAll(three[0].name, " ", "_")+".crt\n"), nil, 1))
			return s
		}, args: []string{"any", "test", "3"}, fail: "ca-certificates.conf names mozilla/emSign_ECC_Root_CA_-_G3.crt 0 times, but this step needs exactly one line naming it"},
		{name: "nothing else in the bundle", store: func(t *testing.T) caFilterStore {
			m := map[string][]byte{}
			for _, r := range roots {
				m[r.fp] = r.pem
			}
			return newCAFilterStore(t, m, bundleOnly)
		}, args: []string{"bundle", "test", "4"}, fail: "would hold no certificate"},
		{name: "another file beside a bundle-only store's bundle", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(roots...), bundleOnly)
			caFilterWrite(t, filepath.Join(filepath.Dir(s.bundlePath()), "extra.pem"), others["Test_Root_00"])
			return s
		}, args: []string{"bundle", "test", "4"}, fail: "extra.pem exists, but this image's store should hold only"},
		{name: "one of the four in another file", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), bundleOnly)
			caFilterWrite(t, filepath.Join(filepath.Dir(s.bundlePath()), "extra.pem"), roots[0].pem)
			return s
		}, args: []string{"any", "test", "3"}, fail: "extra.pem still holds emSign ECC Root CA - G3", late: true},
		{name: "one of the four behind an absolute link", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), bundleOnly)
			// The target exists only under s.root. On this host the link leads nowhere.
			// The script must read the target inside s.root to find the file.
			caFilterWrite(t, filepath.Join(s.root, "opt", "ca-filter-test", "root.pem"), roots[0].pem)
			caFilterLink(t, "/opt/ca-filter-test/root.pem", filepath.Join(filepath.Dir(s.bundlePath()), "root.pem"))
			return s
		}, args: []string{"any", "test", "3"}, fail: "root.pem) still holds emSign ECC Root CA - G3", late: true},
		{name: "a directory in /etc/ssl/certs", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), bundleOnly)
			caFilterMkdir(t, filepath.Join(filepath.Dir(s.bundlePath()), "sub"))
			return s
		}, args: []string{"any", "test", "3"}, fail: "sub is a directory; this step cannot check what it holds", late: true},
		{name: "another root path Go reads", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), bundleOnly)
			caFilterMkdir(t, filepath.Join(s.root, "etc", "pki", "tls", "certs"))
			return s
		}, args: []string{"any", "test", "3"}, fail: "certs exists. Go reads roots from it, and this step cannot check that path"},
		{name: "cert.pem not a link", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), bundleOnly)
			link := filepath.Join(s.root, "etc", "ssl", "cert.pem")
			if err := os.Remove(link); err != nil {
				t.Fatal(err)
			}
			caFilterWrite(t, link, roots[0].pem)
			return s
		}, args: []string{"any", "test", "3"}, fail: "cert.pem is not a link to"},
		{name: "a source file with two certificates", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), alpineFull)
			src := filepath.Join(s.root, "usr", "share", "ca-certificates", "mozilla", "emSign_ECC_Root_CA_-_G3.crt")
			caFilterWrite(t, src, append(append([]byte{}, roots[0].pem...), others["Test_Root_00"]...))
			return s
		}, args: []string{"any", "test", "3"}, fail: "holds other certificates beside one of the four"},
		{name: "Alpine's package with relative links", store: func(t *testing.T) caFilterStore { return newCAFilterStore(t, with(three...), alpineFullRelative) }, args: []string{"any", "test", "3"}, keep: 20},
		{name: "a link to a file not there yet is kept", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), bundleOnly)
			caFilterLink(t, "/run/secrets/operator.pem", filepath.Join(filepath.Dir(s.bundlePath()), "operator.pem"))
			return s
		}, args: []string{"any", "test", "3"}, keep: 20, after: caFilterKept("operator.pem")},
		{name: "Alpine's package with a hidden link to one of the four", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), alpineFull)
			caFilterLink(t, "/usr/share/ca-certificates/mozilla/emSign_ECC_Root_CA_-_G3.crt", filepath.Join(filepath.Dir(s.bundlePath()), ".root.pem"))
			return s
		}, args: []string{"any", "test", "3"}, keep: 20, after: caFilterGone(".root.pem")},
		{name: "Alpine's package with a link spelled with . and //", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), alpineFull)
			caFilterLink(t, "./..//certs/../../../usr/share/ca-certificates//mozilla/./emSign_ECC_Root_CA_-_G3.crt", filepath.Join(filepath.Dir(s.bundlePath()), "dots.pem"))
			return s
		}, args: []string{"any", "test", "3"}, keep: 20, after: caFilterGone("dots.pem")},
		{name: "Alpine's package with a link to a file not there yet", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), alpineFull)
			caFilterLink(t, "/run/secrets/corp-ca.pem", filepath.Join(filepath.Dir(s.bundlePath()), "corp-ca.pem"))
			return s
		}, args: []string{"any", "test", "3"}, keep: 20, after: caFilterKept("corp-ca.pem")},
		{name: "a bundle-only store with a file named from two dots", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(roots...), bundleOnly)
			caFilterWrite(t, filepath.Join(filepath.Dir(s.bundlePath()), "..x.pem"), others["Test_Root_00"])
			return s
		}, args: []string{"bundle", "test", "4"}, fail: "..x.pem exists, but this image's store should hold only"},
		{name: "one of the four in a file named from two dots", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), bundleOnly)
			caFilterWrite(t, filepath.Join(filepath.Dir(s.bundlePath()), "..aaa.pem"), roots[3].pem)
			return s
		}, args: []string{"any", "test", "3"}, fail: "..aaa.pem still holds AAA Certificate Services", late: true},
		{name: "a same-named certificate in another file", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), bundleOnly)
			caFilterWrite(t, filepath.Join(filepath.Dir(s.bundlePath()), "extra.pem"), sameName)
			return s
		}, args: []string{"any", "test", "3"}, fail: "extra.pem still holds a certificate with another fingerprint that contains the name emSign ECC Root CA - G3", late: true},
		{name: "a bundle-only store with a hidden file", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(roots...), bundleOnly)
			caFilterWrite(t, filepath.Join(filepath.Dir(s.bundlePath()), ".x.pem"), others["Test_Root_00"])
			return s
		}, args: []string{"bundle", "test", "4"}, fail: ".x.pem exists, but this image's store should hold only"},
		{name: "a bundle-only store with a link to a file not there yet", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(roots...), bundleOnly)
			caFilterLink(t, "missing.pem", filepath.Join(filepath.Dir(s.bundlePath()), "link.pem"))
			return s
		}, args: []string{"bundle", "test", "4"}, fail: "link.pem exists, but this image's store should hold only"},
		{name: "one of the four in a hidden file", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), bundleOnly)
			caFilterWrite(t, filepath.Join(filepath.Dir(s.bundlePath()), ".aaa.pem"), roots[3].pem)
			return s
		}, args: []string{"any", "test", "3"}, fail: ".aaa.pem still holds AAA Certificate Services", late: true},
		{name: "a hidden link to a file not there yet is kept", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), bundleOnly)
			caFilterLink(t, "/run/secrets/operator.pem", filepath.Join(filepath.Dir(s.bundlePath()), ".operator.pem"))
			return s
		}, args: []string{"any", "test", "3"}, keep: 20, after: caFilterKept(".operator.pem")},
		{name: "a leftover link where the new bundle is written", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), bundleOnly)
			dir := filepath.Dir(s.bundlePath())
			caFilterWrite(t, filepath.Join(dir, "mine.pem"), others["Test_Root_01"])
			caFilterLink(t, "mine.pem", s.bundlePath()+".new")
			return s
		}, args: []string{"any", "test", "3"}, keep: 20, after: func(t *testing.T, s caFilterStore) {
			b, err := os.ReadFile(filepath.Join(filepath.Dir(s.bundlePath()), "mine.pem"))
			if err != nil || !bytes.Equal(b, others["Test_Root_01"]) {
				t.Errorf("mine.pem was overwritten or lost: %v", err)
			}
			if fi, err := os.Lstat(s.bundlePath()); err != nil || !fi.Mode().IsRegular() {
				t.Errorf("the bundle is not a regular file: %v", err)
			}
		}},
		{name: "a line break in a name in /etc/ssl/certs", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), bundleOnly)
			caFilterWrite(t, filepath.Join(filepath.Dir(s.bundlePath()), "x\ny.pem"), roots[0].pem)
			return s
		}, args: []string{"any", "test", "3"}, fail: "certs holds a line break, which this step cannot check. Rename or remove it"},
		{name: "a link whose target ends in a line break", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), bundleOnly)
			caFilterWrite(t, filepath.Join(s.root, "etc", "ssl", "evil\n"), roots[1].pem)
			caFilterLink(t, "../evil\n", filepath.Join(filepath.Dir(s.bundlePath()), "zz.pem"))
			return s
		}, args: []string{"any", "test", "3"}, fail: "zz.pem leads to a path that holds a line break, which this step cannot check"},
		{name: "a line break in a source file's name", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), alpineFull)
			caFilterWrite(t, filepath.Join(s.root, "usr", "local", "share", "ca-certificates", "a\nb.crt"), others["Test_Root_00"])
			return s
		}, args: []string{"any", "test", "3"}, fail: "ca-certificates holds a line break, which this step cannot check. Rename or remove it"},
		{name: "a same-named AAA certificate whose name is UTF-16 text", store: func(t *testing.T) caFilterStore {
			m := with(three...)
			m["zz"] = caFilterBMPCert(t, "AAA Certificate Services")
			return newCAFilterStore(t, m, bundleOnly)
		}, args: []string{"any", "test", "3"}, fail: "contains the name AAA Certificate Services"},
		{name: "a same-named emSign ECC Root CA - C3", store: func(t *testing.T) caFilterStore {
			m := with(three...)
			m["zz"] = caFilterCert(t, "emSign ECC Root CA - C3")
			return newCAFilterStore(t, m, bundleOnly)
		}, args: []string{"any", "test", "3"}, fail: "contains the name emSign ECC Root CA - C3"},
		{name: "a same-named certificate whose name is UTF-16 text", store: func(t *testing.T) caFilterStore {
			m := with(three...)
			m["zz"] = sameNameBMP
			return newCAFilterStore(t, m, bundleOnly)
		}, args: []string{"any", "test", "3"}, fail: "a certificate with another fingerprint that contains the name emSign ECC Root CA - G3"},
		{name: "one of the four with carriage returns in another file", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), bundleOnly)
			caFilterWrite(t, filepath.Join(filepath.Dir(s.bundlePath()), "extra.pem"), bytes.ReplaceAll(roots[3].pem, []byte("\n"), []byte("\r\n")))
			return s
		}, args: []string{"any", "test", "3"}, fail: "extra.pem holds something other than PEM CERTIFICATE blocks, such as comment text, another PEM type or CRLF line endings. Remove that content, convert the file to LF line endings", late: true},
		{name: "a loop of links", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), alpineFull)
			dir := filepath.Dir(s.bundlePath())
			caFilterLink(t, "loop2.pem", filepath.Join(dir, "loop1.pem"))
			caFilterLink(t, "loop1.pem", filepath.Join(dir, "loop2.pem"))
			return s
		}, args: []string{"any", "test", "3"}, fail: "leads through more than 16 links"},
		{name: "a source directory that is a link", store: func(t *testing.T) caFilterStore {
			s := newCAFilterStore(t, with(three...), alpineFull)
			share := filepath.Join(s.root, "usr", "share", "ca-certificates")
			if err := os.Rename(share, share+"-real"); err != nil {
				t.Fatal(err)
			}
			caFilterLink(t, "ca-certificates-real", share)
			return s
		}, args: []string{"any", "test", "3"}, fail: "ca-certificates is a link; this step needs the trust-store directories to be real directories"},
		{name: "a count that is not a number", store: func(t *testing.T) caFilterStore { return newCAFilterStore(t, with(three...), bundleOnly) }, args: []string{"any", "test", ""}, fail: "EXPECTED must be a number, or = and a number, not ''"},
		{name: "a layout it does not know", store: func(t *testing.T) caFilterStore { return newCAFilterStore(t, with(three...), bundleOnly) }, args: []string{"Bundle", "test", "5"}, fail: "LAYOUT must be bundle or any, not 'Bundle'"},
	}
	for shellName, shell := range caFilterShells(t) {
		t.Run(shellName, func(t *testing.T) {
			for _, c := range cases {
				t.Run(c.name, func(t *testing.T) {
					s := c.store(t)
					before, err := os.ReadFile(s.bundlePath())
					if err != nil {
						t.Fatal(err)
					}
					out, code := runCAFilter(t, shell, script, s, c.args...)
					if c.fail != "" {
						if code == 0 || !strings.Contains(out, c.fail) || !strings.Contains(out, "the image was not built") {
							t.Fatalf("exit %d, want a refusal saying %q and that the image was not built:\n%s", code, c.fail, out)
						}
						after, err := os.ReadFile(s.bundlePath())
						if err != nil {
							t.Fatal(err)
						}
						if !c.late && !bytes.Equal(before, after) {
							t.Errorf("the refusal changed the bundle")
						}
						return
					}
					if code != 0 {
						t.Fatalf("exit %d:\n%s", code, out)
					}
					fps := s.bundleFingerprints(t)
					removed := countCAFilterRoots(s.certs, roots)
					expected := strings.TrimPrefix(c.args[2], "=")
					if want := fmt.Sprintf("ca-filter: test removed %d of the four (expected %s)", removed, expected); !strings.Contains(out, want) {
						t.Errorf("the output lacks %q:\n%s", want, out)
					}
					if goesOn := strconv.Itoa(removed) != expected; goesOn != strings.Contains(out, "The build goes on.") {
						t.Errorf("removed %d of %s expected, and the output says the build goes on: %v, want %v:\n%s", removed, expected, !goesOn, goesOn, out)
					}
					if len(fps) != c.keep {
						t.Errorf("the bundle keeps %d certificates, want %d", len(fps), c.keep)
					}
					for _, r := range roots {
						if slices.Contains(fps, r.fp) {
							t.Errorf("the bundle still holds %s", r.name)
						}
					}
					checkCAFilterFullStore(t, s, roots)
					if c.after != nil {
						c.after(t, s)
					}
				})
			}
		})
	}
}

// checkCAFilterFullStore fails if a store in Alpine's layout still has a source
// file, configuration line or link for one of the four, or a ca-cert or hash
// link that leads nowhere inside its root.
func checkCAFilterFullStore(t *testing.T, s caFilterStore, roots []caFilterRoot) {
	t.Helper()
	conf, err := os.ReadFile(filepath.Join(s.root, "etc", "ca-certificates.conf"))
	if errors.Is(err, os.ErrNotExist) {
		return
	} else if err != nil {
		t.Fatal(err)
	}
	for _, r := range roots {
		file := strings.ReplaceAll(r.name, " ", "_")
		if strings.Contains(string(conf), file) {
			t.Errorf("ca-certificates.conf still names %s", file)
		}
		if _, err := os.Lstat(filepath.Join(s.root, "usr", "share", "ca-certificates", "mozilla", file+".crt")); err == nil {
			t.Errorf("the source file of %s is still there", r.name)
		}
	}
	certsDir := filepath.Dir(s.bundlePath())
	entries, err := os.ReadDir(certsDir)
	if err != nil {
		t.Fatal(err)
	}
	// Each kept certificate has its ca-cert link and hash link; the bundle is one
	// more. A link the test added that points at nothing is not counted.
	live := 0
	for _, e := range entries {
		if _, err := os.Stat(resolveIn(t, s.root, filepath.Join(certsDir, e.Name()))); err == nil {
			live++
		} else if strings.HasPrefix(e.Name(), "ca-cert-") || strings.HasSuffix(e.Name(), ".0") {
			t.Errorf("%s points at nothing: %v", e.Name(), err)
		}
	}
	if want := 2*(len(s.certs)-countCAFilterRoots(s.certs, roots)) + 1; live != want {
		t.Errorf("/etc/ssl/certs holds %d entries that lead to a file, want %d", live, want)
	}
}

// caFilterKept returns a check that the link named name is still in the
// store's /etc/ssl/certs.
func caFilterKept(name string) func(t *testing.T, s caFilterStore) {
	return func(t *testing.T, s caFilterStore) {
		t.Helper()
		if _, err := os.Lstat(filepath.Join(filepath.Dir(s.bundlePath()), name)); err != nil {
			t.Errorf("%s is gone: %v", name, err)
		}
	}
}

// caFilterGone returns a check that the link named name is no longer in the
// store's /etc/ssl/certs.
func caFilterGone(name string) func(t *testing.T, s caFilterStore) {
	return func(t *testing.T, s caFilterStore) {
		t.Helper()
		if _, err := os.Lstat(filepath.Join(filepath.Dir(s.bundlePath()), name)); err == nil {
			t.Errorf("%s is still there", name)
		}
	}
}

// resolveIn follows the links in path as the image would, with an absolute
// target read inside root.
func resolveIn(t *testing.T, root, path string) string {
	t.Helper()
	for range 17 {
		target, err := os.Readlink(path)
		if err != nil {
			return path // not a link
		}
		if filepath.IsAbs(target) {
			path = filepath.Join(root, target)
		} else {
			path = filepath.Join(filepath.Dir(path), target)
		}
	}
	t.Fatalf("%s leads through more than 16 links", path)
	return ""
}

// countCAFilterRoots counts how many of the four certs holds.
func countCAFilterRoots(certs map[string][]byte, roots []caFilterRoot) int {
	n := 0
	for _, p := range certs {
		for _, r := range roots {
			if bytes.Equal(p, r.pem) {
				n++
			}
		}
	}
	return n
}

// TestCAFilterRefusesABrokenCopy shows that the script's own checks catch a
// copy of it that is broken in one place, on a store it otherwise accepts.
func TestCAFilterRefusesABrokenCopy(t *testing.T) {
	script := caFilterText(t)
	requireCAFilterTools(t, script)
	roots := caFilterRoots(t)
	certs := map[string][]byte{}
	for i := range 5 {
		certs[fmt.Sprintf("Test_Root_%02d", i)] = caFilterCert(t, fmt.Sprintf("Test Root %02d", i))
	}
	for _, r := range roots[:3] {
		certs[r.fp] = r.pem
	}
	for _, c := range []struct {
		name, old, new, fail string
	}{
		{"a reader that drops the third certificate", `f = sprintf("%s/b%05d.pem", dir, k)`, `f = (k == 3) ? "/dev/null" : sprintf("%s/b%05d.pem", dir, k)`, "certificates were read"},
		{"a wrong fingerprint for the known certificate", "= bc4d809b15189d78db3e1d8cf4f9726a795da1643ca5f1358e1ddb0edc0d7eb3 ]", "= 0c4d809b15189d78db3e1d8cf4f9726a795da1643ca5f1358e1ddb0edc0d7eb3 ]", "computes a wrong fingerprint for a known certificate"},
		{"no names to look for", `printf '%s\n' "$hexnames" >"$work/names"`, `: >"$work/names"`, "this step's name check does not work"},
	} {
		t.Run(c.name, func(t *testing.T) {
			broken := strings.Replace(script, c.old, c.new, 1)
			if broken == script {
				t.Fatalf("%q is not in the script", c.old)
			}
			out, code := runCAFilter(t, []string{"sh"}, broken, newCAFilterStore(t, certs, bundleOnly), "any", "test", "3")
			if code == 0 || !strings.Contains(out, c.fail) {
				t.Fatalf("exit %d, want a refusal saying %q:\n%s", code, c.fail, out)
			}
		})
	}
}

// The dispatcher image filters a copy of its distroless base's store, and the
// sandbox image filters its own store in place. Each passes its image's name and
// how many of the four its base is expected to hold.
func TestCAFilterInImages(t *testing.T) {
	checkCAFilterImages(t, []caFilterImage{
		{dockerfile: "dispatcher/Dockerfile", source: "image/ca-filter.sh", context: ".", args: "/base bundle cma-dispatcher 4"},
		{dockerfile: "image/Dockerfile", source: "ca-filter.sh", context: "image", args: "/ any cma-sandbox 3"},
	})
}
