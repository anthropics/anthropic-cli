#!/bin/sh
# The fetch container: collect the pod's work item from the dispatcher once
# and leave it for the worker, one value per file, in a directory of the
# shared volume that the worker can read but not change. POSIX sh so it runs
# under busybox.
set -eu

: "${CMA_DISPATCHER_URL:?}"
: "${CMA_CA_FILE:?}"
# The pod proves itself to the dispatcher one of two ways. With
# CMA_ASSIGNMENT_NONCE set (the GKE add-on, which bans every kubelet
# credential in a sandbox pod) the request carries the nonce in a header.
# Otherwise it carries the projected service-account token from
# CMA_SA_TOKEN_FILE, re-read on each call because the kubelet rotates it.
if [ -z "${CMA_ASSIGNMENT_NONCE:-}" ]; then
	: "${CMA_SA_TOKEN_FILE:?}"
fi
SHARED_DIR="${CMA_SHARED_DIR:-/run/cma}"
WORK_ITEM_DIR="$SHARED_DIR/work-item"
# HOME and every file this script reads back live on a volume only this
# container mounts. The session's code can rename entries of the shared
# volume, so a scratch file there could be swapped for a forged work item
# between the write and the read.
FETCH_DIR="${CMA_FETCH_DIR:-/run/cma-fetch}"
HOME="$FETCH_DIR/home"
export HOME
# The worker asks fetch to record that the pod's session ran by creating
# START_REQUEST. fetch answers by writing worker-ran in the work-item
# directory, which the session's code cannot write. COLLECTED, on this
# container's own volume, marks that the work item was written, so a
# restarted fetch reuses it instead of asking the dispatcher again.
START_REQUEST="$SHARED_DIR/worker-starting"
COLLECTED="$FETCH_DIR/collected"
# Caps the poll interval, which is also the worst-case handoff delay for a
# pod waiting in the warm pool.
POLL_BACKOFF_MAX_SECONDS=10
# Bounds one request to the dispatcher. A dispatcher that accepts the
# connection and never answers then costs one poll, which the backoff
# retries. A stop signal reaches the script only between commands, and a
# curl that never returns would hold it off for good. The dispatcher stops
# writing an answer about 10 s after it reads the request, so a longer
# limit gains at most the time it takes to connect. 10 s is also the
# poll cap, and a stop signal therefore waits at most that long, inside
# the 30 s grace period Kubernetes gives the pod by default.
POLL_REQUEST_MAX_SECONDS=10

main() {
	trap 'exit 143' TERM INT
	if [ -n "${ANTHROPIC_ENVIRONMENT_KEY:-}" ]; then
		echo "refusing to start: ANTHROPIC_ENVIRONMENT_KEY is set, and the environment key must never reach a sandbox" >&2
		exit 1
	fi
	# A symlink here would point the ownership check and the chmod at
	# whatever the session's code linked to.
	if [ -L "$HOME" ] || [ -L "$WORK_ITEM_DIR" ]; then
		echo "refusing to start: $HOME or $WORK_ITEM_DIR is a symlink" >&2
		exit 1
	fi
	mkdir -p "$HOME" && chmod 0700 "$HOME"
	if [ -e "$COLLECTED" ]; then
		resume_work_item
	else
		write_work_item
	fi
	while :; do pause 3600; done
}

# Collects the work item once the dispatcher answers, leaves it in the
# work-item directory, then records the worker's run once the worker asks.
write_work_item() {
	mkdir -p "$WORK_ITEM_DIR"
	if [ ! -O "$WORK_ITEM_DIR" ]; then
		echo "refusing to start: $WORK_ITEM_DIR belongs to another user" >&2
		exit 1
	fi
	chmod 0750 "$WORK_ITEM_DIR"
	if [ -e "$WORK_ITEM_DIR/secret" ]; then
		echo "an earlier process in this pod already wrote the work item" >&2
	else
		wait_for_work_item
	fi
	if [ ! -e "$WORK_ITEM_DIR/secret" ]; then
		return 0
	fi
	: > "$COLLECTED"
	cd "$WORK_ITEM_DIR"
	record_run
}

# Resumes after a restart of this container with the work item it already
# wrote. The session's code may be running and can rename the work-item
# directory away, but cannot create a directory this container owns, so an
# owned directory holding the secret is the one this container wrote.
# If it is gone it is not fetched again (the dispatcher would return the
# same work item). record_run uses relative paths inside the directory
# entered here so a swapped-in directory is never written.
resume_work_item() {
	if ! cd "$WORK_ITEM_DIR" 2>/dev/null || [ ! -O . ] || [ ! -e secret ]; then
		echo "this pod's work item is gone; not collecting it again" >&2
		return 0
	fi
	record_run
}

# Waits for the worker's request, then writes worker-ran in the current
# directory. The worker runs no session code before it sees worker-ran. A
# worker-ran already there was written before a restart of this container
# and is kept.
record_run() {
	if [ -e worker-ran ]; then
		return 0
	fi
	until [ -e "$START_REQUEST" ]; do pause 1; done
	: > worker-ran
}

# Polls GET /work-item until 200, then writes the work item. Each poll costs a
# TokenReview, so the wait doubles up to the cap. 404 is "nothing yet", 410 says
# the worker has exited, and 403 repeats until the dispatcher sees the pod run.
wait_for_work_item() {
	delay=2
	while :; do
		# What is read after the request is then this answer's, whatever curl
		# does with the file of an earlier one.
		rm -f "$HOME/response"
		status="$(dispatcher_request GET /work-item "$HOME/response")"
		case "$status" in
		200)
			mv "$HOME/response" "$HOME/work-item.json"
			break
			;;
		404) ;;
		410)
			echo "this pod's work item is closed" >&2
			return 0
			;;
		403) echo "dispatcher answered 403 to GET /work-item ($(refusal_reason "$HOME/response")); retrying in ${delay}s" >&2 ;;
		*) echo "dispatcher answered $status to GET /work-item; retrying in ${delay}s" >&2 ;;
		esac
		pause "$delay"
		delay="$(next_delay "$delay")"
	done
	# The secret comes last: its presence tells the worker the rest is there.
	for field in session_id work_id environment_id base_url secret; do
		value="$(work_item_field "$field")"
		write_value "$field" "$value"
	done
	rm -f "$HOME/work-item.json"
	echo "assigned session $(cat "$WORK_ITEM_DIR/session_id") (work item $(cat "$WORK_ITEM_DIR/work_id"))" >&2
}

# refusal_reason FILE prints why the dispatcher refused, FILE holding its answer.
# It prints the first line only if that line is listed here, so the log line
# that reports a 403 holds nothing the answer's writer chose.
refusal_reason() {
	[ -s "$1" ] || { printf 'no reason given'; return; }
	reason="$(head -c 256 "$1" | head -n 1)"
	case "$reason" in
	'token does not belong to the sandbox service account' | \
		'token is not bound to a pod' | \
		'token is not bound to a pod: the request has no resolvable source address' | \
		'token is not bound to a pod: the sandbox has no running pod' | \
		'token is not bound to a pod: request did not come from the sandbox pod')
		printf '%s' "$reason"
		;;
	*) printf 'a reason this script does not know' ;;
	esac
}

work_item_field() {
	value="$(jq -r --arg k "$1" '.[$k] // empty' "$HOME/work-item.json")"
	if [ -z "$value" ]; then
		echo "the dispatcher's work item is missing $1" >&2
		exit 1
	fi
	printf '%s' "$value"
}

# Writes one value, readable by the pod's group and nobody else, in place
# only once complete.
write_value() {
	(umask 037 && printf '%s' "$2" > "$WORK_ITEM_DIR/.$1.tmp")
	chmod 0440 "$WORK_ITEM_DIR/.$1.tmp"
	mv "$WORK_ITEM_DIR/.$1.tmp" "$WORK_ITEM_DIR/$1"
}

# The script is the container's first process, so a stop signal reaches it
# only between commands; sleeping in the background lets the trap run at once.
pause() {
	sleep "$1" &
	wait "$!"
}

next_delay() {
	if [ "$(($1 * 2))" -lt "$POLL_BACKOFF_MAX_SECONDS" ]; then
		echo "$(($1 * 2))"
	else
		echo "$POLL_BACKOFF_MAX_SECONDS"
	fi
}

# dispatcher_request METHOD PATH OUTFILE prints the HTTP status; curl itself
# prints 000 when no connection was made or no status line arrived in time.
# A stall after the status line leaves a short or missing body, so the move
# or the value reads fail and the script exits. The pod authenticates by the
# assignment nonce when it has one, otherwise by the projected token re-read
# on each call. -q, first, keeps curl from reading any .curlrc.
dispatcher_request() {
	if [ -n "${CMA_ASSIGNMENT_NONCE:-}" ]; then
		auth="X-CMA-Assignment-Nonce: $CMA_ASSIGNMENT_NONCE"
	else
		auth="Authorization: Bearer $(cat "$CMA_SA_TOKEN_FILE")"
	fi
	curl -q -sS -o "$3" -w '%{http_code}' --cacert "$CMA_CA_FILE" \
		--max-time "$POLL_REQUEST_MAX_SECONDS" \
		-H "$auth" -X "$1" "$CMA_DISPATCHER_URL$2" || true
}

main "$@"
