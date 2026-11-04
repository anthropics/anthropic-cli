package cmd

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-cli/v2/internal/apiquery"
	"github.com/anthropics/anthropic-cli/v2/internal/requestflag"
	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var betaOrganizationTelemetryDestinationsCreate = requestflag.WithInnerFlags(cli.Command{
	Name:    "create",
	Usage:   "Register an OpenTelemetry destination for the organization.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "endpoint-url",
			Usage:    "The base URL of your collector's OTLP/HTTP endpoint, for example `https://gateway.example.com`.\n\nThe URL must use `https` and a hostname, not an IP address. The port, if given, must be 443. The URL must not contain a query string, a fragment, or a username and password. At most 2,048 characters.\n\nA path prefix, such as `/otlp`, is allowed. Don't include the per-signal path (`/v1/traces`), because the API appends it on delivery. If your collector takes traces on another path, also set `traces.endpoint_url`.",
			Required: true,
			BodyPath: "endpoint_url",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "workspaces",
			Usage:    "Which of the organization's workspaces a telemetry destination covers, as a create or update request sets it.",
			Required: true,
			BodyPath: "workspaces",
		},
		&requestflag.Flag[*string]{
			Name:     "display-name",
			Usage:    "Human-readable name for the destination. Optional. 1 to 255 characters.",
			BodyPath: "display_name",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "headers",
			Usage:    "Request headers that the API sends with every delivery, typically your collector's authentication header (for example `authorization`).\n\nThe API never returns the values. Reads show only `header_names`.\n\nAt most 20 headers. Names are at most 128 characters. Two names can't differ only in case. Values are printable ASCII, at most 4,096 characters, and can't begin or end with a space. Names and values together are at most 16,384 bytes.\n\nReserved names are rejected, for example `host`, `cookie`, `user-agent`, and names that start with `content-`, `proxy-`, `x-forwarded-`, or `anthropic-`. For this check, `_` in a name counts as `-`.",
			BodyPath: "headers",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "traces",
			Usage:    "The settings for the traces that a telemetry destination receives, as a create or update request sets them. An omitted field takes its default on create and keeps its stored value on update.",
			BodyPath: "traces",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `telemetry-destinations-2026-08-11` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationTelemetryDestinationsCreate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"workspaces": {
		&requestflag.InnerFlag[string]{
			Name:       "workspaces.type",
			Usage:      `Allowed values: "all", "include", "exclude".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[any]{
			Name:       "workspaces.workspace-ids",
			InnerField: "workspace_ids",
		},
	},
	"traces": {
		&requestflag.InnerFlag[*string]{
			Name:       "traces.content-mode",
			Usage:      "What the spans that the API exports to a telemetry destination contain.\n\nSee the [OpenTelemetry export guide](https://platform.claude.com/docs/en/manage-claude/opentelemetry#choose-a-content-mode) to learn what the spans contain in each content mode.",
			InnerField: "content_mode",
		},
		&requestflag.InnerFlag[bool]{
			Name:       "traces.enabled",
			Usage:      "Whether the destination receives traces. Defaults to true on create. It cannot be null.",
			InnerField: "enabled",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "traces.endpoint-url",
			Usage:      "The full URL of your collector's traces endpoint, for a collector that takes traces on a path other than `/v1/traces`. The API posts traces to it exactly as given and appends nothing.\n\nThe same rules as for `endpoint_url` apply, except that any path is allowed. It must have the same scheme, host, and port as `endpoint_url`. At most 2,048 characters.\n\nOn create, `null` is the same as leaving it out. On update, `null` removes it. The API then posts traces to `endpoint_url` with `/v1/traces` appended. A change of this URL alone keeps the registered headers.",
			InnerField: "endpoint_url",
		},
		&requestflag.InnerFlag[*int64]{
			Name:       "traces.max-content-bytes",
			Usage:      "The most bytes of content that one exported span can contain, from 1 up to the API's maximum. `null` sets no bound, so only the API's maximum applies.\n\nOn create, leaving it out is the same as `null`. The error for a value above the maximum states the maximum.",
			InnerField: "max_content_bytes",
		},
	},
})

var betaOrganizationTelemetryDestinationsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Read a destination's configuration, state, and most recent delivery results.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "destination-id",
			Usage:     "ID of the destination (`tdst_...`).",
			Required:  true,
			PathParam: "destination_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `telemetry-destinations-2026-08-11` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationTelemetryDestinationsRetrieve,
	HideHelpCommand: true,
}

var betaOrganizationTelemetryDestinationsUpdate = requestflag.WithInnerFlags(cli.Command{
	Name:    "update",
	Usage:   "Change a telemetry destination's settings, its state, or which workspaces it\ncovers.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:        "destination-id",
			Usage:       "ID of the destination (`tdst_...`).",
			Required:    true,
			PathParam:   "destination_id",
			DataAliases: []string{"id"},
		},
		&requestflag.Flag[*string]{
			Name:     "display-name",
			Usage:    "New human-readable name, 1 to 255 characters. Null clears it.",
			BodyPath: "display_name",
		},
		&requestflag.Flag[string]{
			Name:     "endpoint-url",
			Usage:    "The new base URL of the destination's OTLP/HTTP endpoint. The same rules as on create apply.\n\nA change of hostname or port first removes every registered header, so send the headers for the new endpoint in the same request. When the destination has a `traces.endpoint_url`, such a change must also set `traces.endpoint_url` to a URL with the new hostname and port, or to `null`. A change of the path alone keeps the headers.",
			BodyPath: "endpoint_url",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "headers",
			Usage:    "Headers to set or remove. The API merges them into the destination's headers by name.\n\nNames compare case-insensitively. A string value sets or replaces that header. A `null` value removes it. Headers that you omit don't change. The field itself can't be `null`.\n\nA request that sets a header value must also include `endpoint_url`. The current URL is accepted.\n\nThe resulting set must stay within the same limits as on create. When the request sets a value, each stored value that it doesn't set counts as 4,096 bytes toward the 16,384-byte limit, because you can't read its real length. So to change one header of a destination that has four or more other headers, send the values of the others again in the same request.",
			BodyPath: "headers",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "state",
			Usage:    "The state to set on a telemetry destination, by `type` alone.\n\nA caller can't set the `reason` of a disabled state. Only the API sets it.",
			BodyPath: "state",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "traces",
			Usage:    "The settings for the traces that a telemetry destination receives, as a create or update request sets them. An omitted field takes its default on create and keeps its stored value on update.",
			BodyPath: "traces",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "workspaces",
			Usage:    "Which of the organization's workspaces a telemetry destination covers, as a create or update request sets it.",
			BodyPath: "workspaces",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `telemetry-destinations-2026-08-11` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationTelemetryDestinationsUpdate,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"state": {
		&requestflag.InnerFlag[string]{
			Name:       "state.type",
			Usage:      `Allowed values: "enabled", "disabled".`,
			InnerField: "type",
		},
	},
	"traces": {
		&requestflag.InnerFlag[*string]{
			Name:       "traces.content-mode",
			Usage:      "What the spans that the API exports to a telemetry destination contain.\n\nSee the [OpenTelemetry export guide](https://platform.claude.com/docs/en/manage-claude/opentelemetry#choose-a-content-mode) to learn what the spans contain in each content mode.",
			InnerField: "content_mode",
		},
		&requestflag.InnerFlag[bool]{
			Name:       "traces.enabled",
			Usage:      "Whether the destination receives traces. Defaults to true on create. It cannot be null.",
			InnerField: "enabled",
		},
		&requestflag.InnerFlag[*string]{
			Name:       "traces.endpoint-url",
			Usage:      "The full URL of your collector's traces endpoint, for a collector that takes traces on a path other than `/v1/traces`. The API posts traces to it exactly as given and appends nothing.\n\nThe same rules as for `endpoint_url` apply, except that any path is allowed. It must have the same scheme, host, and port as `endpoint_url`. At most 2,048 characters.\n\nOn create, `null` is the same as leaving it out. On update, `null` removes it. The API then posts traces to `endpoint_url` with `/v1/traces` appended. A change of this URL alone keeps the registered headers.",
			InnerField: "endpoint_url",
		},
		&requestflag.InnerFlag[*int64]{
			Name:       "traces.max-content-bytes",
			Usage:      "The most bytes of content that one exported span can contain, from 1 up to the API's maximum. `null` sets no bound, so only the API's maximum applies.\n\nOn create, leaving it out is the same as `null`. The error for a value above the maximum states the maximum.",
			InnerField: "max_content_bytes",
		},
	},
	"workspaces": {
		&requestflag.InnerFlag[string]{
			Name:       "workspaces.type",
			Usage:      `Allowed values: "all", "include", "exclude".`,
			InnerField: "type",
		},
		&requestflag.InnerFlag[any]{
			Name:       "workspaces.workspace-ids",
			InnerField: "workspace_ids",
		},
	},
})

var betaOrganizationTelemetryDestinationsList = cli.Command{
	Name:    "list",
	Usage:   "List the organization's destinations, newest first, optionally only those that\ncover given workspaces.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of destinations to return per page, between 1 and 1,000. Defaults to 20.",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "page",
			Usage:     "Opaque cursor from a previous response's `next_page`. Leave it out to get the first page.",
			QueryPath: "page",
		},
		&requestflag.Flag[[]string]{
			Name:      "workspace-id",
			Usage:     "Return only the destinations that cover at least one of these workspaces (`wrkspc_...`).\n\nSend at most 100 IDs, as `workspace_ids[]`. A destination covers a workspace when its `workspaces.type` is `all`, when its `include` list has the workspace, or when its `exclude` list doesn't have it. Destinations in either state can match. An ID that isn't the ID of one of the organization's workspaces doesn't match any destination, not even one whose type is `all`.",
			QueryPath: "workspace_ids",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `telemetry-destinations-2026-08-11` in this header.",
			HeaderPath: "anthropic-beta",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleBetaOrganizationTelemetryDestinationsList,
	HideHelpCommand: true,
}

var betaOrganizationTelemetryDestinationsDelete = cli.Command{
	Name:    "delete",
	Usage:   "Delete a destination and its workspace list.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "destination-id",
			Usage:     "ID of the destination (`tdst_...`).",
			Required:  true,
			PathParam: "destination_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `telemetry-destinations-2026-08-11` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationTelemetryDestinationsDelete,
	HideHelpCommand: true,
}

var betaOrganizationTelemetryDestinationsTest = cli.Command{
	Name:    "test",
	Usage:   "Send one test span, named `anthropic.test`, to the destination with its\nregistered headers, once and without retries, and report the destination's\nanswer. The response is 200 whether or not the destination answered with a 2xx\nstatus: its `result` has the type `success`, or the type `failure` and the\nreason in `error`. The destination's configuration, `state` and\n`delivery_health` do not change. A disabled destination can be tested too. Wait\n10 seconds between tests of one destination; a request that comes sooner may be\nanswered with status 429.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "destination-id",
			Usage:     "ID of the destination to test (starts with `tdst_`).",
			Required:  true,
			PathParam: "destination_id",
		},
		&requestflag.Flag[[]string]{
			Name:       "beta",
			Usage:      "This endpoint is in beta: requests must send `telemetry-destinations-2026-08-11` in this header.",
			HeaderPath: "anthropic-beta",
		},
	},
	Action:          handleBetaOrganizationTelemetryDestinationsTest,
	HideHelpCommand: true,
}

func handleBetaOrganizationTelemetryDestinationsCreate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		ApplicationJSON,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaOrganizationTelemetryDestinationNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.TelemetryDestinations.New(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "beta:organization:telemetry-destinations create",
		Transform:      transform,
	})
}

func handleBetaOrganizationTelemetryDestinationsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("destination-id") && len(unusedArgs) > 0 {
		cmd.Set("destination-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaOrganizationTelemetryDestinationGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.TelemetryDestinations.Get(
		ctx,
		cmd.Value("destination-id").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "beta:organization:telemetry-destinations retrieve",
		Transform:      transform,
	})
}

func handleBetaOrganizationTelemetryDestinationsUpdate(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("destination-id") && len(unusedArgs) > 0 {
		cmd.Set("destination-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		ApplicationJSON,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaOrganizationTelemetryDestinationUpdateParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.TelemetryDestinations.Update(
		ctx,
		cmd.Value("destination-id").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "beta:organization:telemetry-destinations update",
		Transform:      transform,
	})
}

func handleBetaOrganizationTelemetryDestinationsList(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaOrganizationTelemetryDestinationListParams{}

	format := "explore"
	explicitFormat := cmd.Root().IsSet("format")
	if explicitFormat {
		format = cmd.Root().String("format")
	}
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Beta.Organization.TelemetryDestinations.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:telemetry-destinations list",
			Transform:      transform,
		})
	} else {
		iter := client.Beta.Organization.TelemetryDestinations.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "beta:organization:telemetry-destinations list",
			Transform:      transform,
		})
	}
}

func handleBetaOrganizationTelemetryDestinationsDelete(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("destination-id") && len(unusedArgs) > 0 {
		cmd.Set("destination-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaOrganizationTelemetryDestinationDeleteParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.TelemetryDestinations.Delete(
		ctx,
		cmd.Value("destination-id").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "beta:organization:telemetry-destinations delete",
		Transform:      transform,
	})
}

func handleBetaOrganizationTelemetryDestinationsTest(ctx context.Context, cmd *cli.Command) error {
	client := anthropic.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("destination-id") && len(unusedArgs) > 0 {
		cmd.Set("destination-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	params := anthropic.BetaOrganizationTelemetryDestinationTestParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Beta.Organization.TelemetryDestinations.Test(
		ctx,
		cmd.Value("destination-id").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "beta:organization:telemetry-destinations test",
		Transform:      transform,
	})
}
