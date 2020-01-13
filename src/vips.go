package nuggan

import (
	"github.com/davidbyttow/govips/pkg/vips"
	"io"
	"log"
)

// Reads an image from the input, and then crops it using the given parameters.
//
// - input: Image reader
// - x: Crop origin X (>= 0)
// - y: Crop origin Y (>= 0)
// - width: Crop width (or -1 if none)
// - height: Crop height (or -1 if none)
//
func Crop(
	input io.Reader,
	x int,
	y int,
	width int,
	height int) (*vips.ImageRef, error) {

	image, err1 := vips.LoadImage(input)

	if err1 != nil {
		return nil, err1
	}

	// ---

	origWidth := image.Width()
	origHeight := image.Height()

	// (0 < nx < origWidth) && (0 < ny < origHeight)
	nx := 0

	if x < 0 || x >= origWidth {
		log.Printf("WARN: Crop x %d defaulted to %d: expected > 0 and < %d\n", x, nx, origWidth)
	} else {
		nx = x
	}

	ny := 0

	if y < 0 || y >= origHeight {
		log.Printf("WARN: Crop y %d defaulted to %d: expected > 0 and < %d\n", y, ny, origHeight)
	} else {
		ny = y
	}

	nw := origWidth - nx

	if width < 0 || ((nx + width) > origWidth) {
		log.Printf("WARN: Crop width %d defaulted to %d: expected > 0 and (%d + %d) <= %d\n", width, nw, nx, width, origWidth)
	} else {
		nw = width
	}

	nh := origHeight - ny

	if height < 0 || ((ny + height) > origHeight) {
		log.Printf("WARN: Crop height %d defaulted to %d: expected > 0 and (%d + %d) <= %d\n", height, nh, ny, height, origHeight)
	} else {
		nh = height
	}

	// ---

	cropped, err2 := vips.ExtractArea(image.Image(), nx, ny, nw, nh)

	if err2 != nil {
		return nil, err2
	}

	// Reset ImageRef with cropped underlying image
	image.SetImage(cropped)

	return image, nil
}

// Scale down the given image (to a smaller size),
// and write the result to the given writer.
//
// - image: In-memory image reference
// - width: Resize width; Ignored if > image width.
// - height: Resize height; Ignored if < 0 or > image height.
// - compression: Compression level (>= 0);  Ignored if < 0.
// - output: Result writer
//
func ScaleDown(
	image *vips.ImageRef,
	width int,
	height int,
	compression int,
	output io.Writer) error {

	rw := float64(width)
	rh := float64(height)

	imgh := image.Height()
	imgw := image.Width()

	ih := float64(imgh)
	iw := float64(imgw)

	imgTx := vips.NewTransform().Image(image)

	var scale float64 = 1

	if (rh < 0 || rh <= ih) && rw <= iw {
		ws := rw / iw

		var hs float64 = 1

		if rh > 0 {
			hs = rh / ih
		}

		if ws < hs {
			scale = ws
		} else {
			scale = hs
		}
	} else {
		log.Printf("WARN: Scale defaults to %f: expected width(%f < %f) and height(%f < 0 or < %f)\n", scale, rh, ih, rw, iw)
	}

	var finalTx *vips.Transform = nil

	if compression > 0 {
		finalTx = imgTx.Compression(compression)
	} else {
		finalTx = imgTx
	}

	_, _, err3 := finalTx.
		Scale(scale).
		StripMetadata().
		Output(output).
		Apply()

	return err3
}

// Only strips image (no other transformation).
func Strip(image *vips.ImageRef, output io.Writer) error {
	_, _, err := vips.NewTransform().
		Image(image).
		StripMetadata().
		Output(output).
		Apply()

	return err
}
