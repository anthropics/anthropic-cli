package sandbox

import (
	"context"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestWaitForRunningPod(t *testing.T) {
	objects := ownedGraph(map[string]string{LabelManagedBy: ManagedByValue})
	pod := objects[2].(*corev1.Pod)
	c, _, kube := newClient(t, "v1beta1", objects...)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	if _, _, err := c.WaitForRunningPod(ctx, "claim-a", 10*time.Millisecond); err == nil {
		t.Fatal("a Pending pod must not satisfy the wait")
	}

	pod.Status.Phase = corev1.PodRunning
	if _, err := kube.CoreV1().Pods(ns).UpdateStatus(context.Background(), pod, metav1.UpdateOptions{}); err != nil {
		t.Fatal(err)
	}
	sb, got, err := c.WaitForRunningPod(context.Background(), "claim-a", 10*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if got.UID != "pod-a-uid" || sb.Name != "sb-a" {
		t.Errorf("pod UID = %s, sandbox = %+v", got.UID, sb)
	}
}

func TestWaitForRunningPodRefusesTerminatingOrUnreadyPods(t *testing.T) {
	objects := ownedGraph(map[string]string{LabelManagedBy: ManagedByValue})
	pod := objects[2].(*corev1.Pod)
	pod.Status.Phase = corev1.PodRunning
	now := metav1.Now()
	pod.DeletionTimestamp = &now
	pod.Finalizers = []string{"test/hold"}
	c, _, _ := newClient(t, "v1beta1", objects...)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if _, _, err := c.WaitForRunningPod(ctx, "claim-a", 10*time.Millisecond); err == nil {
		t.Error("a terminating pod must not satisfy the wait")
	}

	objects = ownedGraph(map[string]string{LabelManagedBy: ManagedByValue})
	delete(objects[1].(*unstructured.Unstructured).Object, "status")
	objects[2].(*corev1.Pod).Status.Phase = corev1.PodRunning
	c, _, _ = newClient(t, "v1beta1", objects...)
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if _, _, err := c.WaitForRunningPod(ctx, "claim-a", 10*time.Millisecond); err == nil {
		t.Error("a v1beta1 Sandbox without Ready=True must not satisfy the wait")
	}
}

// A session's token is keyed to the pod PodOf resolves, so a pod that only
// wears the Sandbox's pod name — created by someone else while the real pod
// was absent — must read as no pod at all.
func TestPodOfRefusesAPodTheSandboxDoesNotOwn(t *testing.T) {
	for tc, mutate := range map[string]func(*corev1.Pod){
		"no owner":             func(p *corev1.Pod) { p.OwnerReferences = nil },
		"another owner":        func(p *corev1.Pod) { p.OwnerReferences[0].UID = "someone-else" },
		"owner not controller": func(p *corev1.Pod) { p.OwnerReferences[0].Controller = nil },
	} {
		objects := ownedGraph(map[string]string{LabelManagedBy: ManagedByValue})
		mutate(objects[2].(*corev1.Pod))
		c, _, _ := newClient(t, "v1beta1", objects...)
		sb, err := c.GetSandbox(context.Background(), "sb-a")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := c.PodOf(context.Background(), sb); !apierrors.IsNotFound(err) {
			t.Errorf("%s: PodOf = %v, want the look-alike pod answered as not found", tc, err)
		}
	}
}

func TestPodOfHonorsPodNameAnnotation(t *testing.T) {
	objects := ownedGraph(map[string]string{LabelManagedBy: ManagedByValue})
	sb := objects[1].(*unstructured.Unstructured)
	sb.SetAnnotations(map[string]string{podNameAnnotation: "warm-pod-7"})
	pod := objects[2].(*corev1.Pod)
	pod.Name = "warm-pod-7"
	c, _, _ := newClient(t, "v1beta1", objects...)
	got, err := c.GetSandbox(context.Background(), "sb-a")
	if err != nil || got.PodName != "warm-pod-7" {
		t.Fatalf("sandbox = %+v, %v", got, err)
	}
	if p, err := c.PodOf(context.Background(), got); err != nil || p.UID != "pod-a-uid" {
		t.Errorf("PodOf = %v, %v", p, err)
	}
}
