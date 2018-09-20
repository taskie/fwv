package fwv

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestCSV2FWV01(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(csv01))
	w := bytes.NewBufferString("")
	app := Application{
		Mode:                    "c2f",
		UseWidth:                false,
		EastAsianAmbiguousWidth: 1,
		Whitespaces:             " ",
	}
	app.Run(r, w)
	assertEqualForEachLine(t, w.String(), fwv01)
}

func TestCSV2FWVUseWidth01(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(csv01))
	w := bytes.NewBufferString("")
	app := NewApplication("c2f")
	err := app.Run(r, w)
	if err != nil {
		t.Fatal(err)
	}
	assertEqualForEachLine(t, w.String(), fwvUseWidth01)
}

func TestFWV2CSV01(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(fwv01))
	w := bytes.NewBufferString("")
	app := NewApplication("f2c")
	app.UseWidth = false
	err := app.Run(r, w)
	if err != nil {
		t.Fatal(err)
	}
	assertEqualForEachLine(t, w.String(), csv01)
}

func TestFWV2CSVUseWidth01(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(fwvUseWidth01))
	w := bytes.NewBufferString("")
	app := NewApplication("f2c")
	err := app.Run(r, w)
	if err != nil {
		t.Fatal(err)
	}
	assertEqualForEachLine(t, w.String(), csv01)
}
