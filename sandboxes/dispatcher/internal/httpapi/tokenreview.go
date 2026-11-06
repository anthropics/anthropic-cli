package httpapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	authnv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

// TokenAuth authenticates a caller by the bearer token it presents, verified
// through a TokenReview.
type TokenAuth struct {
	Reviewer TokenReviewer
}

// Authenticate reads the bearer token and resolves the pod it was issued to.
func (a TokenAuth) Authenticate(r *http.Request) (Identity, error) {
	token, found := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !found || token == "" {
		return Identity{}, fmt.Errorf("%w: a bearer token is required", ErrUnauthenticated)
	}
	return a.Reviewer.Review(r.Context(), token)
}

// PodUIDExtraKey is the TokenReview user-extra key that names the pod a
// projected token was issued to.
const PodUIDExtraKey = "authentication.kubernetes.io/pod-uid"

var (
	// ErrUnauthenticated means the token did not verify, or not for the
	// dispatcher's audience.
	ErrUnauthenticated = errors.New("token is not valid for the dispatcher")
	// ErrNotSandbox means the token verified but belongs to another identity.
	ErrNotSandbox = errors.New("token does not belong to the sandbox service account")
	// ErrNoPod means the token is not bound to a pod.
	ErrNoPod = errors.New("token is not bound to a pod")
)

// Identity is the pod a verified token was issued to.
type Identity struct {
	PodUID types.UID
}

// TokenReviewer verifies a bearer token and names the pod holding it.
type TokenReviewer interface {
	Review(ctx context.Context, token string) (Identity, error)
}

// KubeTokenReviewer verifies tokens with the TokenReview API.
type KubeTokenReviewer struct {
	kube     kubernetes.Interface
	audience string
	username string
}

// NewKubeTokenReviewer accepts tokens minted for audience to the
// serviceAccount of namespace. The "default" account is refused because
// every pod without an explicit account holds one of its tokens.
func NewKubeTokenReviewer(kube kubernetes.Interface, namespace, serviceAccount, audience string) (*KubeTokenReviewer, error) {
	if kube == nil {
		panic("httpapi.NewKubeTokenReviewer: nil Kubernetes client")
	}
	if serviceAccount == "" || serviceAccount == "default" {
		return nil, fmt.Errorf("sandbox service account must be a dedicated account, not %q", serviceAccount)
	}
	if audience == "" {
		return nil, errors.New("token audience is required")
	}
	return &KubeTokenReviewer{
		kube:     kube,
		audience: audience,
		username: "system:serviceaccount:" + namespace + ":" + serviceAccount,
	}, nil
}

func (r *KubeTokenReviewer) Review(ctx context.Context, token string) (Identity, error) {
	review := &authnv1.TokenReview{Spec: authnv1.TokenReviewSpec{Token: token, Audiences: []string{r.audience}}}
	result, err := r.kube.AuthenticationV1().TokenReviews().Create(ctx, review, metav1.CreateOptions{})
	if err != nil {
		return Identity{}, fmt.Errorf("token review: %w", err)
	}
	if !result.Status.Authenticated || !slices.Contains(result.Status.Audiences, r.audience) {
		return Identity{}, ErrUnauthenticated
	}
	if result.Status.User.Username != r.username {
		return Identity{}, ErrNotSandbox
	}
	uids := result.Status.User.Extra[PodUIDExtraKey]
	if len(uids) != 1 || uids[0] == "" {
		return Identity{}, ErrNoPod
	}
	return Identity{PodUID: types.UID(uids[0])}, nil
}
