package nuggan

import (
	"errors"
	"fmt"
	"github.com/davidbyttow/govips/pkg/vips"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ImageReferer struct {
	Url       string
	UserAgent string
}

type ImageRequest struct {
	Path    string
	Method  string
	Referer ImageReferer
}

type ImageResponse struct {
	SetStatusCode func(int)
	SetHeader     func(string, string)
	Body          io.Writer
}

// Routes:
//
//   HEAD /:routePrefix/:cropX/:cropY/:cropWidth/:cropHeight/:resizeWidth/:resizeHeight/:compressionLevel/:base64Ref
//
//   GET  /:routePrefix/:cropX/:cropY/:cropWidth/:cropHeight/:resizeWidth/:resizeHeight/:compressionLevel/:base64Ref
//
func Service(conf Config) func(*ImageRequest, *ImageResponse) {
	decodeMediaUrl := DecodeMediaUrl(conf)

	return func(req *ImageRequest, resp *ImageResponse) {
		path := strings.Split(req.Path, "/")
		fsz := len(path)

		if fsz < 9 {
			badRequest(resp, fmt.Sprintf(
				"Unexpected request to '%s'", req.Path))
			return
		} else {
			log.Printf("INFO: Serving /%s: %s\n", path[1], path[2:])

			base64Ref := path[9]

			if !strings.HasPrefix(base64Ref, "_") && conf.Strict {
				msg := fmt.Sprintf("Base64url '%s' cannot be specified in strict mode", base64Ref)

				forbidden(resp, msg)
				return
			}

			// ---

			// crop x offset (mandatory)
			x, err := strconv.Atoi(path[2])

			if err != nil {
				msg := fmt.Sprintf(
					"Invalid crop x offset: %s",
					err.Error())

				badRequest(resp, msg)
				return
			}

			// crop y offset (mandatory)
			y, err := strconv.Atoi(path[3])

			if err != nil {
				msg := fmt.Sprintf(
					"Invalid crop y offset: %s",
					err.Error())

				badRequest(resp, msg)
				return
			}

			// crop width
			cropW := -1

			if path[4] != "-" {
				wp, err := strconv.Atoi(path[4])

				if err != nil {
					msg := fmt.Sprintf(
						"Invalid crop width: %s",
						err.Error())

					badRequest(resp, msg)
					return
				}

				cropW = wp
			}

			// crop height
			cropH := -1

			if path[5] != "-" {
				h, err := strconv.Atoi(path[5])

				if err != nil {
					msg := fmt.Sprintf(
						"Invalid crop height: %s",
						err.Error())

					badRequest(resp, msg)
					return
				}

				cropH = h
			}

			// resize width
			resizeW := -1

			if path[6] != "-" {
				vw, err := strconv.Atoi(path[6])

				if err != nil {
					msg := fmt.Sprintf(
						"Invalid resize width '%s': %v",
						path[6], err.Error())

					badRequest(resp, msg)
					return
				}

				resizeW = vw
			}

			// resize height
			resizeH := -1

			if path[7] != "-" {
				vh, err := strconv.Atoi(path[7])

				if err != nil {
					msg := fmt.Sprintf(
						"Invalid resize height '%s': %v", path[7], err.Error())

					badRequest(resp, msg)
					return
				}

				resizeH = vh
			}

			// compression level
			compressionLevel := -1

			if path[8] != "-" {
				cl, err := strconv.Atoi(path[8])

				if err != nil {
					msg := fmt.Sprintf(
						"Invalid compression level '%s': %v", path[8], err.Error())

					badRequest(resp, msg)
					return
				}

				compressionLevel = cl
			}

			// media
			mediaUrl, err := decodeMediaUrl(base64Ref)

			if err != nil {
				writeError(resp, err)
				return
			}

			log.Printf("INFO: Resolve backend URL: '%s'\n",
				mediaUrl)

			// Fetch image from public HTTP URL
			imgResp, err := http.Get(mediaUrl)

			if err != nil {
				writeError(resp, err)
				return
			}

			defer imgResp.Body.Close()

			status := imgResp.StatusCode

			if status == 404 {
				msg := fmt.Sprintf(
					"Media not found: %s", base64Ref)

				imageNotFound(
					req.Referer,
					resp,
					errors.New(msg),
					resizeW,
					resizeH)

				return
			}

			if status != 200 {
				msg := fmt.Sprintf(
					"Fails to fetch media '%s': %d",
					base64Ref, status)

				imageNotFound(
					req.Referer,
					resp,
					errors.New(msg),
					resizeW,
					resizeH)

				return
			}

			// Prepare headers
			origEtag := base64Ref

			for name, vs := range imgResp.Header {
				for _, v := range vs {
					if name == "Etag" {
						if v[0] == '"' { // unquote
							origEtag =
								v[1 : len(v)-1]
						}
					}

					if name == "Date" ||
						name == "Last-Modified" {
						resp.SetHeader(name, v)
					}
				}
			}

			etag := path[:9]
			etag[0] = origEtag

			resp.SetHeader("Etag", strings.Join(etag, "/"))

			// ---

			if req.Method == "HEAD" {
				return
			}

			// ---

			// Apply crop
			croppedImg, err := Crop(
				imgResp.Body, x, y, cropW, cropH)

			if err != nil {
				writeError(resp, err)
				return
			}

			defer croppedImg.Close()

			imgFmt := croppedImg.Format()

			resp.SetHeader(
				"Content-Type",
				fmt.Sprintf(
					"image/%s", vips.ImageTypes[imgFmt]))

			resp.SetHeader(
				"Content-Disposition",
				fmt.Sprintf("inline; filename=\"%s%s\"",
					base64Ref, imgFmt.OutputExt()))

			// Output image on response
			var rerr error = nil

			if resizeW > 0 {
				rerr = ScaleDown(
					croppedImg,
					resizeW,
					resizeH,
					compressionLevel,
					resp.Body)

			} else {
				rerr = Strip(croppedImg, resp.Body)
			}

			if rerr != nil {
				writeError(resp, rerr)
				return
			}
		}
	}
}

func imageNotFound(
	referer ImageReferer,
	resp *ImageResponse,
	err error,
	width int,
	height int) {

	resp.SetStatusCode(404)

	vw := width
	vh := height

	if vw < 1 {
		vw = 1
	}

	if vh < 1 {
		vh = 1
	}

	placeholder, err2 := vips.Gaussnoise(vw, vh)

	if err2 != nil {
		writeError(resp, err2)
		return
	}

	// ---

	msg := fmt.Sprintf("Image not found: %s", err.Error())

	log.Printf("ERROR: %s {referer: %v}\n", msg, referer)

	image := vips.NewImageRef(
		placeholder, vips.ImageTypeGIF)

	resp.SetHeader("Content-Type", "image/gif")

	resp.SetHeader("Content-Disposition",
		"inline; filename=\"not-found.gif\"")

	resp.SetHeader("Cache-Control",
		"public, no-cache, no-store, must-revalidate")

	err4 := Strip(image, resp.Body)

	if err4 != nil {
		writeError(resp, err4)
		return
	}
}

func writeError(resp *ImageResponse, err error) {
	resp.SetStatusCode(500)

	msg := err.Error()

	log.Printf("WARNING: Internal error: %s\n", msg)

	resp.SetHeader("Content-Type", "text/plain")

	fmt.Fprintf(resp.Body, msg)
}

func badRequest(resp *ImageResponse, msg string) {
	resp.SetStatusCode(400)

	log.Printf("WARNING: Bad request: %s\n", msg)

	resp.SetHeader("Content-Type", "text/plain")

	fmt.Fprintf(resp.Body, msg)
}

func forbidden(resp *ImageResponse, msg string) {
	resp.SetStatusCode(403)

	log.Printf("WARNING: Forbidden: %s\n", msg)

	resp.SetHeader("Content-Type", "text/plain")

	fmt.Fprintf(resp.Body, msg)
}
