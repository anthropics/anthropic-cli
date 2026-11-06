// Package httpapi serves the one route a sandbox pod calls, collecting its
// work item. Callers are authenticated either through a Kubernetes
// TokenReview of the pod's projected service-account token (the upstream
// path) or through an assignment nonce bound to the pod's network position
// (the add-on path).
package httpapi
