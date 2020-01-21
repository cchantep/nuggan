package nuggan

import (
	//"os"
	"reflect"
	"strings"
	"testing"
)

func TestEmptyGroups(t *testing.T) {
	_, err := LoadConfig(strings.NewReader(`
groupedBaseUrls = []
`))

	expected := "No URL group configured"

	if err == nil || err.Error() != expected {
		t.Errorf("Expected error '%s': %v", expected, err)
	}
}

func TestEmptyGroup(t *testing.T) {
	_, err := LoadConfig(strings.NewReader(`
groupedBaseUrls = [
  [
    "https://upload.wikimedia.org/wikipedia/commons"
  ],
  []
]
`))

	expected := "URL group #1 is empty"

	if err == nil || err.Error() != expected {
		t.Errorf("Expected error '%s': %v", expected, err)
	}
}

func TestValidConfig(t *testing.T) {
	got, err := LoadConfig(strings.NewReader(`
groupedBaseUrls = [
  [
    "https://upload.wikimedia.org/wikipedia/commons"
  ],
  [
    "https://cdn0.iconfinder.com/data/icons",
    "https://cdn1.iconfinder.com/data/icons"
  ]
]
`))

	if err != nil {
		t.Error(err.Error())
	}

	// ---

	expected := Config{
		GroupedBaseUrls: [][]HttpUrl{
			{
				"https://upload.wikimedia.org/wikipedia/commons",
			},
			{
				"https://cdn0.iconfinder.com/data/icons",
				"https://cdn1.iconfinder.com/data/icons",
			},
		},
		RoutePrefix: "/optimg",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("%s != %s\n", got, expected)
	}
}

func TestRoutePrefixConfig(t *testing.T) {
	got, err := LoadConfig(strings.NewReader(`
groupedBaseUrls = [
  [
    "https://upload.wikimedia.org/wikipedia/commons"
  ]
]
routePrefix = "custom"
cacheControl = "max-age=3600, s-maxage=7200"
`))

	if err != nil {
		t.Error(err.Error())
	}

	// ---

	expected := Config{
		GroupedBaseUrls: [][]HttpUrl{
			{
				"https://upload.wikimedia.org/wikipedia/commons",
			},
		},
		RoutePrefix:  "/custom",
		CacheControl: "max-age=3600, s-maxage=7200",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("%s != %s\n", got, expected)
	}
}

func TestInvalidRoutePrefixConfig(t *testing.T) {
	_, err := LoadConfig(strings.NewReader(`
groupedBaseUrls = [
  [
    "https://upload.wikimedia.org/wikipedia/commons"
  ]
]
routePrefix = "in/alid"
`))

	expected := "Invalid route prefix contains '/': in/alid"

	if err == nil || err.Error() != expected {
		t.Errorf("Expected error for invalid route prefix: %v", err)
	}
}
