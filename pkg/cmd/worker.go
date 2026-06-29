package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/lib/environments"
	"github.com/anthropics/anthropic-sdk-go/tools/agenttoolset"
	"github.com/urfave/cli/v3"
)

var workerCommand = cli.Command{
	Name:     "beta:worker",
	Category: "SELF-HOSTED",
	Usage:    "Run a self-hosted environment worker (poll for work and/or run tools).",
	Suggest:  true,
	Commands: []*cli.Command{
		&workerPollCommand,
		&workerRunCommand,
	},
}

var workerPollCommand = cli.Command{
	Name:    "poll",
	Usage:   "Long-poll an environment for work and run an in-process session tool runner for each session.",
	Suggest: true,
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "environment-id", Required: true, Sources: cli.EnvVars("ANTHROPIC_ENVIRONMENT_ID")},
		&cli.StringFlag{Name: "environment-key", Required: true, Sources: cli.EnvVars("ANTHROPIC_ENVIRONMENT_KEY")},
		&cli.StringFlag{Name: "worker-id", Sources: cli.EnvVars("ANTHROPIC_WORKER_ID")},
		&cli.StringFlag{Name: "base-url", Sources: cli.EnvVars("ANTHROPIC_BASE_URL")},
		&cli.StringFlag{Name: "on-work", Usage: "Script to exec for each work item instead of the in-process runner. Receives ANTHROPIC_{WORK_ID,ENVIRONMENT_ID,SESSION_ID,ENVIRONMENT_KEY} in env and the work JSON on stdin. Empty or 'in-process' = built-in runner."},
		&cli.StringFlag{Name: "workdir", Value: "."},
		&cli.BoolFlag{Name: "unrestricted-paths", Usage: "let the file tools read/write outside the workdir (the workdir check is a guardrail for the file tools only, not a sandbox, and is not respected by bash)"},
		&cli.DurationFlag{Name: "max-idle", Value: anthropic.DefaultMaxIdle, Usage: "stop this long after the session goes idle with stop_reason end_turn; 0 = no timeout"},
		&cli.StringFlag{Name: "log-format", Value: "text"},
	},
	Action:          handleWorkerPoll,
	HideHelpCommand: true,
}

var workerRunCommand = cli.Command{
	Name:    "run",
	Usage:   "Attach to a session and execute agent.tool_use events locally. Intended as a container ENTRYPOINT.",
	Suggest: true,
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "session-id", Required: true, Sources: cli.EnvVars("ANTHROPIC_SESSION_ID")},
		&cli.StringFlag{Name: "environment-key", Required: true, Sources: cli.EnvVars("ANTHROPIC_ENVIRONMENT_KEY")},
		&cli.StringFlag{Name: "work-id", Required: true, Sources: cli.EnvVars("ANTHROPIC_WORK_ID")},
		&cli.StringFlag{Name: "environment-id", Required: true, Sources: cli.EnvVars("ANTHROPIC_ENVIRONMENT_ID")},
		&cli.StringFlag{Name: "base-url", Sources: cli.EnvVars("ANTHROPIC_BASE_URL")},
		&cli.StringFlag{Name: "workdir", Value: "."},
		&cli.BoolFlag{Name: "unrestricted-paths", Usage: "let the file tools read/write outside the workdir (the workdir check is a guardrail for the file tools only, not a sandbox, and is not respected by bash)"},
		&cli.DurationFlag{Name: "max-idle", Value: anthropic.DefaultMaxIdle, Usage: "stop this long after the session goes idle with stop_reason end_turn; 0 = no timeout"},
		&cli.StringFlag{Name: "log-format", Value: "text"},
	},
	Action:          handleWorkerRun,
	HideHelpCommand: true,
}

func handleWorkerPoll(ctx context.Context, cmd *cli.Command) error {
	logger := newWorkerLogger(cmd.String("log-format"))
	client := newWorkerClient(extraClientFlagsFromCmd(cmd))

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	maxIdle := cmd.Duration("max-idle")

	// With no --on-work script (or "in-process"), the SDK's EnvironmentWorker
	// owns the whole loop: poll for work, set up the AgentToolContext +
	// download the session agent's skills, run a SessionToolRunner while
	// heartbeating the work-item lease, force-stop the work on exit, and loop
	// to the next item. The environment key authorizes both the work-poll
	// calls and the session-level calls.
	if script := cmd.String("on-work"); script == "" || script == "in-process" {
		worker := environments.NewEnvironmentWorker(client, environments.EnvironmentWorkerOptions{
			EnvironmentID:     cmd.String("environment-id"),
			EnvironmentKey:    cmd.String("environment-key"),
			WorkerID:          cmd.String("worker-id"),
			Workdir:           cmd.String("workdir"),
			UnrestrictedPaths: cmd.Bool("unrestricted-paths"),
			MaxIdle:           &maxIdle,
			Logger:            logger,
		})
		return worker.Run(ctx)
	}

	// With an --on-work script, use the control-plane-only WorkPoller and exec
	// the script for each claimed item; the script (typically launching a
	// container that runs `ant beta:worker run`) is responsible for servicing
	// the session and heartbeating the lease.
	poller := environments.NewWorkPoller(ctx, client, environments.WorkPollerOptions{
		EnvironmentID:  cmd.String("environment-id"),
		EnvironmentKey: cmd.String("environment-key"),
		WorkerID:       cmd.String("worker-id"),
		Logger:         logger,
	})
	defer poller.Close()

	runScript := newOnWorkRunner(cmd, logger)

	for poller.Next() {
		work := poller.Current()
		if err := runScript(ctx, work); err != nil {
			if errors.Is(err, errFatalWorker) {
				return err
			}
			logger.WarnContext(ctx, "work handler exited with error",
				slog.String("work_id", work.ID), slog.Any("error", err))
		}
	}
	if err := poller.Err(); err != nil {
		return fmt.Errorf("poller: %w", err)
	}
	return nil
}

func handleWorkerRun(ctx context.Context, cmd *cli.Command) error {
	logger := newWorkerLogger(cmd.String("log-format"))
	client := newWorkerClient(extraClientFlagsFromCmd(cmd))

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// `beta:worker run` is handed an already-claimed `session` work item by
	// `beta:worker poll --on-work` (forwarded via env vars / flags). Hand the
	// per-item flow — build the AgentToolContext + download the session agent's
	// skills, run a SessionToolRunner while heartbeating the work-item lease,
	// force-stop the work on exit — to the SDK's EnvironmentWorker. That's the
	// same composition `beta:worker poll`'s in-process path runs, just for a single
	// item instead of a poll loop.
	env := &agenttoolset.AgentToolContext{
		Workdir:           cmd.String("workdir"),
		UnrestrictedPaths: cmd.Bool("unrestricted-paths"),
	}
	toolset := agenttoolset.BetaAgentToolset20260401(env)
	defer agenttoolset.CloseAll(toolset)

	maxIdle := cmd.Duration("max-idle")
	worker := environments.NewEnvironmentWorker(client, environments.EnvironmentWorkerOptions{
		Tools:             toolset,
		Workdir:           cmd.String("workdir"),
		UnrestrictedPaths: cmd.Bool("unrestricted-paths"),
		MaxIdle:           &maxIdle,
		Logger:            logger,
	})
	return worker.HandleItem(ctx, environments.HandleItemOptions{
		WorkID:         cmd.String("work-id"),
		EnvironmentID:  cmd.String("environment-id"),
		SessionID:      cmd.String("session-id"),
		EnvironmentKey: cmd.String("environment-key"),
	})
}

// errFatalWorker signals the poll loop should exit instead of looping.
// Used for operator-configuration errors (an unrunnable --on-work script
// path) where retrying the next claimed work item would just burn through
// the queue producing identical failures.
var errFatalWorker = errors.New("fatal worker error")

// newOnWorkRunner returns the per-work handler for `beta:worker poll --on-work`: it
// execs the given script with the same env vars `beta:worker run` reads, so the
// script can launch a container or subprocess and re-enter via run.
func newOnWorkRunner(cmd *cli.Command, logger *slog.Logger) func(context.Context, *anthropic.BetaSelfHostedWork) error {
	script := cmd.String("on-work")
	environmentKey := cmd.String("environment-key")

	return func(ctx context.Context, work *anthropic.BetaSelfHostedWork) error {
		if work.Data.Type != "session" {
			// Don't hand the script a work payload it can't recognize (it would
			// see an empty ANTHROPIC_SESSION_ID and have no way to tell that
			// apart from a malformed session payload).
			logger.WarnContext(ctx, "unsupported work type, skipping",
				slog.String("type", string(work.Data.Type)),
				slog.String("work_id", work.ID))
			return nil
		}
		sessionID := work.Data.ID
		c := exec.CommandContext(ctx, script)
		c.Env = append(os.Environ(),
			"ANTHROPIC_WORK_ID="+work.ID,
			"ANTHROPIC_ENVIRONMENT_ID="+work.EnvironmentID,
			"ANTHROPIC_SESSION_ID="+sessionID,
			"ANTHROPIC_ENVIRONMENT_KEY="+environmentKey,
		)
		// Fold child stdout into our stderr so a stderr-based log consumer sees
		// one unified stream. `ant beta:worker poll` produces no stdout itself, so
		// we're not stepping on a pipeable output channel.
		c.Stdout = os.Stderr
		c.Stderr = os.Stderr
		c.Cancel = func() error { return c.Process.Signal(syscall.SIGTERM) }
		c.WaitDelay = 30 * time.Second
		if buf, err := json.Marshal(work); err == nil {
			c.Stdin = bytes.NewReader(buf)
		} else {
			logger.WarnContext(ctx, "marshal work for stdin failed, script will see empty stdin",
				slog.String("work_id", work.ID), slog.Any("error", err))
		}
		logger.InfoContext(ctx, "spawning on-work script", slog.String("script", script), slog.String("work_id", work.ID))
		err := c.Run()
		if err != nil {
			// exec.Error means the script could not be started at all (file not
			// found, not executable) — that's operator configuration, retrying
			// every claim will just burn through the queue. Fail hard.
			var execErr *exec.Error
			if errors.As(err, &execErr) {
				return fmt.Errorf("%w: on-work script %q is unrunnable: %w", errFatalWorker, script, err)
			}
			logger.WarnContext(ctx, "on-work script exited non-zero", slog.String("work_id", work.ID), slog.Any("error", err))
		}
		return err
	}
}

// newWorkerClient builds the SDK client the worker uses for all calls. The
// environment helpers attach the environment key per-request, so the client
// itself carries no credential, only the client-level flag overrides.
func newWorkerClient(extra extraClientFlags) anthropic.Client {
	return anthropic.NewClient(extra.requestOptions()...)
}

func newWorkerLogger(format string) *slog.Logger {
	var h slog.Handler
	if format == "json" {
		h = slog.NewJSONHandler(os.Stderr, nil)
	} else {
		h = slog.NewTextHandler(os.Stderr, nil)
	}
	return slog.New(h)
}
