package nuggan

import (
	"github.com/davidbyttow/govips/pkg/vips"
	"log"
	"net/http"
)

func standaloneHandler(conf Config) func(http.ResponseWriter, *http.Request) {
	serve := Service(conf)

	return func(w http.ResponseWriter, req *http.Request) {
		request := ImageRequest{
			Path:    req.URL.Path,
			Method:  req.Method,
			Referer: httpReferer(req),
		}

		headers := w.Header()
		resp := ImageResponse{
			SetStatusCode: func(code int) {
				w.WriteHeader(code)
			},
			SetHeader: func(k string, v string) {
				headers.Set(k, v)
			},
			Body: w,
		}

		serve(&request, &resp)
	}
}

func StandaloneServer(bind string, conf Config) {
	urlPrefix := conf.RoutePrefix + "/"

	http.HandleFunc(urlPrefix, standaloneHandler(conf))

	log.Printf("Starting standalone server on '%s' ...\n\n\tConfiguration: %v\n\n", bind, conf)

	// Setup govips
	vips.Startup(nil)

	defer vips.Shutdown()

	http.ListenAndServe(bind, nil)
}

func httpReferer(req *http.Request) ImageReferer {
	r := req.Referer()
	userAgent := req.Header.Get("User-Agent")

	if r != "" {
		return ImageReferer{Url: r, UserAgent: userAgent}
	} else {
		return ImageReferer{
			Url:       req.Header.Get("Referrer"),
			UserAgent: userAgent,
		}
	}
}
