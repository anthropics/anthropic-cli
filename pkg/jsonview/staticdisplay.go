package jsonview

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/muesli/reflow/truncate"
	"github.com/tidwall/gjson"
)

const (
	tabWidth = 2
	maxDepth = 4
)

var (
	keyStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("75")).Bold(false)
	stringValueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("113"))
	numberValueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("215"))
	boolValueStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("207"))
	nullValueStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Italic(true)
	bulletStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
	containerStyle   = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("63")).
				Padding(0, 1)
)

func formatJSON(jsonStr string, width int) string {
	if !gjson.Valid(jsonStr) {
		return nullValueStyle.Render("Invalid JSON")
	}
	return formatResult(gjson.Parse(jsonStr), 0, width)
}

func formatResult(result gjson.Result, indent, width int) string {
	switch result.Type {
	case gjson.String:
		str := result.Str
		if str == "" {
			return nullValueStyle.Render("(empty)")
		}
		if lipgloss.Width(str) > width {
			str = truncate.String(str, uint(width-1)) + "…"
		}
		return stringValueStyle.Render(str)
	case gjson.Number:
		return numberValueStyle.Render(result.Raw)
	case gjson.True:
		return boolValueStyle.Render("yes")
	case gjson.False:
		return boolValueStyle.Render("no")
	case gjson.Null:
		return nullValueStyle.Render("null")
	case gjson.JSON:
		if result.IsArray() {
			return formatJSONArray(result, indent, width)
		}
		return formatJSONObject(result, indent, width)
	default:
		return stringValueStyle.Render(result.String())
	}
}

func isSingleLine(result gjson.Result, indent int) bool {
	return !(result.IsObject() || result.IsArray()) || (indent >= maxDepth)
}

func formatJSONArray(result gjson.Result, indent, width int) string {
	items := result.Array()
	if len(items) == 0 {
		return nullValueStyle.Render("(none)")
	}

	if indent >= maxDepth {
		return bulletStyle.Render(formatArray(result))
	}

	numberWidth := lipgloss.Width(fmt.Sprintf("%d. ", len(items)))

	var formattedItems []string
	for i, item := range items {
		number := fmt.Sprintf("%d.", i+1)
		numbering := getIndent(indent) + bulletStyle.Render(number)

		// If the item will be a one-liner, put it inline after the numbering,
		// otherwise it starts with a newline and goes below the numbering.
		itemWidth := width
		if isSingleLine(item, indent+1) {
			// Add right-padding:
			numbering += strings.Repeat(" ", numberWidth-lipgloss.Width(number))
			itemWidth = width - lipgloss.Width(numbering)
		}
		value := formatResult(item, indent+1, itemWidth)
		formattedItems = append(formattedItems, numbering+value)
	}
	return "\n" + strings.Join(formattedItems, "\n")
}

func formatJSONObject(result gjson.Result, indent, width int) string {
	keys := result.Get("@keys").Array()
	if len(keys) == 0 {
		return nullValueStyle.Render("(empty)")
	}

	if indent >= maxDepth {
		short := formatObject(result)
		if lipgloss.Width(short) > width {
			short = truncate.String(short, uint(width-1)) + "…"
		}
		return bulletStyle.Render(short)
	}

	var items []string
	for _, key := range keys {
		value := result.Get(key.String())
		keyStr := getIndent(indent) + keyStyle.Render(key.String()+":")
		// If item will be a one-liner, put it inline after the key, otherwise
		// it starts with a newline and goes below the key.
		itemWidth := width
		if isSingleLine(value, indent+1) {
			keyStr += " "
			itemWidth = width - lipgloss.Width(keyStr)
		}
		formattedValue := formatResult(value, indent+1, itemWidth)
		items = append(items, keyStr+formattedValue)
	}

	return "\n" + strings.Join(items, "\n")
}

func getIndent(indent int) string {
	return strings.Repeat(" ", indent*tabWidth)
}

func RenderJSON(title string, jsonStr string) string {
	width, _, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		width = 80
	}
	width -= containerStyle.GetBorderLeftSize() + containerStyle.GetBorderRightSize() +
		containerStyle.GetPaddingLeft() + containerStyle.GetPaddingRight()
	content := strings.TrimLeft(formatJSON(jsonStr, width), "\n")
	return titleStyle.Render(title) + "\n" + containerStyle.Render(content)
}

func DisplayJSON(title string, jsonStr string) {
	fmt.Println(RenderJSON(title, jsonStr))
}
