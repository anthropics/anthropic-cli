package sandbox

import (
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestWorkerStatus(t *testing.T) {
	handoff := time.Date(2026, 9, 3, 12, 0, 0, 0, time.UTC)
	ended := func(at time.Time) corev1.ContainerState {
		return corev1.ContainerState{Terminated: &corev1.ContainerStateTerminated{FinishedAt: metav1.NewTime(at)}}
	}
	running := corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}
	for name, tc := range map[string]struct {
		statuses   []corev1.ContainerStatus
		ok, exited bool
		restarts   int32
	}{
		"no worker yet":          {nil, false, false, 0},
		"running first run":      {[]corev1.ContainerStatus{{Name: WorkerContainer, State: running}}, true, false, 0},
		"exited, not restarted":  {[]corev1.ContainerStatus{{Name: WorkerContainer, State: ended(handoff.Add(time.Minute))}}, true, true, 0},
		"exited and restarted":   {[]corev1.ContainerStatus{{Name: WorkerContainer, State: running, LastTerminationState: ended(handoff.Add(time.Minute)), RestartCount: 1}}, true, true, 1},
		"crashed before handoff": {[]corev1.ContainerStatus{{Name: "fetch"}, {Name: WorkerContainer, State: running, LastTerminationState: ended(handoff.Add(-time.Minute)), RestartCount: 1}}, true, false, 1},
	} {
		pod := &corev1.Pod{Status: corev1.PodStatus{ContainerStatuses: tc.statuses}}
		status, ok := WorkerStatusOf(pod)
		if ok != tc.ok || status.RestartCount != tc.restarts || status.ExitedSince(handoff) != tc.exited {
			t.Errorf("%s: status = %+v ok = %v exited = %v", name, status, ok, status.ExitedSince(handoff))
		}
	}
}
