package core

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"
)

// URLFetcher materializes a URL-referenced resource into a local directory and
// reports the immutable revision it resolved to.
type URLFetcher interface {
	// Fetch downloads rawURL. When pin carries a revision the fetcher must
	// return exactly that rather than whatever the URL's branch currently
	// points at.
	Fetch(ctx context.Context, rawURL string, pin URLPin) (*FetchedURL, error)
}

// URLPin is what the lockfile remembers about a URL so a later apply lands
// on the same content without asking the host again.
type URLPin struct {
	// Revision is the immutable identifier the URL resolved to — a commit SHA
	// for GitHub.
	Revision string
	// Subpath is where in that revision the resource lives. A tree URL does
	// not say where the branch name ends and the path begins; this records
	// the answer so a pinned fetch need not rediscover it.
	Subpath string
}

// FetchedURL is a URL-referenced resource materialized on local disk.
type FetchedURL struct {
	// Dir is the local directory the URL was materialized into.
	Dir string
	URLPin
}

// maxArchiveBytes caps a downloaded tarball. URL-sourced resources are
// documentation-sized; anything larger is a misaimed URL, and decompressing it
// unbounded would be a gift to whoever controls that repo.
const maxArchiveBytes = 64 << 20 // 64 MiB

// GitHubFetcher fetches resources from github.com tree URLs. Other hosts get a
// clear error rather than a half-working guess. Build one with
// NewGitHubFetcher; the zero value has no cache to record fetches in.
type GitHubFetcher struct {
	HTTPClient *http.Client
	// Token authenticates against private repositories. Read from GITHUB_TOKEN
	// or GH_TOKEN by NewGitHubFetcher.
	Token string
	// CacheDir holds extracted archives across applies.
	CacheDir string
	// apiBase overrides the GitHub API root; tests point it at a local server.
	apiBase string

	// cache memoizes fetches within one run, keyed by repository, commit and
	// subpath.
	cache map[string]*FetchedURL
}

// NewGitHubFetcher returns a fetcher that keeps extracted archives under
// cacheDir and authenticates with GITHUB_TOKEN, or GH_TOKEN when that is unset.
func NewGitHubFetcher(cacheDir string) *GitHubFetcher {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		token = os.Getenv("GH_TOKEN")
	}
	return &GitHubFetcher{
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
		Token:      token,
		CacheDir:   cacheDir,
		cache:      map[string]*FetchedURL{},
	}
}

const githubAPIBase = "https://api.github.com"

// apiRoot is the GitHub API root this fetcher talks to.
func (g *GitHubFetcher) apiRoot() string {
	if g.apiBase != "" {
		return g.apiBase
	}
	return githubAPIBase
}

// repoEndpoint is the API URL for tail under the repository ref names.
func (g *GitHubFetcher) repoEndpoint(ref githubRef, tail string) string {
	return fmt.Sprintf("%s/repos/%s/%s/%s", g.apiRoot(), url.PathEscape(ref.Owner), url.PathEscape(ref.Repo), tail)
}

var fullSHA = regexp.MustCompile(`^[0-9a-f]{40}$`)

type githubRef struct {
	Owner string
	Repo  string
	// Rest is everything after /tree/: a git ref followed by a path within the
	// repository. Where one ends and the other begins is ambiguous; resolve
	// settles it.
	Rest []string
}

// parseGitHubURL understands the forms a person actually pastes: a repo root,
// a branch tree, or a subdirectory within a branch tree.
func parseGitHubURL(rawURL string) (githubRef, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return githubRef{}, fmt.Errorf("not a valid URL: %w", err)
	}
	host := strings.ToLower(u.Host)
	if host != "github.com" && host != "www.github.com" {
		return githubRef{}, fmt.Errorf("only github.com URLs are supported, got %q", u.Host)
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return githubRef{}, fmt.Errorf("expected a URL like https://github.com/owner/repo/tree/<ref>/<path>")
	}
	ref := githubRef{Owner: parts[0], Repo: strings.TrimSuffix(parts[1], ".git")}
	if len(parts) == 2 {
		return ref, nil
	}
	if parts[2] != "tree" && parts[2] != "blob" {
		return githubRef{}, fmt.Errorf("expected /tree/<ref>/<path> after the repository, got /%s/", parts[2])
	}
	if len(parts) < 4 {
		return githubRef{}, fmt.Errorf("missing a git ref after /%s/", parts[2])
	}
	ref.Rest = parts[3:]
	// A ".." would send the extracted subdirectory outside the cache and
	// let a URL in a config file bundle up any directory on the machine.
	if slices.Contains(ref.Rest, "..") {
		return githubRef{}, fmt.Errorf("path segments may not contain \"..\"")
	}
	return ref, nil
}

// split interprets the first n segments as the git ref and the rest as a path
// within the repository.
func (r githubRef) split(n int) (gitRef, subpath string) {
	return strings.Join(r.Rest[:n], "/"), path.Join(r.Rest[n:]...)
}

// Fetch implements URLFetcher. An unpinned URL is resolved against the GitHub
// API first; a pinned one goes straight to its recorded commit. Every error it
// returns starts with rawURL, so the user can tell which reference failed.
func (g *GitHubFetcher) Fetch(ctx context.Context, rawURL string, pin URLPin) (*FetchedURL, error) {
	fetched, err := g.fetch(ctx, rawURL, pin)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", rawURL, err)
	}
	return fetched, nil
}

// fetch does the work of Fetch; its errors leave out the URL, which Fetch adds.
func (g *GitHubFetcher) fetch(ctx context.Context, rawURL string, pin URLPin) (*FetchedURL, error) {
	ref, err := parseGitHubURL(rawURL)
	if err != nil {
		return nil, err
	}

	if pin.Revision == "" {
		pin, err = g.resolve(ctx, ref)
		if err != nil {
			return nil, err
		}
	}
	// The pin comes from a committed lockfile, so it is exactly as trustworthy
	// as a URL somebody else wrote: a revision lands in a cache path and an API
	// URL, and a subpath picks which directory gets bundled and uploaded.
	sha := pin.Revision
	if !fullSHA.MatchString(sha) {
		return nil, fmt.Errorf("%q is not a full commit SHA", sha)
	}
	if pin.Subpath != "" && !filepath.IsLocal(filepath.FromSlash(pin.Subpath)) {
		return nil, fmt.Errorf("subpath %q escapes the repository", pin.Subpath)
	}

	cacheKey := ref.Owner + "/" + ref.Repo + "@" + sha + ":" + pin.Subpath
	if hit, ok := g.cache[cacheKey]; ok {
		return hit, nil
	}

	root, err := g.cacheEntry(ref, sha)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(root); err != nil {
		if err := g.downloadTarball(ctx, ref, sha, root); err != nil {
			return nil, err
		}
	}
	dir := filepath.Join(root, filepath.FromSlash(pin.Subpath))
	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		return nil, fmt.Errorf("no directory %q at commit %s", pin.Subpath, sha[:12])
	}

	out := &FetchedURL{Dir: dir, URLPin: pin}
	g.cache[cacheKey] = out
	return out, nil
}

// resolve turns a tree URL into a commit SHA and the in-repo path under it.
//
// Where the git ref ends and the path begins is genuinely ambiguous:
// "tree/feature/new-skill/pdf" is branch "feature" with path
// "new-skill/pdf", or branch "feature/new-skill" with path "pdf". GitHub's
// own UI settles it against the repository's refs, and so does this: ask the
// API shortest candidate first — one call for the overwhelmingly common
// "tree/main/...", a few more only for a branch with slashes in it. The
// answer is recorded in the pin so it is only ever asked once.
func (g *GitHubFetcher) resolve(ctx context.Context, ref githubRef) (URLPin, error) {
	if len(ref.Rest) == 0 {
		sha, err := g.resolveCommit(ctx, ref, "HEAD")
		return URLPin{Revision: sha}, err
	}

	var firstErr error
	for n := 1; n <= len(ref.Rest); n++ {
		candidate, subpath := ref.split(n)
		if fullSHA.MatchString(candidate) {
			return URLPin{Revision: candidate, Subpath: subpath}, nil
		}
		sha, err := g.resolveCommit(ctx, ref, candidate)
		if err == nil {
			return URLPin{Revision: sha, Subpath: subpath}, nil
		}
		if firstErr == nil {
			firstErr = err
		}
	}
	// Report the first failure: for the common case of a plain branch name
	// that is the accurate diagnosis, and the later attempts are noise.
	return URLPin{}, firstErr
}

func (g *GitHubFetcher) request(ctx context.Context, endpoint, accept string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", accept)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if g.Token != "" {
		req.Header.Set("Authorization", "Bearer "+g.Token)
	}
	resp, err := g.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		resp.Body.Close()
		hint := ""
		if g.Token == "" {
			switch resp.StatusCode {
			case http.StatusNotFound:
				hint = " (set GITHUB_TOKEN if the repository is private)"
			case http.StatusForbidden:
				hint = " (unauthenticated GitHub API requests are rate-limited; set GITHUB_TOKEN)"
			}
		}
		return nil, fmt.Errorf("GET %s: %s%s: %s", endpoint, resp.Status, hint, strings.TrimSpace(string(body)))
	}
	return resp, nil
}

func (g *GitHubFetcher) resolveCommit(ctx context.Context, ref githubRef, gitRef string) (string, error) {
	resp, err := g.request(ctx, g.repoEndpoint(ref, "commits/"+url.PathEscape(gitRef)), "application/vnd.github+json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var commit struct {
		SHA string `json:"sha"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&commit); err != nil {
		return "", fmt.Errorf("decoding commit for %s/%s@%s: %w", ref.Owner, ref.Repo, gitRef, err)
	}
	if commit.SHA == "" {
		return "", fmt.Errorf("no commit found for %s/%s@%s", ref.Owner, ref.Repo, gitRef)
	}
	return commit.SHA, nil
}

// cacheEntry is the directory the repository at sha is unpacked into. It sits
// under one directory per repository and one per full SHA: two commits can
// share a shorter prefix, and joining owner and repo with a hyphen cannot tell
// a-b/c from a/b-c. Its own name keeps the hyphenated 12-character form: a
// skill at a repository root with no name of its own is displayed under this
// name, and a skill's display name cannot change once it is created.
func (g *GitHubFetcher) cacheEntry(ref githubRef, sha string) (string, error) {
	name := ref.Owner + "-" + ref.Repo + "-" + sha[:12]
	// Owner and repo come from a URL; on Windows a backslash in either could
	// place the entry outside the cache.
	if strings.ContainsAny(name, `/\`) {
		return "", fmt.Errorf("%s/%s is not a valid repository name: it contains a path separator", ref.Owner, ref.Repo)
	}
	// Escaping turns the "/" into "%2F", so the repository is one path element
	// that no other owner and repo can spell.
	repo := url.QueryEscape(ref.Owner + "/" + ref.Repo)
	return filepath.Join(g.CacheDir, repo, sha, name), nil
}

func (g *GitHubFetcher) downloadTarball(ctx context.Context, ref githubRef, sha, dest string) error {
	resp, err := g.request(ctx, g.repoEndpoint(ref, "tarball/"+sha), "application/vnd.github+json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tmp, err := os.MkdirTemp(g.CacheDir, "extract-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	if err := extractTarGz(io.LimitReader(resp.Body, maxArchiveBytes), tmp); err != nil {
		return fmt.Errorf("extracting %s/%s@%s: %w", ref.Owner, ref.Repo, sha[:12], err)
	}

	// GitHub tarballs wrap everything in a single owner-repo-sha directory.
	entries, err := os.ReadDir(tmp)
	if err != nil {
		return err
	}
	root := tmp
	if len(entries) == 1 && entries[0].IsDir() {
		root = filepath.Join(tmp, entries[0].Name())
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	if err := os.Rename(root, dest); err != nil {
		// A concurrent apply may have won the race; that's fine.
		if _, statErr := os.Stat(dest); statErr == nil {
			return nil
		}
		return err
	}
	return nil
}

// extractTarGz unpacks a gzipped tar into dest, refusing any entry that would
// escape it. Archives are attacker-controlled input whenever the URL is.
func extractTarGz(r io.Reader, dest string) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gz.Close()

	destAbs, err := filepath.Abs(dest)
	if err != nil {
		return err
	}
	tr := tar.NewReader(gz)
	var written int64
	for {
		header, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		target := filepath.Join(destAbs, filepath.FromSlash(header.Name))
		// filepath.Join cleans the path, so a "../" entry lands outside destAbs
		// and this check catches it.
		if target != destAbs && !strings.HasPrefix(target, destAbs+string(os.PathSeparator)) {
			return fmt.Errorf("archive entry %q escapes the extraction directory", header.Name)
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			n, err := extractFile(target, io.LimitReader(tr, maxArchiveBytes-written))
			if err != nil {
				return err
			}
			written += n
			if written >= maxArchiveBytes {
				return fmt.Errorf("archive exceeds the %d MiB limit", maxArchiveBytes>>20)
			}
		default:
			// Symlinks and devices have no place in a resource bundle, and
			// following them is how extraction turns into arbitrary write.
			continue
		}
	}
}

// extractFile writes one regular archive entry to target and reports how
// many bytes it consumed from r.
func extractFile(target string, r io.Reader) (int64, error) {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return 0, err
	}
	f, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return 0, err
	}
	n, err := io.Copy(f, r)
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	return n, err
}
