package infrastructures

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"

	"github.com/nfnt/resize"
)

func GenerateThumbnail(file io.Reader, width uint) (*bytes.Buffer, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	thumbnail := resize.Resize(width, 0, img, resize.Lanczos3)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, thumbnail, nil)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
