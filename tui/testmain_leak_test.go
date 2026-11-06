// SPDX-License-Identifier: MIT

package tui

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/MichaelKinsy/PiG/internal/testenv"
)

// hostTerminalEnv names the variables through which the developer's terminal, multiplexer or SSH session changes tui behavior
// (scroll regions, mouse mode, escape timeout, image and color capabilities). Tests that need one set it with t.Setenv.
var hostTerminalEnv = []string{
	"ZELLIJ", "TMUX", "STY", "SSH_CONNECTION", "SSH_CLIENT", "SSH_TTY",
	"TERM_PROGRAM", "TERMINAL_EMULATOR", "COLORTERM", "COLORFGBG", "WT_SESSION",
	"KITTY_WINDOW_ID", "GHOSTTY_RESOURCES_DIR", "WEZTERM_PANE", "ITERM_SESSION_ID",
	"PI_HYPERLINKS", "PI_IMAGE_PROTOCOL", "PI_TRUE_COLOR", "PI_PROGRAM_STATUS", "PIG_RENDER_DEBUG",
}

var hostTerminalEnvPrefixes = []string{"ZELLIJ_", "TMUX_", "WARP_", "HERDR_", "ORCA_", "PI_TUI_", "PIG_TUI_"}

// TestMain runs the tests without the host terminal's environment and fails the run when any test leaks a goroutine.
func TestMain(m *testing.M) {
	for _, entry := range os.Environ() {
		name, _, _ := strings.Cut(entry, "=")
		hasPrefix := func(prefix string) bool { return strings.HasPrefix(name, prefix) }
		if slices.Contains(hostTerminalEnv, name) || slices.ContainsFunc(hostTerminalEnvPrefixes, hasPrefix) {
			_ = os.Unsetenv(name)
		}
	}
	os.Exit(testenv.RunWithoutGoroutineLeaks(m))
}
