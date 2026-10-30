package web

import (
	"fmt"
	"os/exec"
	"runtime"
)

// DefaultOpenBrowser opens url in the user's default browser using the
// platform's own opener, run directly with no shell in between: xdg-open
// on Linux, open on macOS, and rundll32's URL handler on Windows (which,
// unlike cmd.exe's "start" built-in, needs no shell). url is always our
// own locally-constructed "http://host:port/" string and is passed as a
// separate argv entry.
func DefaultOpenBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("opening browser: %w", err)
	}
	go func() { _ = cmd.Wait() }() // reap without blocking the server
	return nil
}
