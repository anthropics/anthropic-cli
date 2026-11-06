package preflight

import (
	"context"
	"fmt"

	authzv1 "k8s.io/api/authorization/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/sandbox"
)

// RoleVerbs is the namespaced access the dispatcher needs, by resource.
type RoleVerbs struct {
	Group, Resource string
	Verbs           []string
}

// RequiredRole is the namespaced access the dispatcher needs on the shape's
// path. The add-on path creates and deletes Sandboxes directly, writes each
// sandbox's nonce Secret (create alone: it never reads one back, and the
// Secret is garbage-collected with its Sandbox), and needs no claims, warm
// pools, or TokenReview; the upstream path drives claims and verifies pod
// tokens.
func RequiredRole(shape sandbox.Shape) []RoleVerbs {
	if shape.Direct() {
		return []RoleVerbs{
			{"agents.x-k8s.io", "sandboxes", []string{"create", "get", "list", "watch", "patch", "delete"}},
			{"extensions.agents.x-k8s.io", "sandboxtemplates", []string{"get"}},
			{"", "secrets", []string{"create"}},
			{"", "pods", []string{"get", "list", "watch"}},
		}
	}
	return []RoleVerbs{
		{"extensions.agents.x-k8s.io", "sandboxclaims", []string{"create", "get", "list", "watch", "delete", "patch"}},
		{"extensions.agents.x-k8s.io", "sandboxwarmpools", []string{"get"}},
		{"agents.x-k8s.io", "sandboxes", []string{"get", "list", "watch", "patch"}},
		{"", "pods", []string{"get", "list", "watch"}},
	}
}

func (c *Checker) rbac(ctx context.Context) error {
	user := "system:serviceaccount:" + c.Options.Namespace + ":" + c.Options.DispatcherServiceAccount
	for _, rv := range RequiredRole(c.Options.Shape) {
		for _, verb := range rv.Verbs {
			if err := c.allowed(ctx, user, authzv1.ResourceAttributes{Namespace: c.Options.Namespace, Group: rv.Group, Resource: rv.Resource, Verb: verb}); err != nil {
				return err
			}
		}
	}
	if c.Options.Shape.Direct() {
		return nil
	}
	return c.allowed(ctx, user, authzv1.ResourceAttributes{Group: "authentication.k8s.io", Resource: "tokenreviews", Verb: "create"})
}

func (c *Checker) allowed(ctx context.Context, user string, attrs authzv1.ResourceAttributes) error {
	groups := []string{"system:serviceaccounts", "system:serviceaccounts:" + c.Options.Namespace, "system:authenticated"}
	review := &authzv1.SubjectAccessReview{Spec: authzv1.SubjectAccessReviewSpec{User: user, Groups: groups, ResourceAttributes: &attrs}}
	result, err := c.Kube.AuthorizationV1().SubjectAccessReviews().Create(ctx, review, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsForbidden(err) {
			return fmt.Errorf("your credentials may not create SubjectAccessReviews, so the dispatcher's access cannot be checked: %w", err)
		}
		return err
	}
	if !result.Status.Allowed {
		return fmt.Errorf("%s may not %s %s%s", user, attrs.Verb, attrs.Resource, groupSuffix(attrs.Group))
	}
	return nil
}

func groupSuffix(group string) string {
	if group == "" {
		return ""
	}
	return "." + group
}
