package kubernetes

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// The script may apply only the manifest it records, with the image at the
// digest it records, on the named cluster, and by default over no other controller.
func TestInstallAgentSandbox(t *testing.T) {
	source, err := os.ReadFile("install-agent-sandbox.sh")
	if err != nil {
		t.Fatal(err)
	}
	recorded := func(name string) string {
		match := regexp.MustCompile(`(?m)^` + name + `=(\S+)$`).FindSubmatch(source)
		if match == nil {
			t.Fatalf("the script sets no %s", name)
		}
		return string(match[1])
	}
	image, version, digest := recorded("image"), recorded("version"), recorded("digest")
	manifest := "containers:\n- image: " + image + ":" + version + "\n"
	sum := sha256.Sum256([]byte(manifest))
	script := regexp.MustCompile(`(?m)^asset_sha256=.*$`).ReplaceAll(source, []byte("asset_sha256="+hex.EncodeToString(sum[:])))

	const mirror = "registry.example.com/controller"
	for name, tc := range map[string]struct {
		args     []string
		manifest string // the file behind --manifest, or what curl downloads
		local    bool
		want     string // the image applied, or for a refusal a part of what it printed
		refused  bool
		there    string // what kubectl get prints for the Deployment
		crds     string // what it prints for the definitions
		getFails bool
		crdsFail bool
	}{
		"a download":          {args: []string{"--context", "prod"}, manifest: manifest, want: image + ":" + version + "@" + digest},
		"a local manifest":    {args: []string{"--context=prod"}, manifest: manifest, local: true, want: image + ":" + version + "@" + digest},
		"a mirrored image":    {args: []string{"--context", "prod", "--image", mirror}, manifest: manifest, want: mirror + "@" + digest},
		"no cluster":          {manifest: manifest, refused: true, want: "name the cluster to install on"},
		"another manifest":    {args: []string{"--context", "prod"}, manifest: manifest + "# changed\n", refused: true, want: "does not match the recorded SHA-256"},
		"another local file":  {args: []string{"--context", "prod"}, manifest: "kind: Secret\n", local: true, refused: true, want: "does not match the recorded SHA-256"},
		"another controller":  {args: []string{"--context", "prod"}, manifest: manifest, there: "example.com/c:v2", refused: true, want: "nothing was changed"},
		"definitions alone":   {args: []string{"--context", "prod"}, manifest: manifest, crds: "  definition sandboxes", refused: true, want: "nothing was changed"},
		"an unreadable state": {args: []string{"--context", "prod"}, manifest: manifest, getFails: true, refused: true, want: "nothing was changed"},
		"unread definitions":  {args: []string{"--context", "prod"}, manifest: manifest, crdsFail: true, refused: true, want: "nothing was changed"},
		"the same controller": {args: []string{"--context", "prod"}, manifest: manifest, there: "x@" + digest, crds: "  definition sandboxes", want: image + ":" + version + "@" + digest},
		"a replacement":       {args: []string{"--context", "prod", "--replace"}, manifest: manifest, getFails: true, want: image + ":" + version + "@" + digest},
		"another image bytes": {args: []string{"--context", "prod", "--image", mirror + "@sha256:0"}, manifest: manifest, refused: true, want: "names the digest"},
	} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			write := func(name, text string) string {
				file := filepath.Join(dir, name)
				if err := os.WriteFile(file, []byte(text), 0o700); err != nil {
					t.Fatal(err)
				}
				return file
			}
			release, applied := write("release", tc.manifest), filepath.Join(dir, "applied")
			write("curl", "#!/bin/sh\nwhile [ \"$1\" != -o ]; do shift; done\ncat \"$RELEASE\" >\"$2\"\n")
			write("kubectl", "#!/bin/sh\n[ \"$1 $2\" = '--context prod' ] || exit 9\ncase \" $* \" in *' get '*) [ -z \"$GET_FAILS\" ] || exit 1; case \" $* \" in *' deployment '*) printf %s \"$THERE\" ;; *) [ -z \"$CRDS_FAIL\" ] || exit 1; printf %s \"$CRDS\" ;; esac; exit ;; esac\ncase \" $* \" in *' apply '*) for last; do :; done; cat \"$last\" >\"$APPLIED\" ;; esac\n")
			args := append([]string{write("script", string(script))}, tc.args...)
			if tc.local {
				args = append(args, "--manifest", release)
			}
			cmd := exec.Command("sh", args...)
			cmd.Env = append(os.Environ(), "PATH="+dir+":"+os.Getenv("PATH"), "RELEASE="+release, "APPLIED="+applied, "THERE="+tc.there, "CRDS="+tc.crds)
			if tc.crdsFail {
				cmd.Env = append(cmd.Env, "CRDS_FAIL=1")
			}
			if tc.getFails {
				cmd.Env = append(cmd.Env, "GET_FAILS=1")
			}
			out, err := cmd.CombinedOutput()

			if tc.refused {
				if err == nil || !strings.Contains(string(out), tc.want) {
					t.Errorf("want a refusal that says %q, got %v:\n%s", tc.want, err, out)
				}
				if _, err := os.Stat(applied); err == nil {
					t.Error("it refused, and yet applied a manifest")
				}
				return
			}
			if err != nil {
				t.Fatalf("%v:\n%s", err, out)
			}
			got, err := os.ReadFile(applied)
			if err != nil {
				t.Fatal(err)
			}
			if want := "containers:\n- image: \"" + tc.want + "\"\n"; string(got) != want {
				t.Errorf("applied %q, want %q", got, want)
			}
		})
	}
}
