package nuggan

import (
	"github.com/davidbyttow/govips/pkg/vips"
	"github.com/valyala/fasthttp"
	"log"
	"strings"
)

func fasthttpHandler(conf Config) func(*fasthttp.RequestCtx) {
	serve := Service(conf)

	prefix := conf.RoutePrefix + "/"

	return func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())

		if !strings.HasPrefix(path, prefix) {
			ctx.Error(
				"Route prefix expected",
				fasthttp.StatusBadRequest)

			return
		}

		// ---

		request := ImageRequest{
			Path:    path,
			Method:  string(ctx.Method()),
			Referer: fasthttpReferer(ctx),
		}

		resp := ImageResponse{
			SetStatusCode: func(code int) {
				ctx.SetStatusCode(code)
			},
			SetHeader: func(k string, v string) {
				ctx.Response.Header.Set(k, v)
			},
			Body: ctx,
		}

		serve(&request, &resp)
	}
}

func FasthttpServer(bind string, conf Config) {
	log.Printf("Starting fasthttp server on '%s' ...\n\n\tConfiguration: %v\n\n", bind, conf)

	// Setup govips
	vips.Startup(nil)

	defer vips.Shutdown()

	fasthttp.ListenAndServe(bind, fasthttpHandler(conf))
}

func fasthttpReferer(ctx *fasthttp.RequestCtx) ImageReferer {
	r := ctx.Referer()
	h := ctx.Request.Header
	userAgent := string(h.Peek("User-Agent"))

	if len(r) < 1 {
		return ImageReferer{Url: string(r), UserAgent: userAgent}
	} else {
		return ImageReferer{
			Url:       string(h.Peek("Referrer")),
			UserAgent: userAgent,
		}
	}
}
