//go:build unix

// Package image tests the pod's shell scripts by running them.
package image

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"
)

// packedSecret is a work item's secret as the queue mints it: URL-safe
// base64 of a JSON payload embedding the sessions token. The scripts must
// pass it through verbatim; the worker's CLI unpacks it.
var packedSecret = base64.RawURLEncoding.EncodeToString([]byte(`{"version":1,"sessions_token":"tok","api_base_url":"https://api"}`))

// The worker container's session code can rename entries of the shared
// volume, so everything fetch reads back must live on its own volume. A
// scratch file inside the shared volume could be swapped for a forged work
// item between fetch's write and its read. The stubs record every file curl
// is told to write and every file jq is told to read, so the test sees the
// scratch while fetch is using it, not the end state after fetch removed it.
func TestFetchKeepsItsScratchOffTheSharedVolume(t *testing.T) {
	run := startFetch(t, nil)
	run.waitForSecret(t)

	entries, err := os.ReadDir(run.sharedDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Name() != "work-item" {
			t.Errorf("the shared volume holds %q; fetch must keep everything but the work item off it", e.Name())
		}
	}
	if secret, err := os.ReadFile(run.secretPath()); err != nil || string(secret) != packedSecret {
		t.Errorf("secret = %q, %v; want the queue's payload verbatim", secret, err)
	}
	for _, log := range []string{"curl-writes", "jq-reads"} {
		paths := run.loggedPaths(t, log)
		if len(paths) == 0 {
			t.Errorf("%s logged nothing; the stub never ran", log)
		}
		for _, p := range paths {
			if !under(p, run.fetchDir) || under(p, run.sharedDir) {
				t.Errorf("%s: %s; fetch must keep its scratch on its own volume %s", log, p, run.fetchDir)
			}
		}
	}
}

// The worker asks fetch to record its run by creating worker-starting in the
// shared volume. The record must land in the work-item directory, which fetch
// alone can write, and only once asked for: a record written before the
// worker asks blocks the first worker as if it were a restarted one. fetch
// must also write the collected marker on its own volume before any session
// code can run.
func TestFetchRecordsTheRunWhereTheSessionCannotWrite(t *testing.T) {
	run := startFetch(t, nil)
	run.waitForSecret(t)
	run.waitForCollected(t)
	// Gives fetch time for its next step after writing the collected marker:
	// a record written without the worker's request would be there by now.
	time.Sleep(500 * time.Millisecond)
	if exists(filepath.Join(run.sharedDir, "work-item", "worker-ran")) {
		t.Fatal("fetch recorded the run before the worker asked")
	}

	run.requestRecord(t)
	if exists(filepath.Join(run.sharedDir, "worker-ran")) {
		t.Error("fetch wrote the record in the shared root, where the session's code can delete it")
	}
	info, err := os.Stat(filepath.Join(run.sharedDir, "work-item"))
	if err != nil {
		t.Fatal(err)
	}
	if mode := info.Mode().Perm(); mode != 0o750 {
		t.Errorf("work-item mode = %o; the record is only as safe as the directory that holds it", mode)
	}
}

// The session's code can rename entries of the shared volume while fetch
// waits for the worker's request, and can plant a symlink under the name
// fetch validated. The record must land in the directory fetch entered, not
// wherever that name leads by then. fetch sleeps between checks for the
// request and finds sleep through PATH, so a sleep stub that logs its
// directory shows when fetch waits inside the work-item directory.
func TestFetchRecordsTheRunInTheDirectoryItEntered(t *testing.T) {
	for _, tool := range []string{"sh", "sleep"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s is not installed", tool)
		}
	}
	run := &scriptRun{sharedDir: t.TempDir(), fetchDir: t.TempDir(), binDir: t.TempDir()}
	realSleep, _ := exec.LookPath("sleep")
	writeStub(t, run.binDir, "sleep", `#!/bin/sh
pwd >> "`+filepath.Join(run.binDir, "sleep-dirs")+`"
exec "`+realSleep+`" "$@"
`)
	// A restarted fetch with its work item in place goes straight to the
	// wait, with no dispatcher call.
	if err := os.WriteFile(filepath.Join(run.fetchDir, "collected"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	writeWorkItem(t, run.sharedDir)
	tokenFile := filepath.Join(run.binDir, "sa-token")
	if err := os.WriteFile(tokenFile, []byte("sa"), 0o600); err != nil {
		t.Fatal(err)
	}
	run.start(t, "fetch.sh",
		"CMA_DISPATCHER_URL=https://dispatcher",
		"CMA_SA_TOKEN_FILE="+tokenFile,
		"CMA_CA_FILE="+tokenFile,
		"CMA_SHARED_DIR="+run.sharedDir,
		"CMA_FETCH_DIR="+run.fetchDir,
	)
	workItem := filepath.Join(run.sharedDir, "work-item")
	run.waitFor(t, "fetch never waited for the worker's request inside the work-item directory", func() bool {
		for _, dir := range run.loggedLines(t, "sleep-dirs") {
			if dir == workItem {
				return true
			}
		}
		return false
	})

	renamed := filepath.Join(run.sharedDir, "work-item.renamed")
	if err := os.Rename(workItem, renamed); err != nil {
		t.Fatal(err)
	}
	decoy := t.TempDir()
	if err := os.Symlink(decoy, workItem); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(run.sharedDir, "worker-starting"), nil, 0o644); err != nil {
		t.Fatal(err)
	}

	inRenamed := filepath.Join(renamed, "worker-ran")
	inDecoy := filepath.Join(decoy, "worker-ran")
	run.waitFor(t, "fetch never recorded the run", func() bool { return exists(inRenamed) || exists(inDecoy) })
	if exists(inDecoy) {
		t.Error("fetch wrote the record through the symlink the session's code planted")
	}
	if !exists(inRenamed) {
		t.Error("fetch did not write the record in the directory it entered")
	}
}

// A fetch container restarted after it collected the work item must not
// collect it again: the dispatcher answers the same pod with the same work
// item until the worker exits, and a fresh work item with no run record
// beside it would let a restarted worker run the session twice. A work-item
// directory that is gone, renamed away by the session's code, stays gone.
func TestFetchCollectsNothingAgainAfterARestart(t *testing.T) {
	run := startFetch(t, func(fetchDir string) {
		if err := os.WriteFile(filepath.Join(fetchDir, "collected"), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	})
	run.waitFor(t, "fetch never reported the missing work item", func() bool {
		return strings.Contains(run.stderr(t), "this pod's work item is gone; not collecting it again")
	})
	if paths := run.loggedPaths(t, "curl-writes"); len(paths) != 0 {
		t.Errorf("fetch called the dispatcher after the handoff: curl wrote %v", paths)
	}
	if exists(filepath.Join(run.sharedDir, "work-item")) {
		t.Error("fetch made a work-item directory after the handoff")
	}
}

// A fetch container restarted with the work item it wrote still in place,
// whether or not it had written the collected marker, records the worker's
// run without calling the dispatcher.
func TestFetchRecordsTheRunAfterItsOwnRestart(t *testing.T) {
	for _, tc := range []struct {
		name      string
		collected bool
	}{
		{name: "collected marker written", collected: true},
		{name: "collected marker not yet written", collected: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			run := startFetch(t, func(fetchDir string) {
				if tc.collected {
					if err := os.WriteFile(filepath.Join(fetchDir, "collected"), nil, 0o644); err != nil {
						t.Fatal(err)
					}
				}
			}, writeWorkItem)
			run.requestRecord(t)
			run.waitForCollected(t)
			if paths := run.loggedPaths(t, "curl-writes"); len(paths) != 0 {
				t.Errorf("fetch called the dispatcher for a work item it had: curl wrote %v", paths)
			}
		})
	}
}

// scriptRun is one run of a sandbox script against stubbed tools.
type scriptRun struct {
	sharedDir, fetchDir, binDir, stderrPath string
	ctx                                     context.Context
	// credential, when set, runs the script as that user rather than the
	// test's.
	credential *syscall.Credential
	// pid is the script's process, for a test that sends it a stop signal
	// as the kubelet does: to that process alone, not its group.
	pid int
	// exited closes once the script's process has been waited for, and
	// waitErr then holds the wait's result.
	exited  <-chan struct{}
	waitErr error
}

// startFetch runs fetch.sh with a curl that answers every request 200 with
// a complete work item. prepareFetch, if given, seeds fetch's own volume
// first, and prepareShared the shared volume. The process group is killed
// when the test ends: fetch.sh idles forever, sleeping in a background
// child.
func startFetch(t *testing.T, prepareFetch func(fetchDir string), prepareShared ...func(t *testing.T, sharedDir string)) *scriptRun {
	t.Helper()
	for _, tool := range []string{"sh", "jq"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s is not installed", tool)
		}
	}
	run := &scriptRun{sharedDir: t.TempDir(), fetchDir: t.TempDir(), binDir: t.TempDir()}
	if prepareFetch != nil {
		prepareFetch(run.fetchDir)
	}
	for _, prepare := range prepareShared {
		prepare(t, run.sharedDir)
	}
	writeStub(t, run.binDir, "curl", `#!/bin/sh
out=""
while [ $# -gt 0 ]; do
	case "$1" in
	-o) out="$2"; shift 2 ;;
	-H) echo "$2" >> "`+filepath.Join(run.binDir, "curl-headers")+`"; shift 2 ;;
	*) shift ;;
	esac
done
echo "$out" >> "`+filepath.Join(run.binDir, "curl-writes")+`"
printf '%s' '{"session_id":"s_1","work_id":"w_1","environment_id":"e_1","base_url":"https://api","secret":"`+packedSecret+`"}' > "$out"
printf 200
`)
	realJQ, _ := exec.LookPath("jq")
	writeStub(t, run.binDir, "jq", `#!/bin/sh
for last; do :; done
echo "$last" >> "`+filepath.Join(run.binDir, "jq-reads")+`"
exec "`+realJQ+`" "$@"
`)
	tokenFile := filepath.Join(run.binDir, "sa-token")
	if err := os.WriteFile(tokenFile, []byte("sa"), 0o600); err != nil {
		t.Fatal(err)
	}
	run.start(t, "fetch.sh",
		"CMA_DISPATCHER_URL=https://dispatcher",
		"CMA_SA_TOKEN_FILE="+tokenFile,
		"CMA_CA_FILE="+tokenFile,
		"CMA_SHARED_DIR="+run.sharedDir,
		"CMA_FETCH_DIR="+run.fetchDir,
	)
	return run
}

// loggedLines returns the lines a stub appended to its log. It returns none
// when the stub never ran, and none while the log is still empty. A stub's >>
// creates its log before its echo writes the line. A wait would otherwise read
// that moment as one empty line.
func (r *scriptRun) loggedLines(t *testing.T, log string) []string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(r.binDir, log))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	text := strings.TrimRight(string(data), "\n")
	if text == "" {
		return nil
	}
	return strings.Split(text, "\n")
}

// A wait for a stub's first line must outlast the moment its log exists but
// is still empty. A wait that ends there reads an empty string as the line.
func TestWaitForALoggedLineOutlastsAnEmptyLog(t *testing.T) {
	// Given: a log that a stub's >> has created but not yet written.
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	t.Cleanup(cancel)
	run := &scriptRun{binDir: t.TempDir(), ctx: ctx}
	log := filepath.Join(run.binDir, "sleeps")
	if err := os.WriteFile(log, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	// When: the test waits for the first line, which lands after the first read.
	reads := 0
	run.waitFor(t, "the line never showed up in the log", func() bool {
		reads++
		if reads == 2 {
			if err := os.WriteFile(log, []byte("3600\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return len(run.loggedLines(t, "sleeps")) > 0
	})

	// Then: the wait ends on the line, not on the empty log.
	if lines := run.loggedLines(t, "sleeps"); len(lines) != 1 || lines[0] != "3600" {
		t.Errorf("the wait ended with the log reading %q; want [\"3600\"]", lines)
	}
}

// On the add-on path fetch has no token file; it presents the assignment
// nonce in a header instead, and starts without CMA_SA_TOKEN_FILE set.
func TestFetchAuthenticatesWithTheNonce(t *testing.T) {
	for _, tool := range []string{"sh", "jq"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s is not installed", tool)
		}
	}
	run := &scriptRun{sharedDir: t.TempDir(), fetchDir: t.TempDir(), binDir: t.TempDir()}
	writeStub(t, run.binDir, "curl", `#!/bin/sh
out=""
while [ $# -gt 0 ]; do
	case "$1" in
	-o) out="$2"; shift 2 ;;
	-H) echo "$2" >> "`+filepath.Join(run.binDir, "curl-headers")+`"; shift 2 ;;
	*) shift ;;
	esac
done
printf '%s' '{"session_id":"s_1","work_id":"w_1","environment_id":"e_1","base_url":"https://api","secret":"`+packedSecret+`"}' > "$out"
printf 200
`)
	realJQ, _ := exec.LookPath("jq")
	writeStub(t, run.binDir, "jq", "#!/bin/sh\nexec \""+realJQ+"\" \"$@\"\n")
	caFile := filepath.Join(run.binDir, "ca")
	if err := os.WriteFile(caFile, []byte("pem"), 0o600); err != nil {
		t.Fatal(err)
	}
	run.start(t, "fetch.sh",
		"CMA_DISPATCHER_URL=https://dispatcher",
		"CMA_ASSIGNMENT_NONCE=the-nonce",
		"CMA_CA_FILE="+caFile,
		"CMA_SHARED_DIR="+run.sharedDir,
		"CMA_FETCH_DIR="+run.fetchDir,
	)
	run.waitForSecret(t)
	headers := run.loggedLines(t, "curl-headers")
	sawNonce := false
	for _, h := range headers {
		if strings.HasPrefix(h, "Authorization:") {
			t.Errorf("fetch sent a bearer token on the nonce path: %q", h)
		}
		if h == "X-CMA-Assignment-Nonce: the-nonce" {
			sawNonce = true
		}
	}
	if !sawNonce {
		t.Errorf("fetch never sent the assignment nonce header; sent %v", headers)
	}
}

// On the token path fetch presents the projected service-account token,
// read from its file on every call. The kubelet rotates the token, and one
// read at startup would go stale while fetch polls. The curl stub answers
// 404 until a request carries the rotated token.
func TestFetchAuthenticatesWithTheServiceAccountToken(t *testing.T) {
	for _, tool := range []string{"sh", "jq"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s is not installed", tool)
		}
	}
	run := &scriptRun{sharedDir: t.TempDir(), fetchDir: t.TempDir(), binDir: t.TempDir()}
	writeStub(t, run.binDir, "curl", `#!/bin/sh
out=""
auth=""
while [ $# -gt 0 ]; do
	case "$1" in
	-o) out="$2"; shift 2 ;;
	-H) echo "$2" >> "`+filepath.Join(run.binDir, "curl-headers")+`"; auth="$2"; shift 2 ;;
	*) shift ;;
	esac
done
if [ "$auth" != "Authorization: Bearer rotated" ]; then
	printf 404
	exit 0
fi
printf '%s' '{"session_id":"s_1","work_id":"w_1","environment_id":"e_1","base_url":"https://api","secret":"`+packedSecret+`"}' > "$out"
printf 200
`)
	realJQ, _ := exec.LookPath("jq")
	writeStub(t, run.binDir, "jq", "#!/bin/sh\nexec \""+realJQ+"\" \"$@\"\n")
	tokenFile := filepath.Join(run.binDir, "sa-token")
	if err := os.WriteFile(tokenFile, []byte("sa"), 0o600); err != nil {
		t.Fatal(err)
	}
	run.start(t, "fetch.sh",
		"CMA_DISPATCHER_URL=https://dispatcher",
		"CMA_SA_TOKEN_FILE="+tokenFile,
		"CMA_CA_FILE="+tokenFile,
		"CMA_SHARED_DIR="+run.sharedDir,
		"CMA_FETCH_DIR="+run.fetchDir,
	)
	run.waitFor(t, "fetch never asked the dispatcher for the work item", func() bool {
		return len(run.loggedLines(t, "curl-headers")) > 0
	})
	if first := run.loggedLines(t, "curl-headers")[0]; first != "Authorization: Bearer sa" {
		t.Fatalf("fetch's first request carried %q; want the token from the file", first)
	}

	// Rotate the token: a new file takes the old one's place.
	rotated := filepath.Join(run.binDir, "sa-token.rotated")
	if err := os.WriteFile(rotated, []byte("rotated"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(rotated, tokenFile); err != nil {
		t.Fatal(err)
	}
	run.waitFor(t, "fetch never sent the rotated token; it must read the token file on each call", func() bool {
		return exists(run.secretPath())
	})
	for _, h := range run.loggedLines(t, "curl-headers") {
		if strings.HasPrefix(h, "X-CMA-Assignment-Nonce:") {
			t.Errorf("fetch sent the nonce header on the token path: %q", h)
		}
	}
}

// writeWorkItem lays out a complete work item as fetch leaves it.
func writeWorkItem(t *testing.T, sharedDir string) {
	t.Helper()
	workItem := filepath.Join(sharedDir, "work-item")
	if err := os.Mkdir(workItem, 0o750); err != nil {
		t.Fatal(err)
	}
	for field, value := range map[string]string{
		"session_id": "s_1", "work_id": "w_1", "environment_id": "e_1", "base_url": "https://api", "secret": packedSecret,
	} {
		if err := os.WriteFile(filepath.Join(workItem, field), []byte(value), 0o440); err != nil {
			t.Fatal(err)
		}
	}
}

func (r *scriptRun) secretPath() string {
	return filepath.Join(r.sharedDir, "work-item", "secret")
}

func (r *scriptRun) waitForSecret(t *testing.T) {
	t.Helper()
	r.waitFor(t, "fetch.sh never wrote the work item", func() bool { return exists(r.secretPath()) })
}

func (r *scriptRun) waitForCollected(t *testing.T) {
	t.Helper()
	r.waitFor(t, "fetch never wrote the collected marker on its own volume", func() bool {
		return exists(filepath.Join(r.fetchDir, "collected"))
	})
}

// requestRecord asks for the run record as the worker does and waits for
// fetch to write it into the work-item directory.
func (r *scriptRun) requestRecord(t *testing.T) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(r.sharedDir, "worker-starting"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(r.sharedDir, "work-item", "worker-ran")
	r.waitFor(t, "fetch never recorded the run", func() bool { return exists(marker) })
}

func (r *scriptRun) waitFor(t *testing.T, failure string, done func() bool) {
	t.Helper()
	for !done() {
		if r.ctx.Err() != nil {
			t.Fatal(failure)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// loggedPaths returns the paths a stub appended to its log, none when the
// stub never ran.
func (r *scriptRun) loggedPaths(t *testing.T, log string) []string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(r.binDir, log))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	return strings.Fields(string(data))
}

func (r *scriptRun) stderr(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(r.stderrPath)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// start runs sh on script from this directory, with binDir first on PATH,
// stderr going to a file, and the test's time limit as the context. The
// test's own ANTHROPIC_ and CMA_ variables are left out. The scripts see
// only the values the test passes. An inherited variable would also hide a
// script that forgets to export its own value of it.
// The process group is killed when the test ends, which takes a sleeping
// background child down with the script.
func (r *scriptRun) start(t *testing.T, script string, env ...string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	t.Cleanup(cancel)
	r.ctx = ctx
	r.stderrPath = filepath.Join(r.binDir, "stderr")
	stderr, err := os.Create(r.stderrPath)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.CommandContext(ctx, "sh", script)
	for _, kv := range os.Environ() {
		if !strings.HasPrefix(kv, "ANTHROPIC_") && !strings.HasPrefix(kv, "CMA_") {
			cmd.Env = append(cmd.Env, kv)
		}
	}
	cmd.Env = append(cmd.Env, "PATH="+r.binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	cmd.Env = append(cmd.Env, env...)
	cmd.Stderr = stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Credential: r.credential}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	r.pid = cmd.Process.Pid
	exited := make(chan struct{})
	r.exited = exited
	go func() {
		r.waitErr = cmd.Wait()
		close(exited)
	}()
	t.Cleanup(func() {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		<-exited
		_ = stderr.Close()
	})
}

func writeStub(t *testing.T, binDir, name, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(binDir, name), []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// under reports whether path is dir or inside it.
func under(path, dir string) bool {
	rel, err := filepath.Rel(dir, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, "../")
}

// completeWorkItem is the dispatcher's answer as the stock curl stub gives
// it: every field the scripts read.
var completeWorkItem = `{"session_id":"s_1","work_id":"w_1","environment_id":"e_1","base_url":"https://api","secret":"` + packedSecret + `"}`

// Both scripts refuse to start with the environment key in their
// environment: the Anthropic credential must never reach a sandbox. fetch
// must refuse before it asks the dispatcher for anything.
func TestFetchRefusesToStartWithTheEnvironmentKey(t *testing.T) {
	run := newFetchRun(t)
	startFetchWith(t, run, completeWorkItem, "ANTHROPIC_ENVIRONMENT_KEY=key")

	if status := run.waitForExit(t, "assigned session"); status != 1 {
		t.Errorf("fetch exited %d, want 1", status)
	}
	if err := run.stderr(t); !strings.Contains(err, "refusing to start: ANTHROPIC_ENVIRONMENT_KEY is set") {
		t.Errorf("fetch did not refuse the environment key: %q", err)
	}
	if paths := run.loggedPaths(t, "curl-writes"); len(paths) != 0 {
		t.Errorf("fetch called the dispatcher with the environment key set: curl wrote %v", paths)
	}
	if exists(filepath.Join(run.sharedDir, "work-item")) {
		t.Error("fetch made the work-item directory with the environment key set")
	}
}

// The session's code can create entries of the shared volume. A symlink
// where fetch expects its own directory would point the ownership check,
// the chmod and every write at whatever the link names. fetch must refuse
// to start on a symlink at either directory it makes, leaving the target
// untouched.
func TestFetchRefusesToStartOnASymlink(t *testing.T) {
	for _, tc := range []struct {
		name string
		link func(run *scriptRun) string
	}{
		{name: "work-item directory", link: func(run *scriptRun) string { return filepath.Join(run.sharedDir, "work-item") }},
		{name: "home directory", link: func(run *scriptRun) string { return filepath.Join(run.fetchDir, "home") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			run := newFetchRun(t)
			// A mode fetch would change: it makes HOME 0700 and the
			// work-item directory 0750.
			target := filepath.Join(t.TempDir(), "target")
			if err := os.Mkdir(target, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(target, tc.link(run)); err != nil {
				t.Fatal(err)
			}
			before, err := os.Stat(target)
			if err != nil {
				t.Fatal(err)
			}
			startFetchWith(t, run, completeWorkItem)

			if status := run.waitForExit(t, "assigned session"); status != 1 {
				t.Errorf("fetch exited %d, want 1", status)
			}
			if err := run.stderr(t); !strings.Contains(err, "is a symlink") {
				t.Errorf("fetch did not refuse the symlink: %q", err)
			}
			after, err := os.Stat(target)
			if err != nil {
				t.Fatal(err)
			}
			if after.Mode() != before.Mode() {
				t.Errorf("fetch changed the link's target from %v to %v", before.Mode(), after.Mode())
			}
			entries, err := os.ReadDir(target)
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 0 {
				t.Errorf("fetch wrote %q through the link", entries[0].Name())
			}
		})
	}
}

// The worker takes the presence of the secret as the whole work item
// being there, so fetch must put the secret in place after every other
// field. The mv stub records each file as it is put in place.
func TestFetchPutsTheSecretInPlaceLast(t *testing.T) {
	run := newFetchRun(t)
	startFetchWith(t, run, completeWorkItem)
	run.waitForSecret(t)

	var fields []string
	for _, p := range run.loggedPaths(t, "mv-targets") {
		if under(p, filepath.Join(run.sharedDir, "work-item")) {
			fields = append(fields, filepath.Base(p))
		}
	}
	if len(fields) != 5 || fields[4] != "secret" {
		t.Errorf("fetch put the fields in place in the order %v; secret must be last of five", fields)
	}
}

// A work item missing a field is not one the worker can run. fetch must
// stop before it puts the secret in place, since the secret is what tells
// the worker the item is complete.
func TestFetchRefusesAWorkItemMissingAField(t *testing.T) {
	run := newFetchRun(t)
	startFetchWith(t, run, `{"session_id":"s_1","work_id":"w_1","environment_id":"e_1","secret":"`+packedSecret+`"}`)

	if status := run.waitForExit(t, "assigned session"); status != 1 {
		t.Errorf("fetch exited %d, want 1", status)
	}
	if err := run.stderr(t); !strings.Contains(err, "the dispatcher's work item is missing base_url") {
		t.Errorf("fetch did not report the missing field: %q", err)
	}
	if exists(run.secretPath()) {
		t.Error("fetch put the secret in place for a work item missing a field")
	}
}

// newFetchRun makes the volumes for a run of fetch.sh, so a test can lay
// them out before startFetchWith starts the script.
func newFetchRun(t *testing.T) *scriptRun {
	t.Helper()
	for _, tool := range []string{"sh", "jq"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s is not installed", tool)
		}
	}
	return &scriptRun{sharedDir: t.TempDir(), fetchDir: t.TempDir(), binDir: t.TempDir()}
}

// startFetchWith runs fetch.sh on run's volumes with env added to its
// environment. Its curl answers every request 200 with workItem and
// records the file it wrote.
func startFetchWith(t *testing.T, run *scriptRun, workItem string, env ...string) {
	t.Helper()
	startFetchWithCurl(t, run, `#!/bin/sh
out=""
while [ $# -gt 0 ]; do
	case "$1" in
	-o) out="$2"; shift 2 ;;
	-H) echo "$2" >> "`+filepath.Join(run.binDir, "curl-headers")+`"; shift 2 ;;
	*) shift ;;
	esac
done
echo "$out" >> "`+filepath.Join(run.binDir, "curl-writes")+`"
printf '%s' '`+workItem+`' > "$out"
printf 200
`, env...)
}

// startFetchWithCurl runs fetch.sh on run's volumes with curl as the body
// of its curl stub and env added to its environment. Unless the test
// installed an mv of its own, its mv records every destination. Its token
// file, which also serves as the CA file, holds "sa".
func startFetchWithCurl(t *testing.T, run *scriptRun, curl string, env ...string) {
	t.Helper()
	writeStub(t, run.binDir, "curl", curl)
	if !exists(filepath.Join(run.binDir, "mv")) {
		realMV, err := exec.LookPath("mv")
		if err != nil {
			t.Fatal(err)
		}
		writeStub(t, run.binDir, "mv", `#!/bin/sh
for last; do :; done
echo "$last" >> "`+filepath.Join(run.binDir, "mv-targets")+`"
exec "`+realMV+`" "$@"
`)
	}
	tokenFile := filepath.Join(run.binDir, "sa-token")
	if err := os.WriteFile(tokenFile, []byte("sa"), 0o600); err != nil {
		t.Fatal(err)
	}
	run.start(t, "fetch.sh", append([]string{
		"CMA_DISPATCHER_URL=https://dispatcher",
		"CMA_SA_TOKEN_FILE=" + tokenFile,
		"CMA_CA_FILE=" + tokenFile,
		"CMA_SHARED_DIR=" + run.sharedDir,
		"CMA_FETCH_DIR=" + run.fetchDir,
	}, env...)...)
}

// waitForExit waits for the script to exit on its own and returns its exit
// status. A script that prints gotPast instead is past the refusal under
// test and fails the test at once, rather than at the test's time limit.
func (r *scriptRun) waitForExit(t *testing.T, gotPast string) int {
	t.Helper()
	r.waitFor(t, "the script kept running", func() bool {
		if strings.Contains(r.stderr(t), gotPast) {
			t.Fatalf("the script got past the refusal: %s", r.stderr(t))
		}
		select {
		case <-r.exited:
			return true
		default:
			return false
		}
	})
	return r.exitStatus(t)
}

// exitStatus returns the exit status of a script that has exited.
func (r *scriptRun) exitStatus(t *testing.T) int {
	t.Helper()
	if r.waitErr == nil {
		return 0
	}
	exit, ok := r.waitErr.(*exec.ExitError)
	if !ok {
		t.Fatal(r.waitErr)
	}
	return exit.ExitCode()
}

// A dispatcher request that fails costs fetch one poll, which the backoff
// retries, whether curl could not connect or the dispatcher took the
// connection and never answered. For the second case fetch must give curl a
// time limit no longer than the 10 s after which the dispatcher stops
// writing an answer. The stub gives up at once under such a limit, as curl
// would at the limit; without one it waits past the test's own deadline, as
// a real curl would wait for good.
func TestFetchRetriesADispatcherRequestThatFails(t *testing.T) {
	for _, tc := range []struct {
		name      string
		firstCall string // what the curl stub does on its first call
	}{
		{name: "curl cannot connect", firstCall: "printf 000; exit 7"},
		{name: "the dispatcher never answers", firstCall: `[ -n "$limit" ] && [ "$limit" -le 10 ] || sleep 60; printf 000; exit 28`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			run := newFetchRun(t)
			failed := filepath.Join(run.binDir, "failed")
			startFetchWithCurl(t, run, `#!/bin/sh
out=""
limit=""
while [ $# -gt 0 ]; do
	case "$1" in
	-o) out="$2"; shift 2 ;;
	--max-time) limit="$2"; shift 2 ;;
	*) shift ;;
	esac
done
echo "$out" >> "`+filepath.Join(run.binDir, "curl-writes")+`"
if [ ! -e "`+failed+`" ]; then
	: > "`+failed+`"
	`+tc.firstCall+`
fi
printf '%s' '`+completeWorkItem+`' > "$out"
printf 200
`)
			run.waitFor(t, "fetch never got past the failed request", func() bool {
				select {
				case <-run.exited:
					if run.ctx.Err() == nil {
						t.Fatalf("fetch exited on the failed request: %s", run.stderr(t))
					}
				default:
				}
				return exists(run.secretPath())
			})
			if err := run.stderr(t); !strings.Contains(err, "dispatcher answered 000 to GET /work-item; retrying") {
				t.Errorf("fetch did not report the failed request: %q", err)
			}
			if calls := run.loggedPaths(t, "curl-writes"); len(calls) != 2 {
				t.Errorf("curl ran %d times; want the failed request and one retry", len(calls))
			}
		})
	}
}

// listedRefusals returns the reasons that refusal_reason of fetch.sh lists.
func listedRefusals(t *testing.T) []string {
	t.Helper()
	script, err := os.ReadFile("fetch.sh")
	if err != nil {
		t.Fatal(err)
	}
	_, body, found := strings.Cut(string(script), "\nrefusal_reason() {\n")
	body, _, closed := strings.Cut(body, "\n}\n")
	if !found || !closed {
		t.Fatal("fetch.sh has no refusal_reason function")
	}
	var listed []string
	for _, m := range regexp.MustCompile(`(?m)^\s*'([^']+)'( \| \\|\))$`).FindAllStringSubmatch(body, -1) {
		listed = append(listed, m[1])
	}
	if len(listed) == 0 {
		t.Fatal("refusal_reason lists no reason, one to a line in single quotes")
	}
	return listed
}

// The dispatcher answers 403 both to a pod it has not yet seen run, which the
// next polls get past, and to a caller it will never serve. Its answer says
// which, so fetch logs the reason and polls on as it does after any other
// answer. It logs the answer's first line only when that is a reason the script
// lists, and a phrase of its own otherwise. So the line that reports a 403 holds
// nothing that the answer's writer chose, the pod's credential in any form
// included. This test says nothing of the script's other lines. The dispatcher's
// tests hold the list to what the dispatcher writes.
//
// The curl stub answers each call 403 with the next of answers, then 200 with
// the work item. For an empty answer it leaves the file alone. An earlier
// answer would then be read again, if the script counted on curl to clear it.
func TestFetchLogsTheReasonOfA403AndPollsOn(t *testing.T) {
	const (
		nonce   = "Zk3vX9qLr0aTn5uYb7cWd1eHf2gJi4kMl6oPp8sQtUw"
		token   = "eyJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6Y21hOnNhbmRib3gifQ.c2lnbmF0dXJl"
		noPod   = "token is not bound to a pod: the sandbox has no running pod"
		unknown = "a reason this script does not know"
	)
	logged := func(reason string, delay int) string {
		return fmt.Sprintf("dispatcher answered 403 to GET /work-item (%s); retrying in %ds\n", reason, delay)
	}
	type refusal struct {
		name    string
		byToken bool
		answers []string
		want    string // what fetch logs ahead of the line on the work item
	}
	cases := []refusal{
		{name: "a listed reason and more lines", answers: []string{noPod + "\nsecond line " + nonce + "\n"}, want: logged(noPod, 2)},
		{name: "a listed reason with no newline", answers: []string{noPod}, want: logged(noPod, 2)},
		{name: "an answer with no body after one with a body", answers: []string{noPod + "\n", ""}, want: logged(noPod, 2) + logged("no reason given", 4)},
		{name: "a reason that is not listed", answers: []string{"the sandbox is on hold\n"}, want: logged(unknown, 2)},
		{name: "a listed reason and more on its line", answers: []string{noPod + " " + nonce + "\n"}, want: logged(unknown, 2)},
		{name: "a listed reason behind an empty line", answers: []string{"\n" + noPod + "\n"}, want: logged(unknown, 2)},
		{name: "a listed reason that ends in a carriage return", answers: []string{noPod + "\r\n"}, want: logged(unknown, 2)},
		{name: "a very long first line", answers: []string{noPod + strings.Repeat("x", 1<<20) + "\n"}, want: logged(unknown, 2)},
		{name: "control bytes and shell syntax", answers: []string{"re\x1b[31mfused\r $(id) `id` *\n"}, want: logged(unknown, 2)},
		{name: "the nonce", answers: []string{"no pod holds " + nonce + "\n"}, want: logged(unknown, 2)},
		{name: "the nonce with a mark after each character", answers: []string{strings.Join(strings.Split(nonce, ""), "!") + "\n"}, want: logged(unknown, 2)},
		{name: "the nonce with a space after each character", answers: []string{strings.Join(strings.Split(nonce, ""), " ") + "\n"}, want: logged(unknown, 2)},
		{name: "the token", byToken: true, answers: []string{"not for Bearer " + token + "\n"}, want: logged(unknown, 2)},
	}
	for _, reason := range listedRefusals(t) {
		cases = append(cases, refusal{name: "listed: " + reason, answers: []string{reason + "\n"}, want: logged(reason, 2)})
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			run := newFetchRun(t)
			for i, answer := range tc.answers {
				if err := os.WriteFile(filepath.Join(run.binDir, fmt.Sprintf("answer-%d", i+1)), []byte(answer), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			realSleep, err := exec.LookPath("sleep")
			if err != nil {
				t.Fatal(err)
			}
			// The waits between these polls are logged and not waited out. Every
			// other wait is real, or the script's later loops would spin.
			writeStub(t, run.binDir, "sleep", `#!/bin/sh
echo "$1" >> "`+filepath.Join(run.binDir, "sleeps")+`"
case "$1" in 2 | 4) exit 0 ;; esac
exec "`+realSleep+`" "$@"
`)
			credential := "CMA_ASSIGNMENT_NONCE=" + nonce
			if tc.byToken {
				tokenFile := filepath.Join(run.binDir, "long-sa-token")
				if err := os.WriteFile(tokenFile, []byte(token), 0o600); err != nil {
					t.Fatal(err)
				}
				credential = "CMA_SA_TOKEN_FILE=" + tokenFile
			}
			startFetchWithCurl(t, run, `#!/bin/sh
out=""
while [ $# -gt 0 ]; do
	case "$1" in
	-o) out="$2"; shift 2 ;;
	*) shift ;;
	esac
done
echo "$out" >> "`+filepath.Join(run.binDir, "curl-writes")+`"
answer="`+run.binDir+`/answer-$(wc -l < "`+filepath.Join(run.binDir, "curl-writes")+`" | tr -d ' ')"
if [ -e "$answer" ]; then
	[ ! -s "$answer" ] || cat "$answer" > "$out"
	printf 403
	exit 0
fi
printf '%s' '`+completeWorkItem+`' > "$out"
printf 200
`, credential)
			// fetch logs the work item after it has put the secret in place.
			run.waitFor(t, "fetch never logged the work item", func() bool { return strings.Contains(run.stderr(t), "assigned session ") })
			if got, _, _ := strings.Cut(run.stderr(t), "assigned session "); got != tc.want {
				t.Errorf("ahead of the work item fetch logged\n%q, want\n%q", got, tc.want)
			}
			if calls := run.loggedPaths(t, "curl-writes"); len(calls) != len(tc.answers)+1 {
				t.Errorf("curl ran %d times, want once for each 403 and once more", len(calls))
			}
			if waits, want := run.loggedLines(t, "sleeps"), []string{"2", "4"}[:len(tc.answers)]; len(waits) < len(want) || !slices.Equal(waits[:len(want)], want) {
				t.Errorf("fetch waited %v between its polls, want %v first", waits, want)
			}
		})
	}
}

// A 410 means the pod's worker has already exited. fetch must stop polling
// on it, writing no work item, rather than retry a request the dispatcher
// will never answer differently. Reaching the hourly idle sleep shows it
// left the poll loop; the curl stub logs each call.
func TestFetchStopsPollingWhenTheWorkItemIsClosed(t *testing.T) {
	run := newFetchRun(t)
	logSleeps(t, run)
	startFetchWithCurl(t, run, `#!/bin/sh
echo called >> "`+filepath.Join(run.binDir, "curl-calls")+`"
printf 410
`)
	run.waitFor(t, "fetch never idled after the closed work item", func() bool {
		if err := run.stderr(t); strings.Contains(err, "retrying") {
			t.Fatalf("fetch retried the closed work item: %s", err)
		}
		return run.sleeping(t, "3600")
	})
	if err := run.stderr(t); !strings.Contains(err, "this pod's work item is closed") {
		t.Errorf("fetch did not report the closed work item: %q", err)
	}
	if calls := run.loggedLines(t, "curl-calls"); len(calls) != 1 {
		t.Errorf("curl ran %d times; want once", len(calls))
	}
	if exists(run.secretPath()) {
		t.Error("fetch wrote a work item for a closed one")
	}
}

// fetch trusts the dispatcher's certificate against the operator's CA: the
// request carries the pod's credential and the answer the session's secret.
// The curl stub answers only a request that names the CA file fetch was
// given.
func TestFetchAsksTheDispatcherWithTheOperatorsCA(t *testing.T) {
	run := newFetchRun(t)
	startFetchWithCurl(t, run, `#!/bin/sh
out=""
ca=""
while [ $# -gt 0 ]; do
	case "$1" in
	-o) out="$2"; shift 2 ;;
	--cacert) ca="$2"; shift 2 ;;
	*) shift ;;
	esac
done
if [ "$ca" != "$CMA_CA_FILE" ]; then
	printf 000
	exit 60
fi
printf '%s' '`+completeWorkItem+`' > "$out"
printf 200
`)
	run.waitFor(t, "fetch never collected the work item", func() bool {
		err := run.stderr(t)
		if strings.Contains(err, "retrying") {
			t.Fatalf("fetch asked the dispatcher without naming the operator's CA file: %s", err)
		}
		return exists(run.secretPath())
	})
}

// Each value fetch writes is readable by the pod's group, which is how the
// worker reads it, and by nobody else. The secret is a credential, and a
// group-writable copy could be replaced in place by the session's code.
func TestFetchWritesEachValueReadableByTheGroupOnly(t *testing.T) {
	run := newFetchRun(t)
	startFetchWith(t, run, completeWorkItem)
	run.waitForSecret(t)

	for _, field := range []string{"session_id", "work_id", "environment_id", "base_url", "secret"} {
		info, err := os.Stat(filepath.Join(run.sharedDir, "work-item", field))
		if err != nil {
			t.Fatal(err)
		}
		if mode := info.Mode().Perm(); mode != 0o440 {
			t.Errorf("%s mode = %o; want 0440", field, mode)
		}
	}
}

// The worker removes its request once it sees the record, so a fetch
// restarted after it recorded the run finds no request to wait for. It must
// keep the record it wrote and go straight to idling: an hourly sleep,
// where the wait for a request sleeps a second at a time.
func TestFetchKeepsTheRecordItWroteBeforeARestart(t *testing.T) {
	run := newFetchRun(t)
	if err := os.WriteFile(filepath.Join(run.fetchDir, "collected"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	writeWorkItem(t, run.sharedDir)
	record := filepath.Join(run.sharedDir, "work-item", "worker-ran")
	if err := os.WriteFile(record, nil, 0o440); err != nil {
		t.Fatal(err)
	}
	logSleeps(t, run)
	startFetchWith(t, run, completeWorkItem)

	run.waitFor(t, "fetch never slept", func() bool { return len(run.loggedLines(t, "sleeps")) > 0 })
	if first := run.loggedLines(t, "sleeps")[0]; first != "3600" {
		t.Errorf("fetch's first sleep was %s s; a restarted fetch with the run recorded must idle, not wait for a request", first)
	}
	if !exists(record) {
		t.Error("fetch removed the record it wrote before its restart")
	}
}

// The collected marker tells a restarted fetch that the work item is in
// place, so it must not be written before the secret, the last value, is
// moved into place: a restart between the two would leave the pod with a
// marker, no secret, and no second request to the dispatcher. The mv stub
// holds the secret's move until the test releases it, so the moment before
// the secret is in place lasts long enough to look at.
func TestFetchMarksTheWorkItemCollectedOnlyOnceItIsWritten(t *testing.T) {
	run := newFetchRun(t)
	realMV, err := exec.LookPath("mv")
	if err != nil {
		t.Fatal(err)
	}
	blocked := filepath.Join(run.binDir, "mv-blocked")
	release := filepath.Join(run.binDir, "release")
	writeStub(t, run.binDir, "mv", `#!/bin/sh
for last; do :; done
case "$last" in
*/secret)
	: > "`+blocked+`"
	until [ -e "`+release+`" ]; do sleep 0.05; done
	;;
esac
exec "`+realMV+`" "$@"
`)
	startFetchWith(t, run, completeWorkItem)
	run.waitFor(t, "fetch never moved the secret into place", func() bool { return exists(blocked) })
	// Everything fetch does before the secret's move has happened by now.
	if exists(filepath.Join(run.fetchDir, "collected")) {
		t.Fatal("fetch marked the work item collected before the secret was in place")
	}
	if err := os.WriteFile(release, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	run.waitForSecret(t)
	run.waitForCollected(t)
}

// fetch is the container's first process, so the kubelet's stop signal
// reaches it alone and a trap runs only between commands. fetch spends its
// life sleeping, so it must sleep in a way the signal interrupts: a pod
// whose fetch held the signal off would be killed at the end of its grace
// period instead. The sleep log shows when the hourly idle sleep starts.
func TestFetchExitsOnAStopSignalWhileItIdles(t *testing.T) {
	run := newFetchRun(t)
	logSleeps(t, run)
	startFetchWith(t, run, completeWorkItem)
	run.waitForSecret(t)
	run.requestRecord(t)
	run.waitFor(t, "fetch never started its idle sleep", func() bool { return run.sleeping(t, "3600") })

	if err := syscall.Kill(run.pid, syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	select {
	case <-run.exited:
	case <-time.After(3 * time.Second):
		t.Fatal("fetch was still running 3 s after the stop signal")
	}
	if status := run.exitStatus(t); status != 143 {
		t.Errorf("fetch exited %d on the stop signal, want 143", status)
	}
}

// logSleeps installs a sleep that logs each sleep's length to the sleeps
// log before sleeping. The scripts find sleep through PATH, so the log
// shows which wait a script is in.
func logSleeps(t *testing.T, run *scriptRun) {
	t.Helper()
	realSleep, err := exec.LookPath("sleep")
	if err != nil {
		t.Fatal(err)
	}
	writeStub(t, run.binDir, "sleep", `#!/bin/sh
echo "$1" >> "`+filepath.Join(run.binDir, "sleeps")+`"
exec "`+realSleep+`" "$@"
`)
}

// sleeping reports whether a script has started a sleep of the given length,
// as logged by the stub logSleeps installs.
func (r *scriptRun) sleeping(t *testing.T, length string) bool {
	t.Helper()
	for _, logged := range r.loggedLines(t, "sleeps") {
		if logged == length {
			return true
		}
	}
	return false
}

// fetch must let curl verify the dispatcher's certificate against the
// operator's CA. A curl told to skip verification with -k or --insecure
// takes any certificate, so --cacert no longer decides who answers.
// Anything in the cluster that can answer https://dispatcher could then
// take the pod's credential and hand back a forged secret. The curl stub
// logs every argument of every call, one call per line.
func TestFetchNeverTellsCurlToSkipCertificateVerification(t *testing.T) {
	run := newFetchRun(t)
	startFetchWithCurl(t, run, `#!/bin/sh
printf '%s\t' "$@" >> "`+filepath.Join(run.binDir, "curl-args")+`"
echo >> "`+filepath.Join(run.binDir, "curl-args")+`"
out=""
while [ $# -gt 0 ]; do
	case "$1" in
	-o) out="$2"; shift 2 ;;
	*) shift ;;
	esac
done
printf '%s' '`+completeWorkItem+`' > "$out"
printf 200
`)
	run.waitForSecret(t)

	calls := run.loggedLines(t, "curl-args")
	if len(calls) == 0 {
		t.Fatal("curl-args logged nothing; the stub never ran")
	}
	// The options fetch uses that take the next argument as their value.
	// Extend it when fetch.sh's curl line gains an option that takes a value.
	takesValue := map[string]bool{"-o": true, "-w": true, "-H": true, "-X": true, "--cacert": true, "--max-time": true}
	for _, call := range calls {
		args := strings.Split(strings.TrimRight(call, "\t"), "\t")
		for i := 0; i < len(args); i++ {
			switch arg := args[i]; {
			case takesValue[arg]:
				i++
			case arg == "--insecure",
				strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && strings.Contains(arg, "k"):
				t.Errorf("fetch passed curl %q, which skips certificate verification and makes --cacert a no-op", arg)
			}
		}
	}
}

// fetch stages each value in a file that the pod's group can at most read
// and nobody else can open, whatever umask the container started with. The
// mode the file is created with is the mode it keeps until fetch sets the
// final 0440. An image or entrypoint that leaves the umask group-writable
// would otherwise let the session's code rewrite the secret in that moment.
// The test starts fetch under such a umask, and its chmod stub holds the
// staged secret's chmod until the test has looked at the file. The
// write_value helper stages all five values the same way. The secret, the
// last one written and the credential, stands in for them.
func TestFetchStagesTheSecretReadableByTheGroupOnlyUnderAGroupWritableUmask(t *testing.T) {
	run := newFetchRun(t)
	realChmod, err := exec.LookPath("chmod")
	if err != nil {
		t.Fatal(err)
	}
	blocked := filepath.Join(run.binDir, "chmod-blocked")
	release := filepath.Join(run.binDir, "release")
	writeStub(t, run.binDir, "chmod", `#!/bin/sh
for last; do :; done
case "$last" in
*/.secret.tmp)
	: > "`+blocked+`"
	until [ -e "`+release+`" ]; do sleep 0.05; done
	;;
esac
exec "`+realChmod+`" "$@"
`)
	// The umask is the process's, and fetch inherits it. The tests in this
	// package run one at a time, so no other test sees it.
	old := syscall.Umask(0o002)
	t.Cleanup(func() { syscall.Umask(old) })
	startFetchWith(t, run, completeWorkItem)

	run.waitFor(t, "fetch never set the staged secret's mode", func() bool { return exists(blocked) })
	info, err := os.Stat(filepath.Join(run.sharedDir, "work-item", ".secret.tmp"))
	if err != nil {
		t.Fatal(err)
	}
	if extra := info.Mode().Perm() &^ 0o640; extra != 0 {
		t.Errorf("the staged secret's mode is %o under a group-writable umask; mode bits %o beyond owner read-write and group read let the session's code rewrite it, or someone outside the group read it, before it is put in place", info.Mode().Perm(), extra)
	}
	if err := os.WriteFile(release, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	run.waitForSecret(t)
}

// fetch refuses to start on a work-item directory another user owns. It
// leaves the directory's mode as it was and writes nothing into it. The
// session's code runs as the worker's user and can make that directory.
func TestFetchRefusesAWorkItemDirectoryAnotherUserOwns(t *testing.T) {
	run := newFetchRun(t)
	workItem := filepath.Join(run.sharedDir, "work-item")
	if err := os.Mkdir(workItem, 0o755); err != nil {
		t.Fatal(err)
	}
	giveToWorker(t, workItem)
	before, err := os.Stat(workItem)
	if err != nil {
		t.Fatal(err)
	}
	startFetchWith(t, run, completeWorkItem)

	if status := run.waitForExit(t, "assigned session"); status != 1 {
		t.Errorf("fetch exited %d, want 1", status)
	}
	if err := run.stderr(t); !strings.Contains(err, "belongs to another user") {
		t.Errorf("fetch did not refuse the other user's directory: %q", err)
	}
	after, err := os.Stat(workItem)
	if err != nil {
		t.Fatal(err)
	}
	if after.Mode() != before.Mode() {
		t.Errorf("fetch changed the directory's mode from %v to %v", before.Mode(), after.Mode())
	}
	entries, err := os.ReadDir(workItem)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("fetch wrote %q into the other user's directory", entries[0].Name())
	}
}

// A restarted fetch treats a work-item directory another user owns as gone.
// It records no run there. A restarted fetch finds the directory by name. The
// session's code can rename fetch's directory away and make its own in its
// place, with a secret in it.
func TestFetchRecordsNoRunInAnotherUsersWorkItemAfterARestart(t *testing.T) {
	run := newFetchRun(t)
	if err := os.WriteFile(filepath.Join(run.fetchDir, "collected"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	writeWorkItem(t, run.sharedDir)
	giveToWorker(t, filepath.Join(run.sharedDir, "work-item"))
	// This is the worker's request. A fetch that trusts the directory
	// answers it by recording the run there.
	if err := os.WriteFile(filepath.Join(run.sharedDir, "worker-starting"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	startFetchWith(t, run, completeWorkItem)

	record := filepath.Join(run.sharedDir, "work-item", "worker-ran")
	run.waitFor(t, "fetch never reported the work item gone", func() bool {
		if exists(record) {
			t.Fatal("fetch recorded the run in another user's work-item directory")
		}
		return strings.Contains(run.stderr(t), "this pod's work item is gone; not collecting it again")
	})
}

// giveToWorker makes workerUID the owner of path. The session's code owns
// what it makes on the shared volume. Giving a file to another user usually
// takes root, so the test skips where the chown is not permitted.
//
// A test that already runs as workerUID skips too. Its chown succeeds but
// leaves path with the test's own user, so no other user owns it.
func giveToWorker(t *testing.T, path string) {
	t.Helper()
	if os.Geteuid() == workerUID {
		t.Skipf("cannot give %s to a second uid here: the test already runs as uid %d", path, workerUID)
	}
	if err := os.Chown(path, workerUID, os.Getgid()); err != nil {
		if errors.Is(err, syscall.EPERM) || errors.Is(err, syscall.EINVAL) {
			t.Skipf("cannot give %s to a second uid here (%v)", path, err)
		}
		t.Fatal(err)
	}
}
