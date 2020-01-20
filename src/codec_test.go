package nuggan

import (
	"testing"
)

// ---

var config1 = Config{
	GroupedBaseUrls: [][]HttpUrl{
		{
			"https://upload.wikimedia.org/wikipedia/commons",
		},
		{
			"https://cdn0.iconfinder.com/data/icons",
			"https://cdn1.iconfinder.com/data/icons",
		},
	},
}

// --- Encode tests

var encode1 = EncodeMediaUrl(config1)

func TestEncodeMediaUrlInGroup1(t *testing.T) {
	got := encode1("https://upload.wikimedia.org/wikipedia/commons/thumb/a/a0/Ansberg_Veitskapelle_1060027-PSD.jpg/1000px-Ansberg_Veitskapelle_1060027-PSD.jpg")
	expected := "_0_L3RodW1iL2EvYTAvQW5zYmVyZ19WZWl0c2thcGVsbGVfMTA2MDAyNy1QU0QuanBnLzEwMDBweC1BbnNiZXJnX1ZlaXRza2FwZWxsZV8xMDYwMDI3LVBTRC5qcGc="

	if got != expected {
		t.Errorf("%s != %s\n", got, expected)
	}
}

func TestEncodeMediaUrlInGroup2(t *testing.T) {
	expected := "_1_L29jdGljb25zLzEwMjQvbWFyay1naXRodWItNTEyLnBuZw=="

	got1 := encode1("https://cdn0.iconfinder.com/data/icons/octicons/1024/mark-github-512.png")

	if got1 != expected {
		t.Errorf("%s != %s\n", got1, expected)
	}

	got2 := encode1("https://cdn1.iconfinder.com/data/icons/octicons/1024/mark-github-512.png")

	if got2 != expected {
		t.Errorf("%s != %s\n", got2, expected)
	}

}

func TestEncodeMediaUrlNoGroup(t *testing.T) {
	got := encode1("https://blog.golang.org/lib/godoc/images/go-logo-blue.svg")
	expected := "aHR0cHM6Ly9ibG9nLmdvbGFuZy5vcmcvbGliL2dvZG9jL2ltYWdlcy9nby1sb2dvLWJsdWUuc3Zn"

	if got != expected {
		t.Errorf("%s != %s\n", got, expected)
	}
}

// --- Decode tests

var decode1 = DecodeMediaUrl(config1)

func TestDecodeMediaUrlInGroup1(t *testing.T) {
	expected := "https://upload.wikimedia.org/wikipedia/commons/thumb/a/a0/Ansberg_Veitskapelle_1060027-PSD.jpg/1000px-Ansberg_Veitskapelle_1060027-PSD.jpg"

	got, err := decode1("_0_L3RodW1iL2EvYTAvQW5zYmVyZ19WZWl0c2thcGVsbGVfMTA2MDAyNy1QU0QuanBnLzEwMDBweC1BbnNiZXJnX1ZlaXRza2FwZWxsZV8xMDYwMDI3LVBTRC5qcGc=")

	if err != nil {
		t.Error(err.Error())
	}

	if got != expected {
		t.Errorf("%s != %s\n", got, expected)
	}
}

func TestDecodeMediaUrlInGroup2(t *testing.T) {
	expected := "https://cdn0.iconfinder.com/data/icons/octicons/1024/mark-github-512.png"

	got, err := decode1("_1_L29jdGljb25zLzEwMjQvbWFyay1naXRodWItNTEyLnBuZw==")

	if err != nil {
		t.Error(err.Error())
	}

	if got != expected {
		t.Errorf("%s != %s\n", got, expected)
	}
}

func TestDecodeMediaUrlNoGroup(t *testing.T) {
	expected := "https://blog.golang.org/lib/godoc/images/go-logo-blue.svg"

	got, err := decode1("aHR0cHM6Ly9ibG9nLmdvbGFuZy5vcmcvbGliL2dvZG9jL2ltYWdlcy9nby1sb2dvLWJsdWUuc3Zn")

	if err != nil {
		t.Error(err.Error())
	}

	if got != expected {
		t.Errorf("%s != %s\n", got, expected)
	}
}

func TestDecodeMediaUrlMissingSeparator(t *testing.T) {
	_, err := decode1("_1L29jdGljb25zLzEwMjQvbWFyay1naXRodWItNTEyLnBuZw==")

	expected := "Invalid base64Ref '_1L29jdGljb25zLzEwMjQvbWFyay1naXRodWItNTEyLnBuZw==': second '_' separator expected after group index"

	if err == nil || err.Error() != expected {
		t.Error("Error must be raised for missing '_' separator")
	}
}

func TestDecodeMediaUrlInvalidGroupIndex(t *testing.T) {
	_, err := decode1("_10_L29jdGljb25zLzEwMjQvbWFyay1naXRodWItNTEyLnBuZw==")

	expected := "Invalid group index: 10"

	if err == nil || err.Error() != expected {
		t.Error("Error must be raised for invalid group index")
	}
}
