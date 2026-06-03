# VHS tape that records the demo GIF embedded in README.md.
#
# To re-record, run ./scripts/record-demo from the repository root. It builds
# the CLI, starts a local mock API server, places an `ant` wrapper pointed at
# the mock in /tmp/ant-demo, and runs vhs against this tape.

Output .github/demo.gif

Set Shell "bash"
Set FontSize 18
Set Width 1200
Set Height 680
Set Padding 24
Set WindowBar Colorful
Set TypingSpeed 30ms
Set CursorBlink false
Set Theme { "name": "Claude", "background": "#262624", "foreground": "#F0EEE6", "cursor": "#D97757", "selection": "#51504A", "black": "#262624", "red": "#E5484D", "green": "#62A87C", "yellow": "#D9A45B", "blue": "#7BA4C7", "magenta": "#C47ED1", "cyan": "#6FBFC9", "white": "#F0EEE6", "brightBlack": "#83827D", "brightRed": "#FF6369", "brightGreen": "#7DC894", "brightYellow": "#EBC06D", "brightBlue": "#96BDE0", "brightMagenta": "#D9A0E8", "brightCyan": "#8BD4DE", "brightWhite": "#FFFFFF" }

# Hidden setup: put the authenticated `ant` wrapper first on PATH.
Hide
Type `PATH="/tmp/ant-demo:$PATH"; clear`
Enter
Show

# Scene 1: send a message to Claude.
Type "ant messages create \"
Enter
Type "  --model claude-opus-4-6 \"
Enter
Type "  --max-tokens 1024 \"
Enter
Type `  --message '{role: user, content: "Hello, Claude"}'`
Sleep 500ms
Enter
Sleep 4.5s

# Scene 2: the same request, extracting just the reply text.
Type "ant messages create \"
Enter
Type "  --model claude-opus-4-6 \"
Enter
Type "  --max-tokens 1024 \"
Enter
Type `  --message '{role: user, content: "Hello, Claude"}' \`
Enter
Type "  --transform content.0.text --raw-output"
Sleep 500ms
Enter
Sleep 3.5s

# Scene 3: browse a resource in the interactive explorer.
Hide
Type "clear"
Enter
Show
Type "ant models retrieve --model-id claude-opus-4-6"
Sleep 500ms
Enter
Sleep 2.5s
Down@300ms 6
Sleep 1s
Right
Sleep 2s
Down@300ms 3
Sleep 1.5s
Type "q"
Sleep 1.5s
