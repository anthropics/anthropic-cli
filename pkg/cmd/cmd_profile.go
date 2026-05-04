package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/anthropics/anthropic-sdk-go/config"
	"github.com/urfave/cli/v3"
)

func init() {
	Command.Commands = append(Command.Commands, &cli.Command{
		Name:    "profile",
		Usage:   "Manage configuration profiles",
		Suggest: true,
		Commands: []*cli.Command{
			{
				Name:      "activate",
				Usage:     "Set the active profile",
				ArgsUsage: "<name>",
				Action:    profileActivate,
			},
			{
				Name:   "list",
				Usage:  "List configured profiles",
				Action: profileList,
			},
			{
				Name:      "get",
				Usage:     "Print the profile config (or one field)",
				ArgsUsage: "[<key>]",
				Action:    profileGet,
			},
			{
				Name:      "set",
				Usage:     "Set a field in the profile config",
				ArgsUsage: "<key> <value>",
				Action:    profileSet,
			},
		},
	})
}

func profileActivate(ctx context.Context, c *cli.Command) error {
	if c.Args().Len() != 1 {
		return fmt.Errorf("usage: ant profile activate <name>")
	}
	name := c.Args().First()
	dir := config.DefaultDir()
	if _, err := os.Stat(config.ProfilePath(dir, name)); err != nil {
		return fmt.Errorf("profile %q doesn't exist. Run `ant auth login --profile %s` to create it", name, name)
	}
	if err := config.SetActiveProfile(dir, name); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "✓ Active profile set to %q.\n", name)
	return nil
}

func profileList(ctx context.Context, c *cli.Command) error {
	dir := config.DefaultDir()
	active, _ := activeProfile(c)
	names, err := config.ListProfiles(dir)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		fmt.Fprintln(os.Stderr, "(no profiles — run `ant auth login`)")
		return nil
	}
	for _, name := range names {
		marker := " "
		if name == active {
			marker = "*"
		}
		fmt.Fprintf(os.Stdout, "%s %s\n", marker, name)
	}
	return nil
}

// configKeys derives the set of valid config field names from the struct's
// JSON tags, so `profile get` can distinguish a typo from a known-but-unset
// field without a hand-maintained list. `scope` is added explicitly — it
// lives nested under authentication in the on-disk shape but the CLI
// surfaces it as a flat profile key.
//
// IMPORTANT: `profile set` has its own hand-maintained switch over the
// writable subset of these keys. When adding a new settable field, update
// both configKeys() (read surface — implicit via reflection) AND the switch
// in profileSet (write surface — must be explicit). The reflection path
// will happily let `profile get <new_field>` work in isolation, which
// masks a missing `profile set` case until a user tries to write.
func configKeys() map[string]bool {
	t := reflect.TypeOf(config.Config{})
	keys := make(map[string]bool, t.NumField()+1)
	for i := 0; i < t.NumField(); i++ {
		tag := strings.SplitN(t.Field(i).Tag.Get("json"), ",", 2)[0]
		if tag != "" && tag != "-" {
			keys[tag] = true
		}
	}
	// version is file-format metadata owned by SaveProfile, not a user knob.
	delete(keys, "version")
	keys["scope"] = true
	return keys
}

func profileGet(ctx context.Context, c *cli.Command) error {
	profile, dir := activeProfile(c)
	cfg, err := loadProfileConfig(dir, profile)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal profile config: %w", err)
	}
	if c.Args().Len() == 0 {
		fmt.Fprintln(os.Stdout, string(data))
		return nil
	}
	key := c.Args().First()
	if !configKeys()[key] {
		return fmt.Errorf("no such field %q in profile config", key)
	}
	if key == "scope" {
		var scope string
		if cfg.AuthenticationInfo != nil && cfg.AuthenticationInfo.UserOAuth != nil {
			scope = cfg.AuthenticationInfo.UserOAuth.Scope
		}
		fmt.Fprintln(os.Stdout, scope)
		return nil
	}
	var m map[string]any
	_ = json.Unmarshal(data, &m)
	if v, ok := m[key]; ok {
		switch v.(type) {
		case map[string]any, []any:
			b, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal %s: %w", key, err)
			}
			fmt.Fprintln(os.Stdout, string(b))
		default:
			fmt.Fprintln(os.Stdout, v)
		}
	} else {
		fmt.Fprintln(os.Stdout)
	}
	return nil
}

func profileSet(ctx context.Context, c *cli.Command) error {
	if c.Args().Len() != 2 {
		return fmt.Errorf("usage: ant profile set <key> <value>")
	}
	key, value := c.Args().Get(0), c.Args().Get(1)
	profile, dir := activeProfile(c)
	cfg, err := loadProfileConfig(dir, profile)
	if err != nil {
		return err
	}
	// LoadConfig resolves CredentialsPath to its default for callers; don't
	// persist that resolution back into the config on round-trip.
	if cfg.AuthenticationInfo != nil && cfg.AuthenticationInfo.CredentialsPath == config.ProfileCredentialsPath(dir, profile) {
		cfg.AuthenticationInfo.CredentialsPath = ""
	}
	switch key {
	case "workspace_id":
		cfg.WorkspaceID = value
	case "base_url":
		cfg.BaseURL = value
	case "organization_id":
		cfg.OrganizationID = value
	case "scope", "client_id", "console_url":
		if cfg.AuthenticationInfo == nil || cfg.AuthenticationInfo.Type != config.AuthenticationTypeUserOAuth || cfg.AuthenticationInfo.UserOAuth == nil {
			return fmt.Errorf("%s can only be set on user_oauth profiles (run `ant auth login` first)", key)
		}
		switch key {
		case "scope":
			cfg.AuthenticationInfo.UserOAuth.Scope = value
		case "client_id":
			cfg.AuthenticationInfo.UserOAuth.ClientID = value
		case "console_url":
			cfg.AuthenticationInfo.UserOAuth.ConsoleURL = strings.TrimRight(value, "/")
		}
	default:
		return fmt.Errorf("unsupported key %q (one of: workspace_id, base_url, organization_id, scope, client_id, console_url)", key)
	}
	if err := config.SaveProfile(dir, profile, cfg); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "✓ Set %s=%s on profile %q.\n", key, value, profile)
	return nil
}
