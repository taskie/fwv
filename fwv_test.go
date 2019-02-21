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
	conv := Converter{
		Reader:                  r,
		Writer:                  w,
		FromType:                "csv",
		UseWidth:                false,
		EastAsianAmbiguousWidth: 1,
		Whitespaces:             " ",
	}
	conv.Convert()
	assertEqualForEachLine(t, w.String(), fwv01)
}

func TestCSV2FWVUseWidth01(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(csv01))
	w := bytes.NewBufferString("")
	conv := NewConverter(w, r, "csv", "fwv")
	err := conv.Convert()
	if err != nil {
		t.Fatal(err)
	}
	assertEqualForEachLine(t, w.String(), fwvUseWidth01)
}

func TestFWV2CSV01(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(fwv01))
	w := bytes.NewBufferString("")
	conv := NewConverter(w, r, "fwv", "csv")
	conv.UseWidth = false
	err := conv.Convert()
	if err != nil {
		t.Fatal(err)
	}
	assertEqualForEachLine(t, w.String(), csv01)
}

func TestFWV2CSVUseWidth01(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(fwvUseWidth01))
	w := bytes.NewBufferString("")
	conv := NewConverter(w, r, "fwv", "csv")
	err := conv.Convert()
	if err != nil {
		t.Fatal(err)
	}
	assertEqualForEachLine(t, w.String(), csv01)
}
