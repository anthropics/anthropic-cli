package sandbox

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	kubefake "k8s.io/client-go/kubernetes/fake"
)

func TestNewRequiresAnEnvironmentID(t *testing.T) {
	_, err := New(dynamicfake.NewSimpleDynamicClient(runtime.NewScheme()), kubefake.NewSimpleClientset(), Options{})
	if err == nil {
		t.Error("a client without an environment id must be refused")
	}
}
