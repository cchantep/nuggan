package nuggan

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/davidbyttow/govips/pkg/vips"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestScaleDownWidthScale(t *testing.T) {
	expectedMd5 := []string{
		"47016225f65d22a44ce804c3523dfa66", // osx
		"c875a3184d68481c6bef7e765c22a61f", // linux
	}
	inputFile := "image1.jpg"

	testIt := ScaleDownTest(95, 150)
	// w_scale{92 / width(250)} < h_scale{150 / height(340)}

	checkVipsResult(t, inputFile, expectedMd5, testIt)
}

func TestScaleDownHeightScale(t *testing.T) {
	expectedMd5 := []string{
		"34d7663f3400ee1a188dca17fa1179fe", // osx
		"1ad075576225e530e30988633655296a", // linux
	}
	inputFile := "image1.jpg"

	testIt := ScaleDownTest(125, 91)
	// w_scale{125 / width(250)} < h_scale{91 / height(340)}

	checkVipsResult(t, inputFile, expectedMd5, testIt)
}

func TestScaleDownNegativeHeight(t *testing.T) {
	expectedMd5 := []string{
		"47016225f65d22a44ce804c3523dfa66", // osx
		"c875a3184d68481c6bef7e765c22a61f", // linux
	}
	inputFile := "image1.jpg"

	testIt := ScaleDownTest(95, -1)

	checkVipsResult(t, inputFile, expectedMd5, testIt)
}

// ---

func ScaleDownTest(width int, height int) func(io.Reader, io.Writer) error {
	return func(input io.Reader, output io.Writer) error {
		img, err := vips.LoadImage(input)

		if err != nil {
			return err
		}

		return ScaleDown(img, width, height, -1, output)
	}
}

// ---

func checkVipsResult(t *testing.T, inputFile string, expectedMd5 []string, testIt func(io.Reader, io.Writer) error) {
	wd, err1 := os.Getwd()

	if err1 != nil {
		t.Error(err1.Error())
	}

	input, err2 := os.Open(fmt.Sprintf("%s/../test/%s", wd, inputFile))

	if err2 != nil {
		t.Error(err2.Error())
	}

	defer input.Close()

	output, err3 := ioutil.TempFile("", "test")

	if err3 != nil {
		t.Error(err3.Error())
	}

	t.Logf("%s using temporary file %s\n", t.Name(), output.Name())

	defer output.Close()

	// Comment this one to be able to check the temporary file
	// in case of processing change
	//defer os.Remove(output.Name())

	// ---

	err5 := testIt(input, output)

	if err5 != nil {
		t.Error(err5.Error())
	}

	bytes, err6 := ioutil.ReadFile(output.Name())

	if err6 != nil {
		t.Error(err6.Error())
	}

	// ---

	md5Sum := md5.Sum(bytes)
	res := hex.EncodeToString(md5Sum[:])

	var m = -1

	for i, v := range expectedMd5 {
		if v == res {
			m = i
			break
		}
	}

	t.Logf("Output for %s: %s %d %v\n", t.Name(), res, m, expectedMd5)

	if m == -1 {
		t.Errorf("MD5(%s) != MD5(%s)", res, expectedMd5)
	}
}
