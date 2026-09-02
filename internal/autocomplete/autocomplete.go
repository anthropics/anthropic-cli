package autocomplete

import (
	"context"
	"embed"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/urfave/cli/v3"
)

type CompletionStyle string

const (
	CompletionStyleZsh        CompletionStyle = "zsh"
	CompletionStyleBash       CompletionStyle = "bash"
	CompletionStylePowershell CompletionStyle = "pwsh"
	CompletionStyleFish       CompletionStyle = "fish"
)

//go:embed shellscripts
var autoCompleteFS embed.FS

// shellScriptFiles maps each supported shell to its embedded script template.
// Templates contain the literal token `__APPNAME__`, replaced at render time.
var shellScriptFiles = map[CompletionStyle]string{
	CompletionStyleBash:       "shellscripts/bash_autocomplete.bash",
	CompletionStyleZsh:        "shellscripts/zsh_autocomplete.zsh",
	CompletionStyleFish:       "shellscripts/fish_autocomplete.fish",
	CompletionStylePowershell: "shellscripts/pwsh_autocomplete.ps1",
}

// SupportedShells returns the set of shells for which a completion script can
// be rendered, in a deterministic order suitable for help text and errors.
func SupportedShells() []CompletionStyle {
	return []CompletionStyle{
		CompletionStyleBash,
		CompletionStyleZsh,
		CompletionStyleFish,
		CompletionStylePowershell,
	}
}

// RenderCompletionScript returns the shell completion script for the given
// shell, with the template's __APPNAME__ token replaced by appName. Returns
// an error for unsupported shells or embed read failures (the latter would
// indicate a build problem, not user error).
func RenderCompletionScript(shell CompletionStyle, appName string) (string, error) {
	path, ok := shellScriptFiles[shell]
	if !ok {
		return "", unsupportedShellErr(string(shell))
	}
	b, err := autoCompleteFS.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(b), "__APPNAME__", appName), nil
}

// BuildCompletionCommand returns a `completion <shell>` cli.Command tree for
// attachment under any urfave/cli root. Leaves are generated for every shell
// in SupportedShells; the parent's Action handles `<app> completion` (no
// shell) and `<app> completion <unknown>` with a unified usage error.
//
// appName is interpolated into the install-hints rendered in the parent's
// Description.
func BuildCompletionCommand(appName string) *cli.Command {
	shells := SupportedShells()
	leaves := make([]*cli.Command, 0, len(shells))
	for _, shell := range shells {
		leaves = append(leaves, &cli.Command{
			Name:            string(shell),
			Usage:           fmt.Sprintf("Output a %s completion script for %s", shell, appName),
			HideHelpCommand: true,
			Action:          renderShellAction(shell),
		})
	}
	return &cli.Command{
		Name:            "completion",
		Usage:           "Generate shell completion scripts",
		Description:     completionDescription(appName, shells),
		Suggest:         true,
		HideHelpCommand: true,
		Commands:        leaves,
		Action: func(ctx context.Context, c *cli.Command) error {
			return cli.Exit(usageError(appName, c.Args().First()), 1)
		},
	}
}

func renderShellAction(shell CompletionStyle) cli.ActionFunc {
	return func(ctx context.Context, c *cli.Command) error {
		script, err := RenderCompletionScript(shell, c.Root().Name)
		if err != nil {
			return cli.Exit(err, 1)
		}
		if _, err := c.Writer.Write([]byte(script)); err != nil {
			return cli.Exit(err, 1)
		}
		return nil
	}
}

func shellList(shells []CompletionStyle) string {
	parts := make([]string, len(shells))
	for i, s := range shells {
		parts[i] = string(s)
	}
	return strings.Join(parts, ", ")
}

func unsupportedShellErr(shell string) error {
	return fmt.Errorf("unsupported shell %q (supported: %s)", shell, shellList(SupportedShells()))
}

// usageError returns the message for `<app> completion` invoked with no shell
// or with an unrecognized one. The two branches share the canonical shell
// list so the wording cannot drift.
func usageError(appName, shell string) string {
	if shell == "" {
		return fmt.Sprintf("%s completion <shell> — specify one of: %s", appName, shellList(SupportedShells()))
	}
	return unsupportedShellErr(shell).Error()
}

func completionDescription(appName string, shells []CompletionStyle) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Print a shell completion script for %s to stdout.\n\n", appName)
	b.WriteString("Install it by writing the output to the location your shell loads completions\nfrom:\n\n")
	for _, s := range shells {
		switch s {
		case CompletionStyleBash:
			fmt.Fprintf(&b, "  # bash (Linux)\n  %s completion bash > ~/.local/share/bash-completion/completions/%s\n\n", appName, appName)
			fmt.Fprintf(&b, "  # bash (macOS, with Homebrew bash-completion@2)\n  %s completion bash > \"$(brew --prefix)/etc/bash_completion.d/%s\"\n\n", appName, appName)
		case CompletionStyleZsh:
			fmt.Fprintf(&b, "  # zsh — pick any directory that's already on $fpath\n  %s completion zsh > \"${fpath[1]}/_%s\"\n\n", appName, appName)
		case CompletionStyleFish:
			fmt.Fprintf(&b, "  # fish\n  %s completion fish > ~/.config/fish/completions/%s.fish\n\n", appName, appName)
		case CompletionStylePowershell:
			fmt.Fprintf(&b, "  # PowerShell\n  %s completion pwsh >> \"$PROFILE.CurrentUserAllHosts\"\n\n", appName)
		}
	}
	b.WriteString("Or source it directly in your shell rc for a one-shot install:\n\n")
	fmt.Fprintf(&b, "  echo 'source <(%s completion bash)' >> ~/.bashrc\n", appName)
	fmt.Fprintf(&b, "  echo 'source <(%s completion zsh)'  >> ~/.zshrc", appName)
	return b.String()
}

type ShellCompletion struct {
	Name  string
	Usage string
}

func NewShellCompletion(name string, usage string) ShellCompletion {
	return ShellCompletion{Name: name, Usage: usage}
}

type ShellCompletionBehavior int

const (
	ShellCompletionBehaviorDefault ShellCompletionBehavior = iota
	ShellCompletionBehaviorFile                            = 10
	ShellCompletionBehaviorNoComplete
)

type CompletionResult struct {
	Completions []ShellCompletion
	Behavior    ShellCompletionBehavior
}

func isFlag(arg string) bool {
	return strings.HasPrefix(arg, "-")
}

func findFlag(cmd *cli.Command, arg string) *cli.Flag {
	name := strings.TrimLeft(arg, "-")
	for _, flag := range cmd.Flags {
		if vf, ok := flag.(cli.VisibleFlag); ok && !vf.IsVisible() {
			continue
		}

		if slices.Contains(flag.Names(), name) {
			return &flag
		}
	}
	return nil
}

func findChild(cmd *cli.Command, name string) *cli.Command {
	for _, c := range cmd.Commands {
		if !c.Hidden && c.Name == name {
			return c
		}
	}
	return nil
}

type shellCompletionBuilder struct {
	completionStyle CompletionStyle
}

func (scb *shellCompletionBuilder) createFromCommand(input string, command *cli.Command, result []ShellCompletion) []ShellCompletion {
	matchingNames := make([]string, 0, len(command.Names()))

	for _, name := range command.Names() {
		if strings.HasPrefix(name, input) {
			matchingNames = append(matchingNames, name)
		}
	}

	if scb.completionStyle == CompletionStyleBash {
		index := strings.LastIndex(input, ":") + 1
		if index > 0 {
			for _, name := range matchingNames {
				result = append(result, NewShellCompletion(name[index:], command.Usage))
			}
			return result
		}
	}

	for _, name := range matchingNames {
		result = append(result, NewShellCompletion(name, command.Usage))
	}
	return result
}

func (scb *shellCompletionBuilder) createFromFlag(input string, flag *cli.Flag, result []ShellCompletion) []ShellCompletion {
	matchingNames := make([]string, 0, len((*flag).Names()))

	for _, name := range (*flag).Names() {
		withPrefix := ""
		if len(name) == 1 {
			withPrefix = "-" + name
		} else {
			withPrefix = "--" + name
		}

		if strings.HasPrefix(withPrefix, input) {
			matchingNames = append(matchingNames, withPrefix)
		}
	}

	usage := ""
	if dgf, ok := (*flag).(cli.DocGenerationFlag); ok {
		usage = dgf.GetUsage()
	}

	for _, name := range matchingNames {
		result = append(result, NewShellCompletion(name, usage))
	}

	return result
}

func GetCompletions(completionStyle CompletionStyle, root *cli.Command, args []string) CompletionResult {
	result := getAllPossibleCompletions(completionStyle, root, args)

	// If the user has not put in a colon, filter out colon commands
	if len(args) > 0 && !strings.Contains(args[len(args)-1], ":") {
		// Nothing with anything after a colon. Create a single entry for groups with the same colon subset
		foundNames := make([]string, 0, len(result.Completions))
		filteredCompletions := make([]ShellCompletion, 0, len(result.Completions))

		for _, completion := range result.Completions {
			name := completion.Name
			firstColonIndex := strings.Index(name, ":")
			if firstColonIndex > -1 {
				name = name[0:firstColonIndex]
				completion.Name = name
				completion.Usage = ""
			}

			if !slices.Contains(foundNames, name) {
				foundNames = append(foundNames, name)
				filteredCompletions = append(filteredCompletions, completion)
			}
		}

		result.Completions = filteredCompletions
	}

	return result
}

func getAllPossibleCompletions(completionStyle CompletionStyle, root *cli.Command, args []string) CompletionResult {
	builder := shellCompletionBuilder{completionStyle: completionStyle}
	completions := make([]ShellCompletion, 0)
	if len(args) == 0 {
		for _, child := range root.Commands {
			completions = builder.createFromCommand("", child, completions)
		}
		return CompletionResult{Completions: completions, Behavior: ShellCompletionBehaviorDefault}
	}

	current := args[len(args)-1]
	preceding := args[0 : len(args)-1]
	cmd := root
	i := 0
	for i < len(preceding) {
		arg := preceding[i]

		if isFlag(arg) {
			flag := findFlag(cmd, arg)
			if flag == nil {
				i++
			} else if docFlag, ok := (*flag).(cli.DocGenerationFlag); ok && docFlag.TakesValue() {
				// All flags except for bool flags take values
				i += 2
			} else {
				i++
			}
		} else {
			child := findChild(cmd, arg)
			if child != nil {
				cmd = child
			}
			i++
		}
	}

	// Check if the previous arg was a flag expecting a value
	if len(preceding) > 0 {
		prev := preceding[len(preceding)-1]
		if isFlag(prev) {
			flag := findFlag(cmd, prev)
			if flag != nil {
				if fb, ok := (*flag).(*cli.StringFlag); ok && fb.TakesFile {
					return CompletionResult{Completions: completions, Behavior: ShellCompletionBehaviorFile}
				} else if docFlag, ok := (*flag).(cli.DocGenerationFlag); ok && docFlag.TakesValue() {
					return CompletionResult{Completions: completions, Behavior: ShellCompletionBehaviorNoComplete}
				}
			}
		}
	}

	// Completing a flag name
	if isFlag(current) {
		for _, flag := range cmd.Flags {
			completions = builder.createFromFlag(current, &flag, completions)
		}
	}

	for _, child := range cmd.Commands {
		if !child.Hidden {
			completions = builder.createFromCommand(current, child, completions)
		}
	}

	return CompletionResult{
		Completions: completions,
		Behavior:    ShellCompletionBehaviorDefault,
	}
}

func ExecuteShellCompletion(ctx context.Context, cmd *cli.Command) error {
	root := cmd.Root()
	args := rebuildColonSeparatedArgs(root.Args().Slice()[1:])

	var completionStyle CompletionStyle
	if style, ok := os.LookupEnv("COMPLETION_STYLE"); ok {
		switch style {
		case "bash":
			completionStyle = CompletionStyleBash
		case "zsh":
			completionStyle = CompletionStyleZsh
		case "pwsh":
			completionStyle = CompletionStylePowershell
		case "fish":
			completionStyle = CompletionStyleFish
		default:
			return cli.Exit("COMPLETION_STYLE must be set to 'bash', 'zsh', 'pwsh', or 'fish'", 1)
		}
	} else {
		return cli.Exit("COMPLETION_STYLE must be set to 'bash', 'zsh', 'pwsh', 'fish'", 1)
	}

	result := GetCompletions(completionStyle, root, args)

	for _, completion := range result.Completions {
		name := completion.Name
		if completionStyle == CompletionStyleZsh {
			name = strings.ReplaceAll(name, ":", "\\:")
		}
		if completionStyle == CompletionStyleZsh && len(completion.Usage) > 0 {
			_, _ = fmt.Fprintf(cmd.Writer, "%s:%s\n", name, completion.Usage)
		} else if completionStyle == CompletionStyleFish && len(completion.Usage) > 0 {
			_, _ = fmt.Fprintf(cmd.Writer, "%s\t%s\n", name, completion.Usage)
		} else {
			_, _ = fmt.Fprintf(cmd.Writer, "%s\n", name)
		}
	}
	return cli.Exit("", int(result.Behavior))
}

// When CLI arguments are passed in, they are separated on word barriers.
// Most commonly this is whitespace but in some cases that may also be colons.
// We wish to allow arguments with colons. To handle this, we append/prepend colons to their neighboring
// arguments.
//
// Example: `rebuildColonSeparatedArgs(["a", "b", ":", "c", "d"])` => `["a", "b:c", "d"]`
func rebuildColonSeparatedArgs(args []string) []string {
	if len(args) == 0 {
		return args
	}

	result := []string{}
	i := 0

	for i < len(args) {
		current := args[i]

		// Keep joining while the next element is ":" or the current element ends with ":"
		for i+1 < len(args) && (args[i+1] == ":" || strings.HasSuffix(current, ":")) {
			if args[i+1] == ":" {
				current += ":"
				i++
				// Check if there's a following element after the ":"
				if i+1 < len(args) && args[i+1] != ":" {
					current += args[i+1]
					i++
				}
			} else {
				break
			}
		}

		result = append(result, current)
		i++
	}

	return result
}
