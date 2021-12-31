package takeout

// Format represents a media format
type Format string

var (
	FormatPNG  Format = "PNG"
	FormatJPEG Format = "JPEG"
	FormatHEIC Format = "HEIC"
	FormatGIF  Format = "GIF"

	FormatMOV Format = "MOV"
	FormatMP4 Format = "MP4"
)

func (f Format) IsPicture() bool {
	return f == FormatPNG || f == FormatHEIC || f == FormatJPEG || f == FormatGIF
}

func (f Format) IsVideo() bool {
	return f == FormatMOV || f == FormatMP4
}
