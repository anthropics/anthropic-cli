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
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/lib/environments"
	"github.com/urfave/cli/v3"
	"go.opentelemetry.io/otel/trace"
)

var workerCommand = cli.Command{
	Name:     "beta:worker",
	Category: "SELF-HOSTED",
	Usage:    "Run a self-hosted environment worker (poll for work and/or run tools).",
	Suggest:  true,
	Commands: []*cli.Command{
		workerPollCommandDef(),
		workerRunCommandDef(),
	},
}

// workerPollCommandDef returns a fresh `poll` subcommand definition, for the
// same reason as workerRunCommandDef.
func workerPollCommandDef() *cli.Command {
	return &cli.Command{
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
			unrestrictedPathsFlag(),
			&cli.DurationFlag{Name: "max-idle", Value: anthropic.DefaultMaxIdle, Usage: "stop this long after the session goes idle with stop_reason end_turn; 0 = no timeout"},
			&cli.StringFlag{Name: "log-format", Value: "text"},
		},
		Action:          handleWorkerPoll,
		HideHelpCommand: true,
	}
}

// workerRunCommandDef returns a fresh `run` subcommand definition — urfave/cli
// caches flag state on a Command across runs, so tests need their own instance.
func workerRunCommandDef() *cli.Command {
	return &cli.Command{
		Name:    "run",
		Usage:   "Attach to a session and execute agent.tool_use events locally, authorized by an environment key or a work secret. Intended as a container ENTRYPOINT.",
		Suggest: true,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "session-id", Required: true, Sources: cli.EnvVars("ANTHROPIC_SESSION_ID")},
			&cli.StringFlag{Name: "environment-key", Sources: cli.EnvVars("ANTHROPIC_ENVIRONMENT_KEY"), Usage: "environment key that authorizes the run; optional when a work secret is provided"},
			&cli.StringFlag{Name: "work-secret-file", Sources: cli.EnvVars("ANTHROPIC_WORK_SECRET_FILE"), Usage: "path to a file whose contents are the work item's secret; keeps the secret out of the environment that child processes inherit"},
			&cli.StringFlag{Name: "work-id", Required: true, Sources: cli.EnvVars("ANTHROPIC_WORK_ID")},
			&cli.StringFlag{Name: "environment-id", Required: true, Sources: cli.EnvVars("ANTHROPIC_ENVIRONMENT_ID")},
			&cli.StringFlag{Name: "base-url", Sources: cli.EnvVars("ANTHROPIC_BASE_URL")},
			&cli.StringFlag{Name: "workdir", Value: "."},
			unrestrictedPathsFlag(),
			&cli.DurationFlag{Name: "max-idle", Value: anthropic.DefaultMaxIdle, Usage: "stop this long after the session goes idle with stop_reason end_turn; 0 = no timeout"},
			&cli.StringFlag{Name: "log-format", Value: "text"},
		},
		Action:          handleWorkerRun,
		HideHelpCommand: true,
	}
}

// unrestrictedPathsFlag stays parseable, hidden from help, so an invocation
// that still passes it gets errUnrestrictedPaths instead of an unknown-flag
// usage error.
func unrestrictedPathsFlag() cli.Flag {
	return &cli.BoolFlag{Name: "unrestricted-paths", Hidden: true, Usage: "no longer supported"}
}

var errUnrestrictedPaths = errors.New("--unrestricted-paths is no longer supported: the file tools always stay inside --workdir and the session's attached memory stores; remove the flag")

func handleWorkerPoll(ctx context.Context, cmd *cli.Command) error {
	if cmd.Bool("unrestricted-paths") {
		return errUnrestrictedPaths
	}
	logger := newWorkerLogger(cmd.String("log-format"))
	client := newWorkerClient(extraClientFlagsFromCmd(cmd))

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// One poll process handles many sessions, one after another, so its spans are
	// not in the trace of the process that started it: each work item gets a
	// root span of its own instead.
	ctx = trace.ContextWithSpanContext(ctx, trace.SpanContext{})

	maxIdle := cmd.Duration("max-idle")

	// With no --on-work script (or "in-process"), the SDK's EnvironmentWorker
	// owns the whole loop: poll for work, set up the AgentToolContext +
	// download the session agent's skills, run a SessionToolRunner while
	// heartbeating the work-item lease, force-stop the work on exit, and loop
	// to the next item. The environment key authorizes both the work-poll
	// calls and the session-level calls.
	if script := cmd.String("on-work"); script == "" || script == "in-process" {
		worker := environments.NewEnvironmentWorker(client, environments.EnvironmentWorkerOptions{
			EnvironmentID:  cmd.String("environment-id"),
			EnvironmentKey: cmd.String("environment-key"),
			WorkerID:       cmd.String("worker-id"),
			Workdir:        cmd.String("workdir"),
			MaxIdle:        &maxIdle,
			Logger:         logger,
		})
		return worker.Run(ctx)
	}

	// With an --on-work script, use the WorkPoller and exec
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
	if cmd.Bool("unrestricted-paths") {
		return errUnrestrictedPaths
	}
	environmentKey, workSecret, err := workerRunCredentials(cmd)
	if err != nil {
		return err
	}

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
	// item instead of a poll loop. The worker must build the toolset: only it
	// learns the session's memory-store mount paths the file tools may reach.
	maxIdle := cmd.Duration("max-idle")
	worker := environments.NewEnvironmentWorker(client, environments.EnvironmentWorkerOptions{
		Workdir: cmd.String("workdir"),
		MaxIdle: &maxIdle,
		Logger:  logger,
	})
	return worker.HandleItem(ctx, environments.HandleItemOptions{
		WorkID:         cmd.String("work-id"),
		EnvironmentID:  cmd.String("environment-id"),
		SessionID:      cmd.String("session-id"),
		EnvironmentKey: environmentKey,
		WorkSecret:     workSecret,
	})
}

// workerRunCredentials resolves what authorizes a `run` invocation: the
// environment key flag, and the whitespace-trimmed contents of
// --work-secret-file as the work item's secret. A set file flag must yield a
// secret, and at least one credential must be present (counting the
// ANTHROPIC_WORK_SECRET variable the SDK reads on its own).
func workerRunCredentials(cmd *cli.Command) (environmentKey, workSecret string, err error) {
	environmentKey = cmd.String("environment-key")
	if path := cmd.String("work-secret-file"); path != "" {
		if err := refuseWorldAccessibleSecretFile(path); err != nil {
			return "", "", err
		}
		buf, err := os.ReadFile(path)
		if err != nil {
			return "", "", fmt.Errorf("read work secret file: %w", err)
		}
		workSecret = strings.TrimSpace(string(buf))
		if workSecret == "" {
			return "", "", fmt.Errorf("work secret file %s is empty", path)
		}
	}
	if environmentKey == "" && workSecret == "" && os.Getenv("ANTHROPIC_WORK_SECRET") == "" {
		return "", "", errors.New("provide an environment key (--environment-key or ANTHROPIC_ENVIRONMENT_KEY) or a work secret (--work-secret-file, ANTHROPIC_WORK_SECRET_FILE, or ANTHROPIC_WORK_SECRET)")
	}
	return environmentKey, workSecret, nil
}

// refuseWorldAccessibleSecretFile rejects a work secret file that every user
// on the machine can read or write. Group access stays allowed — Kubernetes
// deployments write the token 0440 under another owner and the worker reads
// it through the group bit. Skipped on Windows, whose mode bits are synthetic.
func refuseWorldAccessibleSecretFile(path string) error {
	if runtime.GOOS == "windows" {
		return nil
	}
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat work secret file: %w", err)
	}
	if mode := info.Mode().Perm(); mode&0o006 != 0 {
		return fmt.Errorf("work secret file %s (mode %04o) is readable or writable by every user on this machine — run chmod o-rw on it", path, mode)
	}
	return nil
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
