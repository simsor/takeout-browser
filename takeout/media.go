package takeout

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

// Media is a Google Takeout picture or video
type Media struct {
	Title       string
	Description string
	Taken       int64
	Format      Format

	path string
}

// OpenMedia opens a Google Photos Takeout media file, identified by its JSON metadata file.
func OpenMedia(path string) (Media, error) {
	m := Media{}
	if !strings.HasSuffix(strings.ToLower(path), ".json") {
		return m, fmt.Errorf("media must have the JSON file extension")
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return m, fmt.Errorf("ReadFile: %v", err)
	}

	var j struct {
		Title          string `json:"title"`
		Description    string `json:"description"`
		PhotoTakenTime struct {
			Timestamp string `json:"timestamp"`
		} `json:"photoTakenTime"`
	}
	err = json.Unmarshal(data, &j)
	if err != nil {
		return m, fmt.Errorf("json.Unmarshal %s: %v", path, err)
	}

	m.Title = j.Title
	m.Description = j.Description
	m.path = path[:len(path)-len(".json")]
	format, err := m.format()
	if err != nil {
		return m, fmt.Errorf("get media format: %v", err)
	}
	m.Format = format

	taken, err := strconv.ParseInt(j.PhotoTakenTime.Timestamp, 10, 64)
	if err != nil {
		return m, fmt.Errorf("parse timestamp: %v", err)
	}
	m.Taken = taken

	return m, nil
}

func (m Media) format() (Format, error) {
	f := strings.ToLower(m.path)

	if strings.HasSuffix(f, ".jpeg") || strings.HasSuffix(f, ".jpg") {
		return FormatJPEG, nil
	}

	if strings.HasSuffix(f, ".png") {
		return FormatPNG, nil
	}

	if strings.HasSuffix(f, ".heic") {
		return FormatHEIC, nil
	}

	if strings.HasSuffix(f, ".mov") {
		return FormatMOV, nil
	}

	if strings.HasSuffix(f, ".mp4") {
		return FormatMP4, nil
	}

	if strings.HasSuffix(f, ".gif") {
		return FormatGIF, nil
	}

	return "", fmt.Errorf("unknown file format")
}

// Thumbnail returns an image.Image representing the media, of size 300x300 pixels maximum.
func (m Media) Thumbnail() (image.Image, error) {

	f, err := os.Open(m.path)
	if err != nil {
		return nil, fmt.Errorf("open: %v", err)
	}
	defer f.Close()

	if m.Format.IsVideo() {
		data, err := firstFrame(f, 300, 300)
		if err != nil {
			return nil, err
		}

		return jpeg.Decode(bytes.NewReader(data))
	}

	var img image.Image

	if m.Format == FormatPNG {
		img, err = png.Decode(f)
	} else if m.Format == FormatJPEG {
		img, err = jpeg.Decode(f)
	} else if m.Format == FormatHEIC {
		var data []byte
		data, err = convert(f, "jpeg")
		if err != nil {
			return nil, err
		}

		img, err = jpeg.Decode(bytes.NewReader(data))
	} else if m.Format == FormatGIF {
		img, err = gif.Decode(f)
	} else {
		log.Printf("No handler for this media type: %s", m.Format)
		return nil, fmt.Errorf("no handler for this media type: %s", m.Format)
	}

	if err != nil {
		return nil, fmt.Errorf("error decoding %s: %v", m.path, err)
	}

	return resize.Thumbnail(300, 300, img, resize.Bicubic), nil
}

// BrowserSafeMedia writes data to the Writer that *should* be understood by any modern browser
func (m Media) BrowserSafeMedia(w io.Writer) error {
	f, err := os.Open(m.path)
	if err != nil {
		return err
	}
	defer f.Close()

	if m.Format == FormatPNG || m.Format == FormatJPEG || m.Format == FormatMP4 || m.Format == FormatGIF {
		io.Copy(w, f)
	} else if m.Format == FormatHEIC {
		conv, err := convert(f, "jpeg")
		if err != nil {
			return err
		}
		io.Copy(w, bytes.NewReader(conv))
	} else if m.Format == FormatMOV {
		r, err := ffmpeg(f, "h264")
		if err != nil {
			return err
		}

		io.Copy(w, r)
	}

	return nil
}
