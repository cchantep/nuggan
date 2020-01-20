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
	standaloneServer = flag.String("server", "",
		"Bind address (e.g. ':8080') to run a standalone net/http server")

	fasthttpServer = flag.String("fast-http", "",
		"Bind address (e.g. ':8080') to run a standalone fasthttp server")

	lambda = flag.Bool("lambda", false, "Starts the server as an AWS lambda")

	serverConfig = flag.String("server-config", "server.conf",
		"If running as standalone server or lambda, path to configuration")

	encodeUrl = flag.String("encode-url", "", "An image URL to be encoded according the 'groupedBaseUrls' setting in the server configuration (e.g. http://image/url/to/be/encoded/according/server-conf)")

	inputUrl = flag.String("in", "", "Url to load")
	output   = flag.String("out", "output", "File to write out")

	width  = flag.Int("w", 640, "Viewport width")
	height = flag.Int("h", -1, "Viewport height")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nRun as standalone net/http server:\n\n\t%s -server ':8080' -server-config server.conf\n", os.Args[0])

		fmt.Fprintf(os.Stderr, "\nRun as standaline fasthttp server:\n\n\t%s -fast-http ':8080' -server-config server.conf\n", os.Args[0])

		fmt.Fprintf(os.Stderr, "\nRun as AWS lambda function:\n\n\t%s -lambda -server-config server.conf\n", os.Args[0])

		fmt.Fprintf(os.Stderr, "\nEncode an URL according a server configuration:\n\n\t%s -server-config server.conf -encode-url 'http://an/image/url'\n", os.Args[0])

		fmt.Fprintf(os.Stderr, "\nScale down an image locally:\n\n\t%s -in 'http://input/image/url' -out '/path/for/output/image' -w scale_down_width_int -h scale_down_height_int\n", os.Args[0])

		fmt.Fprintf(os.Stderr, "\nDetailed options:\n\n")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "")
	}

	flag.Parse()

	standaloneBind := *standaloneServer
	fasthttpBind := *fasthttpServer
	isLambda := *lambda

	if standaloneBind != "" || fasthttpBind != "" || isLambda {
		// Start as standalone server or lambda

		f, err := os.Open(*serverConfig)

		if err != nil {
			fmt.Fprintf(os.Stderr,
				"Fails to open configuration file: %s\n",
				err.Error())

			flag.Usage()

			return
		}

		conf, err := nuggan.LoadConfig(f)

		if err != nil {
			fmt.Fprintf(os.Stderr,
				"Fails to load configuration: %s\n",
				err.Error())

			flag.Usage()

			return
		}

		if standaloneBind != "" {
			nuggan.StandaloneServer(standaloneBind, conf)
		} else if fasthttpBind != "" {
			nuggan.FasthttpServer(fasthttpBind, conf)
		} else {
			nuggan.Lambda(conf)
		}

		return
	}

	// ---

	if *encodeUrl != "" {
		f, err := os.Open(*serverConfig)

		if err != nil {
			fmt.Fprintf(os.Stderr,
				"Fails to open configuration file: %s\n",
				err.Error())

			flag.Usage()

			return
		}

		conf, err := nuggan.LoadConfig(f)

		if err != nil {
			fmt.Fprintf(os.Stderr,
				"Fails to load configuration: %s\n",
				err.Error())

			flag.Usage()

			return
		}

		// ---

		encode := nuggan.EncodeMediaUrl(conf)
		repr := encode(*encodeUrl)

		log.Printf("\nEncode '%s':\n\n\t%s\n\n\te.g. http://localhost:8080/%s/0/0/-/-/-/-/-/%s\n\n", *encodeUrl, repr, conf.RoutePrefix, repr)

		return
	}

	// ---

	if *inputUrl == "" {
		log.Printf("ERROR: No input.\r\n\r\n")
		flag.Usage()
		os.Exit(1)
	}

	// ---

	err := cliScaleDown(*inputUrl, *output, *width, *height)

	if err != nil {
		log.Printf(err.Error())
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
