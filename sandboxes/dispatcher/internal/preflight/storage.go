package preflight

import (
	"context"
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func (c *Checker) storageClass(ctx context.Context) error {
	tpl, err := c.template(ctx)
	if err != nil {
		return err
	}
	templates, _, _ := unstructured.NestedSlice(tpl.Object, "spec", "volumeClaimTemplates")
	if len(templates) == 0 {
		return errors.New("the SandboxTemplate has no volumeClaimTemplates; /workspace would not survive a suspend")
	}
	var warnings []error
	for _, raw := range templates {
		vct, _ := raw.(map[string]any)
		className, _, _ := unstructured.NestedString(vct, "spec", "storageClassName")
		err := c.claimStorageClass(ctx, className)
		if err != nil && !errors.Is(err, ErrWarning) {
			return err
		}
		warnings = append(warnings, err)
	}
	return errors.Join(warnings...)
}

// claimStorageClass fails unless the StorageClass className names, or the
// cluster default when it is empty, exists. It warns on a Retain class, since
// deleted sessions then leave volumes behind.
func (c *Checker) claimStorageClass(ctx context.Context, className string) error {
	sc, err := c.findStorageClass(ctx, className)
	if err != nil {
		return err
	}
	if sc.ReclaimPolicy != nil && *sc.ReclaimPolicy == corev1.PersistentVolumeReclaimRetain {
		return fmt.Errorf("%w: StorageClass %s reclaims with Retain, so deleted sessions leave volumes behind", ErrWarning, sc.Name)
	}
	return nil
}

func (c *Checker) findStorageClass(ctx context.Context, name string) (*storagev1.StorageClass, error) {
	if name != "" {
		sc, err := c.Kube.StorageV1().StorageClasses().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("StorageClass %s: %w", name, err)
		}
		return sc, nil
	}
	list, err := c.Kube.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for i := range list.Items {
		// Kubernetes treats either annotation as marking the default class.
		annotations := list.Items[i].Annotations
		if annotations["storageclass.kubernetes.io/is-default-class"] == "true" || annotations["storageclass.beta.kubernetes.io/is-default-class"] == "true" {
			return &list.Items[i], nil
		}
	}
	return nil, errors.New("the volume claim template names no storageClassName and the cluster has no default StorageClass")
}
