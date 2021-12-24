package takeout

import (
	"io"
	"os/exec"
	"runtime"
)

func convert(in io.Reader, format string) ([]byte, error) {
	conv := []string{"convert"}
	if runtime.GOOS == "windows" {
		conv = []string{"magick.exe", "convert"}
	}

	conv = append(conv, "-", format+":-")
	cmd := exec.Command(conv[0], conv[1:]...)

	cmd.Stdin = in

	return cmd.Output()
}
