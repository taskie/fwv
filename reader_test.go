package fwv

import (
	"bufio"
	"strings"
	"testing"
)

func assertEqualStringArray(t *testing.T, actual []string, expected []string) {
	if len(actual) != len(expected) {
		t.Errorf("assertEqualStringArray: actual: %v, expected: %v", actual, expected)
		return
	}
	for i, a := range actual {
		e := expected[i]
		if a != e {
			t.Errorf("assertEqualStringArray: actual: %v, expected: %v", actual, expected)
			break
		}
	}
}

func assertReader(t *testing.T, reader Reader, expectedRecords [][]string) {
	records, err := reader.ReadAll()
	if err != nil {
		t.Errorf("assertReader: %s", err.Error())
		return
	}
	for i, record := range records {
		assertEqualStringArray(t, record, expectedRecords[i])
	}
}

func TestReader01(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(fwv01))
	reader := NewReader(r)
	assertReader(t, reader, records01)
}

func TestReader02(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(fwv02))
	reader := NewReader(r)
	assertReader(t, reader, records02)
}

func TestReaderUseWidth01(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(fwvUseWidth01))
	reader := NewReaderWithWidthCalculator(r, &TextWidthCalculator{
		EastAsianAmbiguousWidth: 2,
	})
	assertReader(t, reader, records01)
}

func TestReaderUseWidth02(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(fwv02))
	reader := NewReaderWithWidthCalculator(r, &TextWidthCalculator{
		EastAsianAmbiguousWidth: 2,
	})
	assertReader(t, reader, records02)
}

func TestReaderUseWidth03(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(fwvUseWidth03))
	reader := NewReaderWithWidthCalculator(r, &TextWidthCalculator{
		EastAsianAmbiguousWidth: 2,
	})
	assertReader(t, reader, records03)
}

func TestReaderUseWidthEaaHalf01(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(fwvUseWidthEaaHalf01))
	reader := NewReaderWithWidthCalculator(r, &TextWidthCalculator{
		EastAsianAmbiguousWidth: 1,
	})
	assertReader(t, reader, records01)
}
