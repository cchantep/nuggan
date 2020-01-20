package nuggan

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type PreparedBase struct {
	Url  string
	Size int
}

/**
 * Returns a function that encodes a given media `url`,
 * using a list of well-known base URLs and base64).
 */
func EncodeMediaUrl(config Config) func(string) string {
	prepared := make([][]PreparedBase, len(config.GroupedBaseUrls))

	for i, g := range config.GroupedBaseUrls {
		preparedGroup := make([]PreparedBase, len(g))

		for j, b := range g {
			preparedGroup[j] = PreparedBase{
				Url:  b,
				Size: len(b),
			}
		}

		prepared[i] = preparedGroup
	}

	return func(url string) string {
		var groupId = -1
		var reqPath = url

		for i, g := range prepared {
			for _, base := range g {
				if strings.HasPrefix(url, base.Url) {
					groupId = i
					reqPath = url[base.Size:]
				}
			}
		}

		if groupId == -1 {
			// Won't start with '_' as not a base64 character
			// (but possible in URL)
			return base64Enc(url)
		} else {
			return fmt.Sprintf("_%d_%s", groupId, base64Enc(reqPath))
		}
	}
}

/**
 * Returns a function that decodes a media URL
 * (from a string `repr`esentation previously produced by `encodeMediaUrl`).
 */
func DecodeMediaUrl(config Config) func(string) (string, error) {
	groupLen := len(config.GroupedBaseUrls)
	prepared := make([]string, groupLen)

	for i, g := range config.GroupedBaseUrls {
		prepared[i] = g[0]
	}

	return func(repr string) (string, error) {
		prefix := -1
		unprefixed := ""

		if repr[0] == '_' {
			idx := strings.Index(repr[1:], "_")

			if idx > 0 {
				unprefixed = repr[idx+2:]
				p := repr[1 : idx+1]

				px, err := strconv.Atoi(p)

				if err != nil {
					return "", err
				}

				prefix = px
			}

			if prefix == -1 {
				return "", errors.New(fmt.Sprintf("Invalid base64Ref '%s': second '_' separator expected after group index", repr))
			}
		}

		if prefix != -1 {
			decoded, err := base64Dec(unprefixed)

			if err != nil {
				return "", err
			} else if prefix < 0 || prefix >= groupLen {
				return "", errors.New(fmt.Sprintf(
					"Invalid group index: %d", prefix))
			} else {
				return (prepared[prefix] + decoded), nil
			}
		} else {
			return base64Dec(repr)
		}
	}
}

// ---

func base64Enc(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func base64Dec(input string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(input)

	if err != nil {
		return "", err
	}

	return string(data), nil
}
