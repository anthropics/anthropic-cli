package sandbox

import (
	"time"

	corev1 "k8s.io/api/core/v1"
)

// WorkerContainer is the container of a sandbox pod that runs the session.
const WorkerContainer = "worker"

// FetchContainer collects the pod's work item from the dispatcher; on the
// add-on path it alone carries the assignment nonce, on the upstream path the
// projected token.
const FetchContainer = "fetch"

// EnvAssignmentNonce names the fetch container's per-sandbox secret on the
// add-on path. The pod presents it to the dispatcher in place of a
// kubelet-issued credential, which the add-on's admission policy forbids.
const EnvAssignmentNonce = "CMA_ASSIGNMENT_NONCE"

// WorkspaceVolume is the pod volume the worker container mounts at /workspace.
const WorkspaceVolume = "workspace"

// NonceSecretKey is the key inside a sandbox's nonce Secret whose value
// EnvAssignmentNonce reads.
const NonceSecretKey = "nonce"

// WorkerStatus is what the kubelet reports about the worker container.
type WorkerStatus struct {
	RestartCount int32
	// Terminated is when the container's current run ended; zero while it
	// runs or waits.
	Terminated time.Time
	// LastTerminated is when the run before the current one ended; zero on
	// the first run.
	LastTerminated time.Time
}

// WorkerStatusOf reads the worker container's status from pod. ok is false
// while the pod reports no such container.
func WorkerStatusOf(pod *corev1.Pod) (status WorkerStatus, ok bool) {
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.Name != WorkerContainer {
			continue
		}
		status.RestartCount = cs.RestartCount
		if t := cs.State.Terminated; t != nil {
			status.Terminated = t.FinishedAt.Time
		}
		if t := cs.LastTerminationState.Terminated; t != nil {
			status.LastTerminated = t.FinishedAt.Time
		}
		return status, true
	}
	return WorkerStatus{}, false
}

// ExitedPast reports whether the worker restarted past baseline or its
// current run has ended.
func (s WorkerStatus) ExitedPast(baseline int32) bool {
	return s.RestartCount > baseline || !s.Terminated.IsZero()
}

// ExitedSince reports whether the worker ended a run at or after since.
func (s WorkerStatus) ExitedSince(since time.Time) bool {
	return endedSince(s.Terminated, since) || endedSince(s.LastTerminated, since)
}

func endedSince(ended, since time.Time) bool {
	return !ended.IsZero() && !ended.Before(since)
}
