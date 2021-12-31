package takeout

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"runtime"
)

func ffmpeg(in io.Reader, format string) (io.Reader, error) {
	conv := []string{"./ffmpeg"}
	if runtime.GOOS == "windows" {
		conv = []string{".\\ffmpeg.exe"}
	}

	conv = append(conv, "-i", "-", "-f", format, "-")
	cmd := exec.Command(conv[0], conv[1:]...)

	cmd.Stdin = in
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	go func() {
		err := cmd.Run()
		if err != nil {
			if exerr, ok := err.(*exec.ExitError); ok {
				log.Printf("ffmpeg convert exit error: %s", exerr.Stderr)
			} else {
				log.Printf("ffmpeg convert error: %v", err)
			}
		}
	}()

	return out, nil
}

func firstFrame(in io.Reader, width, height int) ([]byte, error) {
	conv := []string{"./ffmpeg"}
	if runtime.GOOS == "windows" {
		conv = []string{".\\ffmpeg.exe"}
	}

	conv = append(conv, "-i", "-", "-f", "image2", "-r", "1", "-s", fmt.Sprintf("%dx%d", width, height), "-frames:v", "1", "-")
	cmd := exec.Command(conv[0], conv[1:]...)

	cmd.Stdin = in
	data, err := cmd.Output()
	if err != nil {
		if exerr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("exit error: %s", exerr.Stderr)
		}

		return nil, err
	}

	return data, nil
}
