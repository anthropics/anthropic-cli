// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/anthropics/anthropic-cli/internal/autocomplete"
	"github.com/anthropics/anthropic-cli/internal/requestflag"
	docs "github.com/urfave/cli-docs/v3"
	"github.com/urfave/cli/v3"
)

var (
	Command            *cli.Command
	CommandErrorBuffer bytes.Buffer
)

func init() {
	Command = &cli.Command{
		Name:      "ant",
		Usage:     "CLI for the anthropic API",
		Suggest:   true,
		Version:   Version,
		ErrWriter: &CommandErrorBuffer,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug logging",
			},
			&cli.StringFlag{
				Name:        "base-url",
				DefaultText: "url",
				Usage:       "Override the base URL for API requests",
				Validator: func(baseURL string) error {
					return ValidateBaseURL(baseURL, "--base-url")
				},
			},
			&cli.StringFlag{
				Name:  "format",
				Usage: "The format for displaying response data (one of: " + strings.Join(OutputFormats, ", ") + ")",
				Value: "auto",
				Validator: func(format string) error {
					if !slices.Contains(OutputFormats, strings.ToLower(format)) {
						return fmt.Errorf("format must be one of: %s", strings.Join(OutputFormats, ", "))
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:  "format-error",
				Usage: "The format for displaying error data (one of: " + strings.Join(OutputFormats, ", ") + ")",
				Value: "auto",
				Validator: func(format string) error {
					if !slices.Contains(OutputFormats, strings.ToLower(format)) {
						return fmt.Errorf("format must be one of: %s", strings.Join(OutputFormats, ", "))
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:  "transform",
				Usage: "The GJSON transformation for data output.",
			},
			&cli.StringFlag{
				Name:  "transform-error",
				Usage: "The GJSON transformation for errors.",
			},
			&cli.BoolFlag{
				Name:    "raw-output",
				Aliases: []string{"r"},
				Usage:   "If the result is a string, print it without JSON quotes. This can be useful for making output transforms talk to non-JSON-based systems.",
			},
			&requestflag.Flag[string]{
				Name:    "api-key",
				Sources: cli.EnvVars("ANTHROPIC_API_KEY"),
			},
			&requestflag.Flag[string]{
				Name:    "auth-token",
				Sources: cli.EnvVars("ANTHROPIC_AUTH_TOKEN"),
			},
		},
		Commands: []*cli.Command{
			{
				Name:     "completions",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&completionsCreate,
				},
			},
			{
				Name:     "messages",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&messagesCreate,
					&messagesCountTokens,
				},
			},
			{
				Name:     "messages:batches",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&messagesBatchesCreate,
					&messagesBatchesRetrieve,
					&messagesBatchesList,
					&messagesBatchesDelete,
					&messagesBatchesCancel,
					&messagesBatchesResults,
				},
			},
			{
				Name:     "models",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&modelsRetrieve,
					&modelsList,
				},
			},
			{
				Name:     "beta:models",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&betaModelsRetrieve,
					&betaModelsList,
				},
			},
			{
				Name:     "beta:messages",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&betaMessagesCreate,
					&betaMessagesCountTokens,
				},
			},
			{
				Name:     "beta:messages:batches",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&betaMessagesBatchesCreate,
					&betaMessagesBatchesRetrieve,
					&betaMessagesBatchesList,
					&betaMessagesBatchesDelete,
					&betaMessagesBatchesCancel,
					&betaMessagesBatchesResults,
				},
			},
			{
				Name:     "beta:agents",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&betaAgentsCreate,
					&betaAgentsRetrieve,
					&betaAgentsUpdate,
					&betaAgentsList,
					&betaAgentsArchive,
				},
			},
			{
				Name:     "beta:agents:versions",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&betaAgentsVersionsList,
				},
			},
			{
				Name:     "beta:environments",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&betaEnvironmentsCreate,
					&betaEnvironmentsRetrieve,
					&betaEnvironmentsUpdate,
					&betaEnvironmentsList,
					&betaEnvironmentsDelete,
					&betaEnvironmentsArchive,
				},
			},
			{
				Name:     "beta:sessions",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&betaSessionsCreate,
					&betaSessionsRetrieve,
					&betaSessionsUpdate,
					&betaSessionsList,
					&betaSessionsDelete,
					&betaSessionsArchive,
				},
			},
			{
				Name:     "beta:sessions:events",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&betaSessionsEventsList,
					&betaSessionsEventsSend,
					&betaSessionsEventsStream,
				},
			},
			{
				Name:     "beta:sessions:resources",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&betaSessionsResourcesRetrieve,
					&betaSessionsResourcesUpdate,
					&betaSessionsResourcesList,
					&betaSessionsResourcesDelete,
					&betaSessionsResourcesAdd,
				},
			},
			{
				Name:     "beta:vaults",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&betaVaultsCreate,
					&betaVaultsRetrieve,
					&betaVaultsUpdate,
					&betaVaultsList,
					&betaVaultsDelete,
					&betaVaultsArchive,
				},
			},
			{
				Name:     "beta:vaults:credentials",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&betaVaultsCredentialsCreate,
					&betaVaultsCredentialsRetrieve,
					&betaVaultsCredentialsUpdate,
					&betaVaultsCredentialsList,
					&betaVaultsCredentialsDelete,
					&betaVaultsCredentialsArchive,
				},
			},
			{
				Name:     "beta:memory-stores",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&betaMemoryStoresCreate,
					&betaMemoryStoresRetrieve,
					&betaMemoryStoresUpdate,
					&betaMemoryStoresList,
					&betaMemoryStoresDelete,
					&betaMemoryStoresArchive,
				},
			},
			{
				Name:     "beta:memory-stores:memories",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&betaMemoryStoresMemoriesCreate,
					&betaMemoryStoresMemoriesRetrieve,
					&betaMemoryStoresMemoriesUpdate,
					&betaMemoryStoresMemoriesList,
					&betaMemoryStoresMemoriesDelete,
				},
			},
			{
				Name:     "beta:memory-stores:memory-versions",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&betaMemoryStoresMemoryVersionsRetrieve,
					&betaMemoryStoresMemoryVersionsList,
					&betaMemoryStoresMemoryVersionsRedact,
				},
			},
			{
				Name:     "beta:files",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&betaFilesList,
					&betaFilesDelete,
					&betaFilesDownload,
					&betaFilesRetrieveMetadata,
					&betaFilesUpload,
				},
			},
			{
				Name:     "beta:skills",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&betaSkillsCreate,
					&betaSkillsRetrieve,
					&betaSkillsList,
					&betaSkillsDelete,
				},
			},
			{
				Name:     "beta:skills:versions",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&betaSkillsVersionsCreate,
					&betaSkillsVersionsRetrieve,
					&betaSkillsVersionsList,
					&betaSkillsVersionsDelete,
				},
			},
			{
				Name:     "beta:user-profiles",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&betaUserProfilesCreate,
					&betaUserProfilesRetrieve,
					&betaUserProfilesUpdate,
					&betaUserProfilesList,
					&betaUserProfilesCreateEnrollmentURL,
				},
			},
			{
				Name:            "@manpages",
				Usage:           "Generate documentation for 'man'",
				UsageText:       "ant @manpages [-o ant.1] [--gzip]",
				Hidden:          true,
				Action:          generateManpages,
				HideHelpCommand: true,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "write manpages to the given folder",
						Value:   "man",
					},
					&cli.BoolFlag{
						Name:    "gzip",
						Aliases: []string{"z"},
						Usage:   "output gzipped manpage files to .gz",
						Value:   true,
					},
					&cli.BoolFlag{
						Name:    "text",
						Aliases: []string{"z"},
						Usage:   "output uncompressed text files",
						Value:   false,
					},
				},
			},
			{
				Name:            "__complete",
				Hidden:          true,
				HideHelpCommand: true,
				Action:          autocomplete.ExecuteShellCompletion,
			},
			{
				Name:            "@completion",
				Hidden:          true,
				HideHelpCommand: true,
				Action:          autocomplete.OutputCompletionScript,
			},
		},
		HideHelpCommand: true,
	}
}

func generateManpages(ctx context.Context, c *cli.Command) error {
	manpage, err := docs.ToManWithSection(Command, 1)
	if err != nil {
		return err
	}
	dir := c.String("output")
	err = os.MkdirAll(filepath.Join(dir, "man1"), 0755)
	if err != nil {
		// handle error
	}
	if c.Bool("text") {
		file, err := os.Create(filepath.Join(dir, "man1", "ant.1"))
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.WriteString(manpage); err != nil {
			return err
		}
	}
	if c.Bool("gzip") {
		file, err := os.Create(filepath.Join(dir, "man1", "ant.1.gz"))
		if err != nil {
			return err
		}
		defer file.Close()
		gzWriter := gzip.NewWriter(file)
		defer gzWriter.Close()
		_, err = gzWriter.Write([]byte(manpage))
		if err != nil {
			return err
		}
	}
	fmt.Printf("Wrote manpages to %s\n", dir)
	return nil
}
