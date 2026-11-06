//go:build unix

package image

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

// The record that the worker ran the session blocks a restarted worker only
// where the session's code cannot delete it. The one such place is the
// work-item directory, which fetch owns and the worker's user cannot write. A
// record in the root of the shared volume, which the session's code can
// write, must not block the worker. In this test one user owns everything, so
// a worker that gets past the record refuses the work item instead of running
// it. That refusal is how the test tells a blocked worker from one that got
// through.
func TestWorkerWaitsOnTheRunRecordFetchOwns(t *testing.T) {
	for _, tc := range []struct {
		name    string
		record  string
		blocked bool
	}{
		{name: "record in the work-item directory", record: filepath.Join("work-item", "worker-ran"), blocked: true},
		{name: "no record", blocked: false},
		{name: "record in the shared root", record: "worker-ran", blocked: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			run := prepareWorker(t)
			if tc.record != "" {
				if err := os.WriteFile(filepath.Join(run.sharedDir, tc.record), nil, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			startWorker(t, run)
			refusal := "refusing to run: this container owns or can write to"
			if !tc.blocked {
				run.waitFor(t, "the worker never reached the work item", func() bool {
					return strings.Contains(run.stderr(t), refusal)
				})
				return
			}
			// The worker prints nothing before the trust check, so any
			// output means it got past the record. Removing the record
			// then lets it through, which shows the record is what blocked
			// it. The wait is longer than the worker's 2 s poll.
			select {
			case <-run.exited:
				t.Errorf("the worker exited: %s", run.stderr(t))
			case <-time.After(3 * time.Second):
			}
			if err := run.stderr(t); err != "" {
				t.Errorf("the worker got past the run record: %s", err)
			}
			if err := os.Remove(filepath.Join(run.sharedDir, tc.record)); err != nil {
				t.Fatal(err)
			}
			run.waitFor(t, "the worker was still blocked after the record was removed", func() bool {
				return strings.Contains(run.stderr(t), refusal)
			})
		})
	}
}

// fetch puts the secret in place last, so a work item with every other field
// but no secret is still being written. A worker that starts before fetch
// finishes can find one. The worker must wait, with nothing printed, until
// the secret is in place, and go on to the work item once it is. The sleep
// log shows the worker in its two-second wait for the work item, the only
// wait of that length in the script.
func TestWorkerWaitsUntilTheSecretIsInPlace(t *testing.T) {
	run := prepareWorker(t)
	if err := os.Remove(run.secretPath()); err != nil {
		t.Fatal(err)
	}
	logSleeps(t, run)
	startWorker(t, run)

	run.waitFor(t, "the worker never waited for the secret", func() bool {
		if err := run.stderr(t); err != "" {
			t.Fatalf("the worker got past the work item without its secret: %s", err)
		}
		return run.sleeping(t, "2")
	})
	if err := os.WriteFile(run.secretPath(), []byte(packedSecret), 0o440); err != nil {
		t.Fatal(err)
	}
	run.waitForWorkItemRefusal(t, "the worker was still held after the secret appeared")
}

// The worker hands the secret's path to ant, which reads it as the worker's
// user. A secret that is in place but not readable by that user must hold the
// worker as a missing one does, with nothing printed, until it can be read.
// The sleep log shows the wait, as in TestWorkerWaitsUntilTheSecretIsInPlace.
// A user that can read every file cannot tell an unreadable secret from a
// readable one. The test needs a user without that override, as the pod's
// worker user is.
func TestWorkerWaitsUntilTheSecretIsReadable(t *testing.T) {
	run := prepareWorker(t)
	if err := os.Chmod(run.secretPath(), 0o200); err != nil {
		t.Fatal(err)
	}
	const readOK = 0x4 // R_OK
	if syscall.Access(run.secretPath(), readOK) == nil {
		t.Skip("this user can read a file it made unreadable, so an unreadable secret cannot be told from a readable one")
	}
	logSleeps(t, run)
	startWorker(t, run)

	run.waitFor(t, "the worker never waited for the secret", func() bool {
		if err := run.stderr(t); err != "" {
			t.Fatalf("the worker got past the work item with a secret it cannot read: %s", err)
		}
		return run.sleeping(t, "2")
	})
	if err := os.Chmod(run.secretPath(), 0o440); err != nil {
		t.Fatal(err)
	}
	run.waitForWorkItemRefusal(t, "the worker was still held after the secret became readable")
}

// prepareWorker lays out a complete work item and an nc that finds the
// resolver reachable and every isolation probe unreachable. startWorker
// then runs the script, so a test can add files to the shared volume first.
func prepareWorker(t *testing.T) *scriptRun {
	t.Helper()
	resolv, err := os.ReadFile("/etc/resolv.conf")
	if err != nil || !namesResolver(resolv) {
		t.Skip("worker.sh reads the resolver from /etc/resolv.conf, which names none here")
	}
	run := &scriptRun{sharedDir: t.TempDir(), binDir: t.TempDir()}
	writeStub(t, run.binDir, "nc", `#!/bin/sh
for last; do :; done
[ "$last" = 53 ]
`)
	writeWorkItem(t, run.sharedDir)
	return run
}

func startWorker(t *testing.T, run *scriptRun, env ...string) {
	t.Helper()
	run.start(t, "worker.sh", append([]string{
		"KUBERNETES_SERVICE_HOST=10.0.0.1",
		"KUBERNETES_SERVICE_PORT=443",
		"CMA_SHARED_DIR=" + run.sharedDir,
	}, env...)...)
}

// namesResolver reports whether worker.sh finds a resolver in resolv, the
// contents of /etc/resolv.conf. The worker reads the second field of the
// first line that starts with "nameserver". A commented or indented line
// does not match.
func namesResolver(resolv []byte) bool {
	for _, line := range strings.Split(string(resolv), "\n") {
		if strings.HasPrefix(line, "nameserver") {
			return len(strings.Fields(line)) >= 2
		}
	}
	return false
}

// The skip in prepareWorker must agree with the worker's own read of
// /etc/resolv.conf. A file whose only nameserver line is commented out
// passes a substring match. The suite then runs a worker that finds no
// resolver and waits for one instead of skipping.
func TestNamesResolverReadsAsTheWorkerDoes(t *testing.T) {
	for _, tc := range []struct {
		name   string
		resolv string
		want   bool
	}{
		{name: "a nameserver line", resolv: "search example.test\nnameserver 10.0.0.2\n", want: true},
		{name: "only a commented nameserver", resolv: "# nameserver 10.0.0.2\n", want: false},
		{name: "an indented nameserver", resolv: " nameserver 10.0.0.2\n", want: false},
		{name: "an address-less line before a nameserver line", resolv: "nameserver\nnameserver 10.0.0.2\n", want: false},
		{name: "an empty file", resolv: "", want: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := namesResolver([]byte(tc.resolv)); got != tc.want {
				t.Fatalf("namesResolver(%q) = %v, want %v", tc.resolv, got, tc.want)
			}
		})
	}
}

// The worker refuses to start with the environment key in its environment:
// the Anthropic credential must never reach a sandbox.
func TestWorkerRefusesToStartWithTheEnvironmentKey(t *testing.T) {
	run := prepareWorker(t)
	startWorker(t, run, "ANTHROPIC_ENVIRONMENT_KEY=key")

	if status := run.waitForExit(t, "running session"); status != 1 {
		t.Errorf("the worker exited %d, want 1", status)
	}
	if err := run.stderr(t); !strings.Contains(err, "refusing to start: ANTHROPIC_ENVIRONMENT_KEY is set") {
		t.Errorf("the worker did not refuse the environment key: %q", err)
	}
}

// The worker runs no session code while the pod can still reach an address
// the network policy must block. A pod that can connect to the metadata
// server holds the worker in its isolation check, with nothing printed.
// Once the address is unreachable the worker goes on to the work item.
func TestWorkerWaitsUntilTheProbesAreUnreachable(t *testing.T) {
	run := prepareWorker(t)
	// The nc stub finds every probe reachable while this file exists.
	reachable := filepath.Join(run.binDir, "reachable")
	if err := os.WriteFile(reachable, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	writeStub(t, run.binDir, "nc", `#!/bin/sh
for last; do :; done
[ "$last" = 53 ] || [ -e "`+reachable+`" ]
`)
	startWorker(t, run)

	// The worker's sweep runs at once and again every 5 s; 3 s is enough
	// for the first sweep to have found the probe reachable.
	select {
	case <-run.exited:
		t.Fatalf("the worker exited: %s", run.stderr(t))
	case <-time.After(3 * time.Second):
	}
	if err := run.stderr(t); err != "" {
		t.Fatalf("the worker got past the isolation check with a probe reachable: %s", err)
	}
	if err := os.Remove(reachable); err != nil {
		t.Fatal(err)
	}
	run.waitForWorkItemRefusal(t, "the worker was still held after every probe became unreachable")
}

// worker.sh gives each nc probe a connect timeout of at most one second. The
// isolation probes target addresses the network policy drops. A dropped
// connection gets no refusal. nc with no connect timeout then waits out the
// kernel's retries, about two minutes per probe. A correctly isolated pod
// makes up to three such connections on its one sweep. Its worker would then
// start minutes late. A pod still waiting counts each sweep as 5 s toward the
// 60 s deadline, whatever the probes took. At one second per probe that count
// stays close to the clock. The nc stub answers at once only under a timeout
// of one second, as nc gives up at its deadline. A timeout of 0 reads as none
// to nc, so the stub treats 0 as no timeout. Without one it waits past the
// test's own deadline, as nc would wait on a dropped address.
func TestWorkerGivesEachProbeAOneSecondConnectTimeout(t *testing.T) {
	run := prepareWorker(t)
	// The resolver is reachable and every probe unreachable, as in
	// prepareWorker's stub.
	writeStub(t, run.binDir, "nc", `#!/bin/sh
limit=""
last=""
while [ $# -gt 0 ]; do
	case "$1" in
	-w) limit="$2"; shift 2 ;;
	*) last="$1"; shift ;;
	esac
done
[ -n "$limit" ] && [ "$limit" -ge 1 ] && [ "$limit" -le 1 ] || sleep 60
[ "$last" = 53 ]
`)
	startWorker(t, run)

	run.waitForWorkItemRefusal(t, "the worker never got past its isolation check: an nc probe ran without a connect timeout of at most 1 s")
}

// A work-item directory the worker owns is one the session's code could
// have made, whatever its mode, so the worker must refuse it even when it
// cannot write there. A user that can write every directory cannot tell
// this arm of the check from the other. The test needs a user without
// that override, as the pod's worker user is.
func TestWorkerRefusesAWorkItemDirectoryItOwns(t *testing.T) {
	run := prepareWorker(t)
	workItem := filepath.Join(run.sharedDir, "work-item")
	if err := os.Chmod(workItem, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(workItem, 0o750) })
	const writeOK = 0x2 // W_OK
	if syscall.Access(workItem, writeOK) == nil {
		t.Skip("this user can write a directory it made read-only, so owning and writing cannot be told apart")
	}
	startWorker(t, run)

	run.waitForWorkItemRefusal(t, "the worker never reached the work item")
}

// waitForWorkItemRefusal waits for the worker to refuse the work item on
// the trust check. In these tests one user owns everything, so a worker that
// gets that far refuses the work item instead of running it. A worker that
// reports the session running got past the check, which fails the test at
// once.
func (r *scriptRun) waitForWorkItemRefusal(t *testing.T, failure string) {
	t.Helper()
	r.waitFor(t, failure, func() bool {
		err := r.stderr(t)
		if strings.Contains(err, "running session") {
			t.Fatalf("the worker ran the work item: %s", err)
		}
		return strings.Contains(err, "refusing to run: this container owns or can write to")
	})
}

// A pod whose network is not up yet fails every connection alike, the
// isolation probes included, which proves nothing about its network policy.
// The worker must hold its isolation check until the pod reaches its DNS
// resolver, with nothing printed, and go on to the work item once it does.
func TestWorkerWaitsUntilTheResolverIsReachable(t *testing.T) {
	run := prepareWorker(t)
	// The nc stub finds the resolver unreachable while this file exists.
	down := filepath.Join(run.binDir, "resolver-down")
	if err := os.WriteFile(down, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	writeStub(t, run.binDir, "nc", `#!/bin/sh
for last; do :; done
[ "$last" = 53 ] && [ ! -e "`+down+`" ]
`)
	startWorker(t, run)

	// The worker's sweep runs at once and again every 5 s; 3 s is enough
	// for the first sweep to have found the resolver unreachable.
	select {
	case <-run.exited:
		t.Fatalf("the worker exited: %s", run.stderr(t))
	case <-time.After(3 * time.Second):
	}
	if err := run.stderr(t); err != "" {
		t.Fatalf("the worker got past the isolation check without reaching the resolver: %s", err)
	}
	if err := os.Remove(down); err != nil {
		t.Fatal(err)
	}
	run.waitForWorkItemRefusal(t, "the worker was still held after the resolver became reachable")
}

// The worker must find the resolver reachable before the probes on every
// sweep, not once before its first. A pod can lose its network after one
// sweep and before the next one's probes. Those probes then fail to connect,
// which proves nothing. A resolver that returns after them would pass a
// sweep that exercised no network policy.
//
// In this test the resolver answers every other call. A worker that checks
// it before the probes finds it reachable there and unreachable after them
// on every sweep. So it never passes. At its deadline it refuses to start,
// naming the lost resolver. A worker that checked the resolver once before
// its sweeps would find it reachable after the probes on its second sweep.
// It would then go on to the work item. The sleep stub shortens the worker's
// five-second waits so the deadline passes within the test.
func TestWorkerChecksTheResolverBeforeTheProbesOnEverySweep(t *testing.T) {
	run := prepareWorker(t)
	// Every probe is unreachable. The resolver answers every other call,
	// counting from a reachable first; the stub keeps the count in the file
	// below.
	calls := filepath.Join(run.binDir, "resolver-calls")
	writeStub(t, run.binDir, "nc", `#!/bin/sh
# The worker's last argument is the port: 53 is the resolver, anything else a probe.
for last; do :; done
[ "$last" = 53 ] || exit 1
n=$(($(cat "`+calls+`" 2>/dev/null || echo 0) + 1))
echo "$n" > "`+calls+`"
[ $((n % 2)) = 1 ]
`)
	realSleep, err := exec.LookPath("sleep")
	if err != nil {
		t.Fatal(err)
	}
	writeStub(t, run.binDir, "sleep", `#!/bin/sh
exec "`+realSleep+`" 0.02
`)
	startWorker(t, run)

	if status := run.waitForExit(t, "refusing to run"); status != 1 {
		t.Errorf("the worker exited %d, want 1", status)
	}
	if err := run.stderr(t); !strings.Contains(err, "refusing to start: this pod lost its DNS resolver during the check after 60s") {
		t.Errorf("the worker did not hold every sweep on a resolver lost after the probes: %q", err)
	}
}

// The metadata server is the address the isolation check exists to catch: a
// pod that reaches it leaks the cloud's credentials to the session. The worker
// must stay in its check while that address alone answers, with nothing
// printed, and go on once it stops answering. The nc stub answers by address,
// so a probe list without the metadata server lets the worker through.
func TestWorkerWaitsWhileTheMetadataServerIsReachable(t *testing.T) {
	run := prepareWorker(t)
	// The nc stub finds the metadata server reachable while this file exists,
	// and every other probe unreachable throughout.
	reachable := filepath.Join(run.binDir, "metadata-reachable")
	if err := os.WriteFile(reachable, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	writeStub(t, run.binDir, "nc", `#!/bin/sh
# The worker's last two arguments are the host and the port.
for arg; do host="${port:-}"; port="$arg"; done
[ "$port" = 53 ] || { [ "$host:$port" = 169.254.169.254:80 ] && [ -e "`+reachable+`" ]; }
`)
	startWorker(t, run)

	// The worker's sweep runs at once and again every 5 s; 3 s is enough
	// for the first sweep to have found the metadata server reachable.
	select {
	case <-run.exited:
		t.Fatalf("the worker exited: %s", run.stderr(t))
	case <-time.After(3 * time.Second):
	}
	if err := run.stderr(t); err != "" {
		t.Fatalf("the worker got past the isolation check with the metadata server reachable: %s", err)
	}
	if err := os.Remove(reachable); err != nil {
		t.Fatal(err)
	}
	run.waitForWorkItemRefusal(t, "the worker was still held after the metadata server became unreachable")
}

// workerUID is the pod's worker user, which owns nothing fetch writes. It
// need not exist here: nothing in worker.sh looks the user up.
const workerUID = 1001

// startWorkerAsAnotherUser runs worker.sh as a user that neither owns nor
// can write the work-item directory, like the pod's worker user. That gets
// the worker past its trust check and on to the run. Starting a process as
// another user takes a capability only root usually has; where the test
// lacks it, it skips. The user keeps the test's group, which is how it
// reads the work item, as the pod's group does.
func startWorkerAsAnotherUser(t *testing.T, run *scriptRun, env ...string) {
	t.Helper()
	credential := &syscall.Credential{Uid: workerUID, Gid: uint32(os.Getgid())}
	// The script changes back to this directory by its absolute path, so
	// the user must be able to reach it that way too.
	script, err := filepath.Abs("worker.sh")
	if err != nil {
		t.Fatal(err)
	}
	probe := exec.Command("sh", "-c", `test -r "$1"`, "sh", script)
	probe.SysProcAttr = &syscall.SysProcAttr{Credential: credential}
	if err := probe.Run(); err != nil {
		var exit *exec.ExitError
		if errors.Is(err, syscall.EPERM) || errors.Is(err, syscall.EINVAL) || errors.As(err, &exit) {
			t.Skipf("cannot run worker.sh under a second uid here (%v); the trust check in wait_for_secret refuses a directory the worker owns or can write", err)
		}
		t.Fatal(err)
	}
	// The test's temporary directories sit in a directory only the test's
	// user can enter. The shared volume takes the worker's request and the
	// stubs' directory their logs.
	for _, dir := range []string{run.sharedDir, run.binDir} {
		if err := os.Chmod(filepath.Dir(dir), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(dir, 0o770); err != nil {
			t.Fatal(err)
		}
	}
	// A private temporary directory keeps the user out however the test
	// opens its own directories.
	probe = exec.Command("sh", "-c", `test -w "$1"`, "sh", run.sharedDir)
	probe.SysProcAttr = &syscall.SysProcAttr{Credential: credential}
	if err := probe.Run(); err != nil {
		t.Skipf("the second uid cannot write the shared volume under %s (%v)", filepath.Dir(run.sharedDir), err)
	}
	run.credential = credential
	startWorker(t, run, env...)
}

// recordTheRun answers the worker's request as fetch does: once the worker
// asks, the record is written into the work-item directory.
func (r *scriptRun) recordTheRun(t *testing.T) {
	t.Helper()
	request := filepath.Join(r.sharedDir, "worker-starting")
	r.waitFor(t, "the worker never asked for its run to be recorded", func() bool { return exists(request) })
	if err := os.WriteFile(filepath.Join(r.sharedDir, "work-item", "worker-ran"), nil, 0o440); err != nil {
		t.Fatal(err)
	}
}

// The record that the worker ran is what keeps a restarted worker from
// running the session again. So a worker whose fetch never records the run
// must refuse to run it. The sleep stub shortens the worker's one-second
// waits so the deadline passes within the test.
func TestWorkerRefusesToRunWithoutTheRunRecord(t *testing.T) {
	run := prepareWorker(t)
	realSleep, err := exec.LookPath("sleep")
	if err != nil {
		t.Fatal(err)
	}
	writeStub(t, run.binDir, "sleep", `#!/bin/sh
exec "`+realSleep+`" 0.02
`)
	writeStub(t, run.binDir, "ant", `#!/bin/sh
echo ran >> "`+filepath.Join(run.binDir, "ant-log")+`"
`)
	startWorkerAsAnotherUser(t, run)

	if status := run.waitForExit(t, "worker for session"); status != 1 {
		t.Errorf("the worker exited %d, want 1", status)
	}
	if err := run.stderr(t); !strings.Contains(err, "refusing to run: fetch did not record the run within 60s") {
		t.Errorf("the worker did not refuse the unrecorded run: %q", err)
	}
	// The ant stub's >> creates its log as it runs. Any log at all, even an
	// empty one, means a run.
	if exists(filepath.Join(run.binDir, "ant-log")) {
		t.Error("the worker ran the session without the run record")
	}
}

// fetch polls for the worker's request every second, so it can record the
// run in the last second the worker waits. The worker must run a session
// whose record arrived then. The record is in place, which is all a
// restarted worker needs. The sleep stub shortens the worker's one-second
// waits. It holds the last one the deadline allows until the test has
// written the record, as fetch would have by then. A worker that decides by
// how long it waited, rather than by the record, refuses that run.
func TestWorkerRunsWhenTheRunIsRecordedInTheLastSecondOfItsWait(t *testing.T) {
	run := prepareWorker(t)
	realSleep, err := exec.LookPath("sleep")
	if err != nil {
		t.Fatal(err)
	}
	// The stub counts the worker's one-second waits in the file at waits.
	// On the 60th wait it creates the file at lastWait. The 60th wait is the
	// last one RECORD_DEADLINE_SECONDS in worker.sh allows. The stub runs as
	// the worker's user, which cannot write the work-item directory, so the
	// test writes the record instead.
	waits := filepath.Join(run.binDir, "one-second-waits")
	lastWait := filepath.Join(run.binDir, "last-wait")
	record := filepath.Join(run.sharedDir, "work-item", "worker-ran")
	writeStub(t, run.binDir, "sleep", `#!/bin/sh
[ "$1" = 1 ] || exec "`+realSleep+`" 0.02
n=$(($(cat "`+waits+`" 2>/dev/null || echo 0) + 1))
echo "$n" > "`+waits+`"
if [ "$n" = 60 ]; then
	: > "`+lastWait+`"
	until [ -e "`+record+`" ]; do "`+realSleep+`" 0.01; done
fi
exec "`+realSleep+`" 0.02
`)
	writeStub(t, run.binDir, "ant", `#!/bin/sh
echo ran >> "`+filepath.Join(run.binDir, "ant-log")+`"
`)
	startWorkerAsAnotherUser(t, run)

	run.waitFor(t, "the worker never reached the last second of its wait", func() bool {
		if err := run.stderr(t); strings.Contains(err, "refusing to run") {
			t.Fatalf("the worker gave up before the last second of its wait: %s", err)
		}
		return exists(lastWait)
	})
	if err := os.WriteFile(record, nil, 0o440); err != nil {
		t.Fatal(err)
	}
	run.waitFor(t, "the worker kept running after the run was recorded", func() bool {
		select {
		case <-run.exited:
			return true
		default:
			return false
		}
	})
	if status := run.exitStatus(t); status != 0 {
		t.Errorf("the worker exited %d on a run recorded in the last second of its wait, want 0: %q", status, run.stderr(t))
	}
	if lines := run.loggedLines(t, "ant-log"); len(lines) != 1 || lines[0] != "ran" {
		t.Errorf("the session's worker logged %v; want one run", lines)
	}
}

// The isolation check must end. A pod that never becomes isolated holds the
// worker in the check with nothing printed. Without a deadline the pod would
// sit there for good. Once the deadline passes the worker must refuse to
// start, naming what still blocked it. It must not go on to the work item on
// a pod that can still reach an address the policy must block. In these
// tests a worker that reaches the work item refuses it as one its user owns,
// which is how waitForExit tells the two apart. One row keeps the metadata
// server reachable throughout, the other keeps the resolver unreachable. The
// sleep stub shortens the worker's five-second waits so the deadline passes
// within the test.
func TestWorkerRefusesToStartWhenThePodNeverBecomesIsolated(t *testing.T) {
	for _, tc := range []struct {
		name string
		nc   string // the nc stub's body after the shebang
		want string
	}{
		{name: "metadata server reachable", nc: `# The worker's last two arguments are the host and the port.
for arg; do host="${port:-}"; port="$arg"; done
[ "$port" = 53 ] || [ "$host:$port" = 169.254.169.254:80 ]`, want: "refusing to start: this pod can still reach 169.254.169.254:80 after 60s"},
		{name: "resolver unreachable", nc: `exit 1`, want: "refusing to start: this pod cannot reach its DNS resolver after 60s"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			run := prepareWorker(t)
			writeStub(t, run.binDir, "nc", "#!/bin/sh\n"+tc.nc+"\n")
			realSleep, err := exec.LookPath("sleep")
			if err != nil {
				t.Fatal(err)
			}
			writeStub(t, run.binDir, "sleep", `#!/bin/sh
exec "`+realSleep+`" 0.02
`)
			startWorker(t, run)

			if status := run.waitForExit(t, "refusing to run"); status != 1 {
				t.Errorf("the worker exited %d, want 1", status)
			}
			if err := run.stderr(t); !strings.Contains(err, tc.want) {
				t.Errorf("the worker did not refuse the pod at its deadline: %q", err)
			}
		})
	}
}

// The worker runs the session's worker as a child it can stop. The pod's
// stop signal reaches the script alone, so the script must pass it on and
// wait. Otherwise the worker is killed at the end of the grace period with
// the session left open. The worker's exit status is the pod's, so a failed
// worker shows as a failed pod. The ant stub logs its start and the
// forwarded signal.
func TestWorkerForwardsTheStopSignalAndReturnsTheWorkerStatus(t *testing.T) {
	for _, tc := range []struct {
		name   string
		ant    string // the ant stub's body; it logs "started" once it can be stopped
		signal bool   // send the worker a stop signal once ant has started
		want   int
	}{
		{name: "stop signal", ant: `trap 'echo terminated >> "$log"; exit 0' TERM
echo started >> "$log"
sleep 3600 &
wait $!`, signal: true, want: 143},
		{name: "worker exit status", ant: `echo started >> "$log"
exit 7`, want: 7},
	} {
		t.Run(tc.name, func(t *testing.T) {
			run := prepareWorker(t)
			log := filepath.Join(run.binDir, "ant-log")
			writeStub(t, run.binDir, "ant", `#!/bin/sh
log="`+log+`"
`+tc.ant+`
`)
			startWorkerAsAnotherUser(t, run)
			run.recordTheRun(t)
			run.waitFor(t, "the worker never started the session's worker", func() bool {
				return len(run.loggedLines(t, "ant-log")) > 0
			})
			if tc.signal {
				if err := syscall.Kill(run.pid, syscall.SIGTERM); err != nil {
					t.Fatal(err)
				}
			}

			select {
			case <-run.exited:
			case <-time.After(3 * time.Second):
				t.Fatal("the worker was still running 3 s after the session's worker stopped")
			}
			if status := run.exitStatus(t); status != tc.want {
				t.Errorf("the worker exited %d, want %d", status, tc.want)
			}
			if tc.signal {
				if lines := run.loggedLines(t, "ant-log"); len(lines) != 2 || lines[1] != "terminated" {
					t.Errorf("the session's worker logged %v; want its start and the forwarded stop signal", lines)
				}
			} else if err := run.stderr(t); !strings.Contains(err, "exited with status 7") {
				t.Errorf("the worker did not report its worker's status: %q", err)
			}
		})
	}
}

// The isolation check catches more than the IPv4 metadata server. It also
// probes the Kubernetes API server, which every pod can reach unless the
// network policy blocks it, and the EC2 IPv6 metadata server. The worker must
// stay in its check while either of those alone answers, with nothing
// printed, and go on once it stops answering. The nc stub answers by address,
// so a probe list without that address lets the worker through.
func TestWorkerWaitsWhileTheAPIServerOrTheIPv6MetadataServerIsReachable(t *testing.T) {
	for _, tc := range []struct {
		name string
		// The probe's host and port as the worker hands them to nc. Each is
		// its own argument. The stub checks them apart.
		host, port string
		env        []string // the worker's environment, on top of startWorker's
	}{
		// The API server's address comes from the worker's environment and
		// differs by cluster. This row sets an address of its own for the
		// stub to answer.
		{name: "Kubernetes API server", host: "10.9.8.7", port: "26443", env: []string{"KUBERNETES_SERVICE_HOST=10.9.8.7", "KUBERNETES_SERVICE_PORT=26443"}},
		{name: "EC2 IPv6 metadata server", host: "fd00:ec2::254", port: "80"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			run := prepareWorker(t)
			// The nc stub finds this probe reachable while this file exists,
			// and every other probe unreachable throughout.
			reachable := filepath.Join(run.binDir, "probe-reachable")
			if err := os.WriteFile(reachable, nil, 0o644); err != nil {
				t.Fatal(err)
			}
			writeStub(t, run.binDir, "nc", `#!/bin/sh
# The worker's last two arguments are the host and the port.
for arg; do host="${port:-}"; port="$arg"; done
[ "$port" = 53 ] || { [ "$host" = "`+tc.host+`" ] && [ "$port" = "`+tc.port+`" ] && [ -e "`+reachable+`" ]; }
`)
			startWorker(t, run, tc.env...)

			// The worker's sweep runs at once and again every 5 s; 3 s is
			// enough for the first sweep to have found the probe reachable.
			select {
			case <-run.exited:
				t.Fatalf("the worker exited: %s", run.stderr(t))
			case <-time.After(3 * time.Second):
			}
			if err := run.stderr(t); err != "" {
				t.Fatalf("the worker got past the isolation check with the %s reachable: %s", tc.name, err)
			}
			if err := os.Remove(reachable); err != nil {
				t.Fatal(err)
			}
			run.waitForWorkItemRefusal(t, "the worker was still held after the "+tc.name+" became unreachable")
		})
	}
}

// A sweep can lose the network between its first resolver check and the
// probes. Every probe then fails to connect, which proves nothing. So the
// worker must hold its isolation check, with nothing printed, until a sweep
// finds the resolver reachable after the probes too.
func TestWorkerWaitsUntilTheResolverIsReachableAfterTheProbes(t *testing.T) {
	run := prepareWorker(t)
	// Every probe is unreachable, and the stub records that one ran. While
	// the resolver-flapping file exists, the resolver answers reachable,
	// except once after a probe has run since its last answer. With that
	// file gone, the resolver always answers.
	flapping := filepath.Join(run.binDir, "resolver-flapping")
	probed := filepath.Join(run.binDir, "probed")
	if err := os.WriteFile(flapping, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	writeStub(t, run.binDir, "nc", `#!/bin/sh
# The worker's last argument is the port: 53 is the resolver, anything else a probe.
for last; do :; done
if [ "$last" != 53 ]; then : > "`+probed+`"; exit 1; fi
[ -e "`+flapping+`" ] || exit 0
if [ -e "`+probed+`" ]; then
	rm "`+probed+`"
	exit 1
fi
`)
	startWorker(t, run)

	// The worker's sweep runs at once and again every 5 s; 3 s is enough
	// for the first sweep to have lost the resolver after the probes.
	select {
	case <-run.exited:
		t.Fatalf("the worker exited: %s", run.stderr(t))
	case <-time.After(3 * time.Second):
	}
	if err := run.stderr(t); err != "" {
		t.Fatalf("the worker got past the isolation check on a sweep that lost the resolver after the probes: %s", err)
	}
	if err := os.Remove(flapping); err != nil {
		t.Fatal(err)
	}
	run.waitForWorkItemRefusal(t, "the worker was still held after the resolver stayed reachable through a sweep")
}

// The worker must find the resolver reachable before it runs the probes. A
// sweep that probes first cannot tell a blocked address from a pod whose
// network is down. In this test the resolver answers only while a record that
// a probe ran exists, and each probe writes that record. A worker that checks
// the resolver first never reaches the probes. It must hold its isolation
// check with nothing printed. A worker that probed first would then find the
// resolver reachable and go on to the work item. Once the test writes the
// record itself, the resolver answers before the probes and the worker goes
// on to the work item.
func TestWorkerChecksTheResolverBeforeTheProbes(t *testing.T) {
	run := prepareWorker(t)
	// Every probe is unreachable, and the stub records that one ran. The
	// resolver answers reachable only while that record exists. The resolver
	// is the first nameserver in /etc/resolv.conf, as the worker reads it. A
	// port-53 call to any other host is unreachable.
	probed := filepath.Join(run.binDir, "probed")
	writeStub(t, run.binDir, "nc", `#!/bin/sh
# The worker's last two arguments are the host and the port.
for arg; do host="${port:-}"; port="$arg"; done
if [ "$port" != 53 ]; then : > "`+probed+`"; exit 1; fi
[ "$host" = "$(awk '/^nameserver/ { print $2; exit }' /etc/resolv.conf)" ] && [ -e "`+probed+`" ]
`)
	startWorker(t, run)

	// The worker's sweep runs at once and again every 5 s; 3 s is enough
	// for the first sweep to have found the resolver unreachable.
	select {
	case <-run.exited:
		t.Fatalf("the worker exited: %s", run.stderr(t))
	case <-time.After(3 * time.Second):
	}
	if err := run.stderr(t); err != "" {
		t.Fatalf("the worker got past the isolation check on a sweep that probed before checking the resolver: %s", err)
	}
	// With the record in place the resolver answers before the probes too.
	if err := os.WriteFile(probed, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	run.waitForWorkItemRefusal(t, "the worker was still held after the resolver became reachable before the probes")
}

// The worker refuses a work-item directory its user can write, even one
// fetch owns. That user also runs the session's code. That code could have
// written the fields. Here the directory is writable by its group. The
// worker's user shares that group.
func TestWorkerRefusesAWorkItemDirectoryItCanWrite(t *testing.T) {
	run := prepareWorker(t)
	if err := os.Chmod(filepath.Join(run.sharedDir, "work-item"), 0o770); err != nil {
		t.Fatal(err)
	}
	startWorkerAsAnotherUser(t, run)

	if status := run.waitForExit(t, "running session"); status != 1 {
		t.Errorf("the worker exited %d, want 1", status)
	}
	if err := run.stderr(t); !strings.Contains(err, "refusing to run: this container owns or can write to") {
		t.Errorf("the worker did not refuse a directory it can write: %q", err)
	}
}

// The worker refuses a work item with a field removed or left empty. It
// must not run the session under a blank identity.
func TestWorkerRefusesAWorkItemMissingAField(t *testing.T) {
	for _, field := range []string{"session_id", "work_id", "environment_id", "base_url"} {
		for _, empty := range []bool{false, true} {
			name := field + " removed"
			if empty {
				name = field + " empty"
			}
			t.Run(name, func(t *testing.T) {
				run := prepareWorker(t)
				path := filepath.Join(run.sharedDir, "work-item", field)
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				if empty {
					if err := os.WriteFile(path, nil, 0o440); err != nil {
						t.Fatal(err)
					}
				}
				startWorkerAsAnotherUser(t, run)

				if status := run.waitForExit(t, "running session"); status != 1 {
					t.Errorf("the worker exited %d, want 1", status)
				}
				if err := run.stderr(t); !strings.Contains(err, "the work item is missing "+field) {
					t.Errorf("the worker did not report the missing %s: %q", field, err)
				}
			})
		}
	}
}

// worker.sh starts ant with its fixed arguments, the secret's path and the
// work item's four identity fields. The ant stub logs its arguments and
// those five variables, one per line.
//
// The idle limit is 30m unless the pod sets CMA_MAX_IDLE, which the chart
// sets from sandbox.idleTimeout.
func TestWorkerStartsAntWithTheWorkItem(t *testing.T) {
	for _, tc := range []struct {
		name    string
		env     []string
		maxIdle string
	}{
		{name: "default idle limit", maxIdle: "30m"},
		{name: "idle limit from CMA_MAX_IDLE", env: []string{"CMA_MAX_IDLE=7m"}, maxIdle: "7m"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testWorkerStartsAnt(t, tc.env, tc.maxIdle)
		})
	}
}

func testWorkerStartsAnt(t *testing.T, env []string, maxIdle string) {
	run := prepareWorker(t)
	writeStub(t, run.binDir, "ant", `#!/bin/sh
{
	echo "$*"
	echo "ANTHROPIC_WORK_SECRET_FILE=${ANTHROPIC_WORK_SECRET_FILE-}"
	echo "ANTHROPIC_SESSION_ID=${ANTHROPIC_SESSION_ID-}"
	echo "ANTHROPIC_WORK_ID=${ANTHROPIC_WORK_ID-}"
	echo "ANTHROPIC_ENVIRONMENT_ID=${ANTHROPIC_ENVIRONMENT_ID-}"
	echo "ANTHROPIC_BASE_URL=${ANTHROPIC_BASE_URL-}"
} >> "`+filepath.Join(run.binDir, "ant-log")+`"
`)
	startWorkerAsAnotherUser(t, run, env...)
	run.recordTheRun(t)

	run.waitFor(t, "the worker kept running after ant exited", func() bool {
		select {
		case <-run.exited:
			return true
		default:
			return false
		}
	})
	if status := run.exitStatus(t); status != 0 {
		t.Errorf("the worker exited %d, want 0: %s", status, run.stderr(t))
	}
	want := []string{
		"beta:worker run --workdir /workspace --max-idle " + maxIdle,
		"ANTHROPIC_WORK_SECRET_FILE=" + run.secretPath(),
		"ANTHROPIC_SESSION_ID=s_1",
		"ANTHROPIC_WORK_ID=w_1",
		"ANTHROPIC_ENVIRONMENT_ID=e_1",
		"ANTHROPIC_BASE_URL=https://api",
	}
	if got := run.loggedLines(t, "ant-log"); strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("ant started with\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}
