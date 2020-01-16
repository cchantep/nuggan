package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/davidbyttow/govips/pkg/vips"
	"io"
	"log"
	"net/http"
	"nuggan"
	"os"
	"strings"
)

var (
	inputUrl = flag.String("in", "", "Url to load")
	output   = flag.String("out", "output", "File to write out")

	width  = flag.Int("w", 640, "Viewport width")
	height = flag.Int("h", -1, "Viewport height")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nScale down an image locally:\n\n\t%s -in 'http://input/image/url' -out '/path/for/output/image' -w scale_down_width_int -h scale_down_height_int\n", os.Args[0])

		fmt.Fprintf(os.Stderr, "\nDetailed options:\n\n")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "")
	}

	flag.Parse()

	if *inputUrl == "" {
		log.Fatalf("No input.\r\n\r\n")
		flag.Usage()
		os.Exit(1)
	}

	// ---

	err := cliScaleDown(*inputUrl, *output, *width, *height)

	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(2)
	}

	log.Printf("%s scaled down at %d x %d to %s\n",
		*inputUrl, *width, *height, *output)

}

func cliScaleDown(
	inputUrl string,
	output string,
	width int,
	height int,
) error {
	var reader io.Reader = nil

	if !strings.HasPrefix(inputUrl, "file://") {
		// Fetch image from public HTTP URL
		imgResp, err := http.Get(inputUrl)

		if err != nil {
			return errors.New(
				fmt.Sprintf(
					"Fails to fetch image from %s: %v",
					inputUrl, err))

		}

		defer imgResp.Body.Close()

		if imgResp.StatusCode != 200 {
			return errors.New(
				fmt.Sprintf(
					"Fails to fetch image from %s: %s",
					inputUrl, imgResp.Status))

		}

		reader = imgResp.Body
	} else {
		file, err := os.Open((inputUrl)[7:])

		if err != nil {
			return errors.New(
				fmt.Sprintf(
					"Fails to read image from %s: %v",
					inputUrl, err))

		}

		defer file.Close()

		reader = file
	}

	// ---

	writer, err := os.Create(output)

	if err != nil {
		return errors.New(fmt.Sprintf("Cannot create writer", err))
	}

	// ---

	defer writer.Close()

	vips.Startup(nil)

	defer vips.Shutdown()

	// ---

	image, err := vips.LoadImage(reader)

	if err != nil {
		return errors.New(fmt.Sprintf("Cannot load image", err))
	}

	// ---

	defer image.Close()

	nuggan.ScaleDown(image, width, height, -1, writer)

	vips.Shutdown()

	//vips.PrintObjectReport("nuggan")

	return nil
}
