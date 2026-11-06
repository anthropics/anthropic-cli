#!/bin/sh
# The worker container checks that the known-bad addresses are unreachable.
# It waits for the work item the fetch container leaves in the shared
# directory. Then it runs the CMA worker for that session. Nothing here
# talks to the dispatcher. POSIX sh so it runs under busybox.
set -eu

SHARED_DIR="${CMA_SHARED_DIR:-/run/cma}"
WORK_ITEM_DIR="$SHARED_DIR/work-item"
SECRET_FILE="$WORK_ITEM_DIR/secret"
# The worker runs one session per pod. fetch records the run in the
# work-item directory. The session's code runs as this container's user
# and cannot write that directory. A restarted worker waits on that
# record instead of running the session again. The dispatcher also
# suspends the sandbox when the worker container restarts.
START_REQUEST="$SHARED_DIR/worker-starting"
# How long to wait for fetch to record the run before giving up; fetch
# polls for the request every second.
RECORD_DEADLINE_SECONDS=60
# Endpoints a correctly isolated pod cannot connect to: the IPv4 cloud
# metadata server, the Kubernetes API server, and the EC2 IPv6 metadata
# server. The API server is always listening and always in the cluster's
# Service range.
: "${KUBERNETES_SERVICE_HOST:?}"
: "${KUBERNETES_SERVICE_PORT:?}"
ISOLATION_PROBES="169.254.169.254:80 $KUBERNETES_SERVICE_HOST:$KUBERNETES_SERVICE_PORT fd00:ec2::254:80"
ISOLATION_DEADLINE_SECONDS=60

main() {
	trap terminate TERM INT
	if [ -n "${ANTHROPIC_ENVIRONMENT_KEY:-}" ]; then
		echo "refusing to start: ANTHROPIC_ENVIRONMENT_KEY is set, and the environment key must never reach a sandbox" >&2
		exit 1
	fi
	if ! command -v nc >/dev/null 2>&1; then
		echo "refusing to start: nc is missing, so the network policy cannot be checked" >&2
		exit 1
	fi
	start_dir="$PWD"
	wait_until_isolated
	wait_for_secret
	record_run
	run_worker
}

# Passes only when one sweep finds the cluster's resolver reachable both
# before and after finding every probe unreachable. A pod whose network is
# not up yet fails every connection alike. That proves nothing.
wait_until_isolated() {
	waited=0
	while :; do
		if ! resolver_reachable; then
			blocked_by="this pod cannot reach its DNS resolver"
		elif reachable="$(any_probe_reachable)"; then
			blocked_by="this pod can still reach $reachable"
		elif resolver_reachable; then
			return 0
		else
			blocked_by="this pod lost its DNS resolver during the check"
		fi
		if [ "$waited" -ge "$ISOLATION_DEADLINE_SECONDS" ]; then
			echo "refusing to start: $blocked_by after ${ISOLATION_DEADLINE_SECONDS}s" >&2
			exit 1
		fi
		sleep 5
		waited=$((waited + 5))
	done
}

resolver_reachable() {
	resolver="$(awk '/^nameserver/ { print $2; exit }' /etc/resolv.conf)"
	[ -n "$resolver" ] && nc -z -w 1 "$resolver" 53 2>/dev/null
}

# Prints the first probe endpoint that accepts a TCP connection and returns
# 0; returns 1 when none does. Every probe runs whatever its address family:
# a pod without a route to an address cannot connect to it, which passes.
any_probe_reachable() {
	for probe in $ISOLATION_PROBES; do
		port="${probe##*:}"
		host="${probe%:*}"
		if nc -z -w 1 "$host" "$port" 2>/dev/null; then
			echo "$probe"
			return 0
		fi
	done
	return 1
}

# The secret is the last file fetch writes, so once it exists the others are
# complete. A restarted worker waits here until the dispatcher suspends the
# sandbox. The session's own code runs as this container's user. That user
# can rename entries of the shared directory. So only a directory this user
# neither owns nor can write is the one fetch wrote. Every check from here
# through record_run uses paths relative to the directory entered here. A
# directory the session's code swaps in afterwards is never read.
wait_for_secret() {
	until cd "$WORK_ITEM_DIR" 2>/dev/null && [ -r secret ] && [ ! -e worker-ran ]; do
		sleep 2
	done
	if [ -O . ] || [ -w . ]; then
		echo "refusing to run: this container owns or can write to $WORK_ITEM_DIR, so the work item cannot be trusted" >&2
		exit 1
	fi
	ANTHROPIC_SESSION_ID="$(value session_id)"
	ANTHROPIC_WORK_ID="$(value work_id)"
	ANTHROPIC_ENVIRONMENT_ID="$(value environment_id)"
	ANTHROPIC_BASE_URL="$(value base_url)"
	export ANTHROPIC_SESSION_ID ANTHROPIC_WORK_ID ANTHROPIC_ENVIRONMENT_ID ANTHROPIC_BASE_URL
	echo "running session $ANTHROPIC_SESSION_ID (work item $ANTHROPIC_WORK_ID)" >&2
}

value() {
	if [ ! -s "./$1" ]; then
		echo "the work item is missing $1" >&2
		exit 1
	fi
	cat "./$1"
}

# Asks fetch to record the run and waits for the record, so it is in place
# before any of the session's code runs. A record still missing at the
# deadline means fetch is not running.
record_run() {
	: > "$START_REQUEST"
	waited=0
	while [ ! -e worker-ran ] && [ "$waited" -lt "$RECORD_DEADLINE_SECONDS" ]; do
		sleep 1
		waited=$((waited + 1))
	done
	rm -f "$START_REQUEST"
	# Checked once more after the request is removed, so a record fetch
	# wrote in the last second of the wait still counts.
	if [ ! -e worker-ran ]; then
		echo "refusing to run: fetch did not record the run within ${RECORD_DEADLINE_SECONDS}s" >&2
		exit 1
	fi
}

# Runs the worker as a child and returns the worker's status as the
# script's. A SIGTERM to this script can then be forwarded to the worker.
# The worker's output goes straight to the pod's own.
run_worker() {
	# wait_for_secret left us inside the work-item directory; run the worker
	# from where the script started.
	cd "$start_dir"
	ANTHROPIC_WORK_SECRET_FILE="$SECRET_FILE" ant beta:worker run --workdir /workspace --max-idle "${CMA_MAX_IDLE:-30m}" &
	worker_pid=$!
	worker_status=0
	wait "$worker_pid" || worker_status=$?
	worker_pid=
	echo "worker for session $ANTHROPIC_SESSION_ID exited with status $worker_status" >&2
	return "$worker_status"
}

# Forwards the pod's stop signal to the worker and waits for it. The session
# is then left cleanly. Otherwise the worker would be killed at the end of
# the grace period.
terminate() {
	if [ -n "${worker_pid:-}" ]; then
		kill -TERM "$worker_pid" 2>/dev/null || true
		wait "$worker_pid" || true
	fi
	exit 143
}

main "$@"
