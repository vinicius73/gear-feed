package browser

import (
	"errors"
	"os/exec"
	"runtime"
)

var ErrUnsupportedPlatform = errors.New("unsupported platform")

func OpenURL(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return ErrUnsupportedPlatform
	}
}
