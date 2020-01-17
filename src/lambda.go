package nuggan

import (
	"bytes"
	"encoding/base64"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/davidbyttow/govips/pkg/vips"
	"log"
	"strings"
)

func lambdaHandler(conf Config) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	serve := Service(conf)
	urlPrefix := conf.RoutePrefix + "/"

	return func(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		log.Printf("lambdaHandler: %s '%s'\n",
			event.HTTPMethod, event.Path)

		var output *bytes.Buffer

		headers := make(map[string]string)
		statusCode := 500

		if strings.HasPrefix(event.Path, urlPrefix) {
			request := ImageRequest{
				Path:    event.Path,
				Method:  event.HTTPMethod,
				Referer: lambdaReferer(event),
			}

			log.Printf("Image request: %v\n", request)

			output = new(bytes.Buffer)

			resp := ImageResponse{
				SetStatusCode: func(code int) {
					statusCode = code
				},
				SetHeader: func(k string, v string) {
					headers[k] = v
				},
				Body: output,
			}

			serve(&request, &resp)
		} else {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
			}, nil
		}

		content := base64.StdEncoding.EncodeToString(output.Bytes())

		return events.APIGatewayProxyResponse{
			StatusCode:      statusCode,
			Headers:         headers,
			IsBase64Encoded: true,
			Body:            content,
		}, nil
	}
}

func Lambda(conf Config) {
	log.Printf("Starting lambda ...\n\n\tConfiguration: %v\n\n", conf)

	// Setup govips
	vips.Startup(nil)

	defer vips.Shutdown()

	lambda.Start(lambdaHandler(conf))
}

func lambdaReferer(event events.APIGatewayProxyRequest) ImageReferer {
	userAgent := event.Headers["user-agent"]
	r := event.Headers["referer"]

	if r != "" {
		return ImageReferer{Url: r, UserAgent: userAgent}
	} else {
		return ImageReferer{
			Url:       event.Headers["referrer"],
			UserAgent: userAgent, // alt. name
		}
	}
}
