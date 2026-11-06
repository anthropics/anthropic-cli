package preflight

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// environmentKeyStaysOut fails when the container could read the environment
// key through env, valueFrom, envFrom or a mounted Secret or ConfigMap. An
// unreadable source counts as carrying it, since nothing shows it clean.
func (c *Checker) environmentKeyStaysOut(ctx context.Context, container corev1.Container, volumes map[string]corev1.Volume) error {
	const key = "ANTHROPIC_ENVIRONMENT_KEY"
	carries := func(how string) error {
		return fmt.Errorf("container %q %s; the environment key must never reach a sandbox pod", container.Name, how)
	}
	for _, e := range container.Env {
		if e.Name == key {
			return carries("carries " + key)
		}
		if v := e.ValueFrom; v != nil && ((v.SecretKeyRef != nil && v.SecretKeyRef.Key == key) || (v.ConfigMapKeyRef != nil && v.ConfigMapKeyRef.Key == key)) {
			return carries("reads the " + key + " key into " + e.Name)
		}
	}
	for _, e := range container.EnvFrom {
		switch {
		case e.SecretRef != nil:
			secret, err := c.Kube.CoreV1().Secrets(c.Options.Namespace).Get(ctx, e.SecretRef.Name, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("container %q imports Secret %q through envFrom, and the check could not read it to show it free of the environment key: %w", container.Name, e.SecretRef.Name, err)
			}
			if _, ok := secret.Data[key]; ok {
				return carries("imports Secret " + e.SecretRef.Name + ", which holds " + key + ",")
			}
		case e.ConfigMapRef != nil:
			cm, err := c.Kube.CoreV1().ConfigMaps(c.Options.Namespace).Get(ctx, e.ConfigMapRef.Name, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("container %q imports ConfigMap %q through envFrom, and the check could not read it to show it free of the environment key: %w", container.Name, e.ConfigMapRef.Name, err)
			}
			if _, ok := cm.Data[key]; ok {
				return carries("imports ConfigMap " + e.ConfigMapRef.Name + ", which holds " + key + ",")
			}
		}
	}
	secrets, configMaps := mountedSources(container, volumes)
	for _, name := range secrets {
		secret, err := c.Kube.CoreV1().Secrets(c.Options.Namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("container %q mounts Secret %q, and the check could not read it to show it free of the environment key: %w", container.Name, name, err)
		}
		if _, ok := secret.Data[key]; ok {
			return carries("mounts Secret " + name + ", which holds " + key + ",")
		}
	}
	for _, name := range configMaps {
		cm, err := c.Kube.CoreV1().ConfigMaps(c.Options.Namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("container %q mounts ConfigMap %q, and the check could not read it to show it free of the environment key: %w", container.Name, name, err)
		}
		if _, data := cm.Data[key]; data {
			return carries("mounts ConfigMap " + name + ", which holds " + key + ",")
		}
		if _, binary := cm.BinaryData[key]; binary {
			return carries("mounts ConfigMap " + name + ", which holds " + key + ",")
		}
	}
	return nil
}

// serviceAccountTokenSecretsStayOut fails when the container reads a Secret of
// type kubernetes.io/service-account-token, a long-lived API token that
// automountServiceAccountToken misses. An unreadable Secret counts as one.
func (c *Checker) serviceAccountTokenSecretsStayOut(ctx context.Context, container corev1.Container, volumes map[string]corev1.Volume) error {
	for _, name := range secretsRead(container, volumes) {
		secret, err := c.Kube.CoreV1().Secrets(c.Options.Namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("container %q reads Secret %q, and the check could not read it to show it is not a service-account token: %w", container.Name, name, err)
		}
		if secret.Type == corev1.SecretTypeServiceAccountToken {
			return fmt.Errorf("container %q reads Secret %q of type %s, a long-lived token for the Kubernetes API", container.Name, name, secret.Type)
		}
	}
	return nil
}

// secretsRead names the Secrets the container reads: through plain and
// projected volume mounts, envFrom, and env valueFrom sources.
func secretsRead(container corev1.Container, volumes map[string]corev1.Volume) []string {
	names, _ := mountedSources(container, volumes)
	for _, e := range container.EnvFrom {
		if e.SecretRef != nil {
			names = append(names, e.SecretRef.Name)
		}
	}
	for _, e := range container.Env {
		if v := e.ValueFrom; v != nil && v.SecretKeyRef != nil {
			names = append(names, v.SecretKeyRef.Name)
		}
	}
	return names
}

// mountedSources names the Secrets and ConfigMaps the container mounts,
// through plain and projected volumes.
func mountedSources(container corev1.Container, volumes map[string]corev1.Volume) (secrets, configMaps []string) {
	for _, m := range container.VolumeMounts {
		v := volumes[m.Name]
		switch {
		case v.Secret != nil:
			secrets = append(secrets, v.Secret.SecretName)
		case v.ConfigMap != nil:
			configMaps = append(configMaps, v.ConfigMap.Name)
		case v.Projected != nil:
			for _, src := range v.Projected.Sources {
				switch {
				case src.Secret != nil:
					secrets = append(secrets, src.Secret.Name)
				case src.ConfigMap != nil:
					configMaps = append(configMaps, src.ConfigMap.Name)
				}
			}
		}
	}
	return secrets, configMaps
}
