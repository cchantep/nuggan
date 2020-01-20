package nuggan

import (
	"errors"
	"fmt"
	toml "github.com/pelletier/go-toml"
	"io"
	"strings"
)

type HttpUrl = string

type Config struct {
	GroupedBaseUrls [][]HttpUrl
	RoutePrefix     string // defaulted to '/optimg' is missing
	Strict          bool
}

func (c Config) String() string {
	return fmt.Sprintf("{ GroupedBaseUrls: %v, RoutePrefix: %s, Struct: %v }", c.GroupedBaseUrls, c.RoutePrefix, c.Strict)
}

func LoadConfig(reader io.Reader) (Config, error) {
	config := Config{}
	decoder := toml.NewDecoder(reader)

	err := decoder.Decode(&config)

	if err != nil {
		return config, err
	}

	// ---

	if len(config.GroupedBaseUrls) == 0 {
		return config, errors.New("No URL group configured")
	}

	for i, g := range config.GroupedBaseUrls {
		if len(g) == 0 {
			return config, errors.New(
				fmt.Sprintf("URL group #%d is empty", i))
		}
	}

	config.RoutePrefix = strings.TrimSpace(config.RoutePrefix)

	if config.RoutePrefix == "" {
		config.RoutePrefix = "/optimg"
	} else if strings.Index(config.RoutePrefix, "/") != -1 {
		return config, errors.New(
			fmt.Sprintf("Invalid route prefix contains '/': %s",
				config.RoutePrefix))

	} else {
		config.RoutePrefix = "/" + config.RoutePrefix
	}

	return config, err
}
