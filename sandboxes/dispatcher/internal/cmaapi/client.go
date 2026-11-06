// Package cmaapi is the dispatcher's client for the CMA API: the work queue
// of one environment and the state of that environment's sessions. Every
// call is authorized with the environment key.
package cmaapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
)

// WorkItem is one claimed unit of work: a session that needs a sandbox.
type WorkItem struct {
	ID        string
	SessionID string
	Secret    string
}

// SessionState is what the API says about a session.
type SessionState int

const (
	// SessionAlive covers every session the API returns that is not terminated.
	SessionAlive SessionState = iota
	// SessionEnded is a session whose status is terminated.
	SessionEnded
	// SessionGone is a session the API says, by name, it does not have.
	SessionGone
)

type Client struct {
	api           anthropic.Client
	environmentID string
}

// New builds a client that only ever dials baseURL. A request attempt that
// has not answered within requestTimeout fails with context.DeadlineExceeded.
// The SDK may retry it like any other failed attempt. A hung upstream costs at
// most three requestTimeouts per call plus the SDK's backoff. requestTimeout
// must be positive. Zero would let a hung upstream block every call forever.
//
// The client drops the user name and password url.Parse reads in baseURL.
// Callers log the error of a failed call, and that error quotes the request's
// URL. Go would send neither: every request already has an Authorization
// header, the key's.
func New(baseURL, environmentID, environmentKey string, requestTimeout time.Duration) *Client {
	if requestTimeout <= 0 {
		panic("cmaapi.New: requestTimeout must be positive")
	}
	if u, err := url.Parse(baseURL); err == nil && u.User != nil {
		u.User = nil
		baseURL = u.String()
	}
	return &Client{
		api: anthropic.NewClient(
			option.WithoutEnvironmentDefaults(),
			option.WithBaseURL(baseURL),
			option.WithAuthToken(environmentKey),
			option.WithMaxRetries(2),
			option.WithRequestTimeout(requestTimeout),
		),
		environmentID: environmentID,
	}
}

// VerifyEnvironment fails unless the key can read the environment's queue stats.
func (c *Client) VerifyEnvironment(ctx context.Context) error {
	_, err := c.api.Beta.Environments.Work.Stats(ctx, c.environmentID, anthropic.BetaEnvironmentWorkStatsParams{})
	if err != nil {
		return fmt.Errorf("environment %s: work stats: %w", c.environmentID, withoutKey(err))
	}
	return nil
}

// pollBlockMs is how long one poll blocks server-side waiting for work.
const pollBlockMs = 999

// Poll long-polls the queue once. It returns nil when no work arrived.
// Items not acked within reclaimAfter are handed out again by the queue.
func (c *Client) Poll(ctx context.Context, workerID string, reclaimAfter time.Duration) (*WorkItem, error) {
	work, err := c.api.Beta.Environments.Work.Poll(ctx, c.environmentID, anthropic.BetaEnvironmentWorkPollParams{
		AnthropicWorkerID:  param.NewOpt(workerID),
		BlockMs:            param.NewOpt(int64(pollBlockMs)),
		ReclaimOlderThanMs: param.NewOpt(reclaimAfter.Milliseconds()),
	})
	if err != nil {
		return nil, withoutKey(err)
	}
	if work == nil || work.ID == "" {
		return nil, nil
	}
	return &WorkItem{ID: work.ID, SessionID: work.Data.ID, Secret: work.Secret}, nil
}

func (c *Client) Ack(ctx context.Context, workID string) error {
	_, err := c.api.Beta.Environments.Work.Ack(ctx, workID, anthropic.BetaEnvironmentWorkAckParams{EnvironmentID: c.environmentID})
	return withoutKey(err)
}

// SessionState asks the API about sessionID. Only the API's own 404 or a
// terminated status count as gone; any other failure is returned as an error
// so the caller keeps the session's sandbox.
func (c *Client) SessionState(ctx context.Context, sessionID string) (SessionState, error) {
	session, err := c.api.Beta.Sessions.Get(ctx, sessionID, anthropic.BetaSessionGetParams{})
	if isGone(err, sessionID) {
		return SessionGone, nil
	}
	if err != nil {
		return SessionAlive, withoutKey(err)
	}
	if session.Status == anthropic.BetaManagedAgentsSessionStatusTerminated {
		return SessionEnded, nil
	}
	return SessionAlive, nil
}

// isGone reports whether err is the API's own answer that sessionID does not
// exist. A proxy can answer 404 too, and the API answers not_found_error with
// another message to an organization without the beta.
func isGone(err error, sessionID string) bool {
	var apiErr *anthropic.Error
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusNotFound || apiErr.Type() != anthropic.ErrorTypeNotFoundError {
		return false
	}
	var body struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	return json.Unmarshal([]byte(apiErr.RawJSON()), &body) == nil && body.Error.Message == "Session not found: "+sessionID
}
