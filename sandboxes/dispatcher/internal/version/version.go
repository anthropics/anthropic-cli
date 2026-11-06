// Package version reports which sandboxes release a binary was built from.
package version

// release is the sandboxes release, set at link time by the image's build.
// `go build` from a checkout stamps no release into the binary.
var release string

// String is the release, or "development" for a build that is not one.
func String() string {
	if release == "" {
		return "development"
	}
	return release
}
