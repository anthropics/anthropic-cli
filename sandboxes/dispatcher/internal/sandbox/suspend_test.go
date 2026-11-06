package sandbox

import (
	"context"
	"strings"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	k8stesting "k8s.io/client-go/testing"
)

func TestSetSuspendedPerVersion(t *testing.T) {
	for version, want := range map[string]struct {
		field              string
		suspended, running any
	}{
		"v1beta1":  {"operatingMode", "Suspended", "Running"},
		"v1alpha1": {"replicas", int64(0), int64(1)},
	} {
		sb := object("agents.x-k8s.io/"+version, "Sandbox", "sb-1", "sb-uid", map[string]any{"spec": map[string]any{"podTemplate": map[string]any{}}})
		c, dyn, _ := newClient(t, version, sb)
		storedField := func() any {
			stored, _ := dyn.Resource(c.shape.Sandboxes()).Namespace(ns).Get(context.Background(), "sb-1", metav1.GetOptions{})
			v, _, _ := unstructured.NestedFieldNoCopy(stored.Object, "spec", want.field)
			return v
		}
		if err := c.Suspend(context.Background(), "sb-1"); err != nil {
			t.Fatal(err)
		}
		if got, _ := c.GetSandbox(context.Background(), "sb-1"); !got.Suspended || storedField() != want.suspended {
			t.Errorf("%s: after Suspend spec.%s = %v, want %v (read back as suspended: %v)", version, want.field, storedField(), want.suspended, got.Suspended)
		}
		if err := c.Resume(context.Background(), "sb-1", time.Millisecond); err != nil {
			t.Fatal(err)
		}
		if got, _ := c.GetSandbox(context.Background(), "sb-1"); got.Suspended || storedField() != want.running {
			t.Errorf("%s: after Resume spec.%s = %v, want %v (read back as suspended: %v)", version, want.field, storedField(), want.running, got.Suspended)
		}
	}
}

func TestSuspendAndResumeFailWhenTheClusterIgnoresTheField(t *testing.T) {
	sb := object("agents.x-k8s.io/v1beta1", "Sandbox", "sb-1", "sb-uid", map[string]any{"spec": map[string]any{"podTemplate": map[string]any{}}})
	c, dyn, _ := newClient(t, "v1beta1", sb)
	dyn.PrependReactor("patch", "sandboxes", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, sb.DeepCopy(), nil
	})
	resume := func(ctx context.Context, name string) error { return c.Resume(ctx, name, time.Millisecond) }
	for name, change := range map[string]func(context.Context, string) error{"Suspend": c.Suspend, "Resume": resume} {
		if err := change(context.Background(), "sb-1"); err == nil || !strings.Contains(err.Error(), "api-version") {
			t.Errorf("%s against a cluster that drops the field: err = %v", name, err)
		}
	}
}

func TestResumeWaitsForThePodFromBeforeTheSuspend(t *testing.T) {
	objects := ownedGraph(map[string]string{LabelManagedBy: ManagedByValue})
	pod := objects[2].(*corev1.Pod)
	pod.Status.Phase = corev1.PodRunning
	c, _, kube := newClient(t, "v1beta1", objects...)
	if err := c.Suspend(context.Background(), "sb-a"); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if err := c.Resume(ctx, "sb-a", 10*time.Millisecond); err == nil {
		t.Fatal("Resume must wait while the old pod is still running")
	}
	if sb, _ := c.GetSandbox(context.Background(), "sb-a"); !sb.Suspended {
		t.Fatal("a Resume that gave up must leave the sandbox suspended")
	}

	if err := kube.CoreV1().Pods(ns).Delete(context.Background(), "sb-a", metav1.DeleteOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := c.Resume(context.Background(), "sb-a", 10*time.Millisecond); err != nil {
		t.Fatal(err)
	}
	if sb, _ := c.GetSandbox(context.Background(), "sb-a"); sb.Suspended {
		t.Error("Resume after the old pod is gone must clear the suspend")
	}
}
