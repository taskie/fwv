package fwv

import (
	"bufio"
	"errors"
	"regexp"
	"strings"
	"testing"
)

func assertEqual(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Errorf("assertEqual: actual: %v, expected: %v", actual, expected)
	}
}

func assertNotEqual(t *testing.T, actual interface{}, notExpected interface{}) {
	if actual == notExpected {
		t.Errorf("assertNotEqual: must not be equal to %v", notExpected)
	}
}

func assertEqualForEachLine(t *testing.T, actual string, expected string) {
	actualScanner := bufio.NewScanner(strings.NewReader(actual))
	expectedScanner := bufio.NewScanner(strings.NewReader(expected))
	for actualScanner.Scan() && expectedScanner.Scan() {
		// XXX: ignoring trailing spaces...
		actual := strings.TrimRight(actualScanner.Text(), " ")
		expected := strings.TrimRight(expectedScanner.Text(), " ")
		assertEqual(t, actual, expected)
	}
}

func assertWriter(t *testing.T, writer Writer, records [][]string, expectedFWV string) {
	actual := ""
	writer.ForEach(records, func(line string) error {
		actual += line + "\n"
		return nil
	})
	assertEqualForEachLine(t, actual, expectedFWV)
}

func TestWriter01(t *testing.T) {
	writer := NewWriter(nil)
	assertWriter(t, writer, records01, fwv01)
}

func TestWriter02(t *testing.T) {
	writer := NewWriter(nil)
	assertWriter(t, writer, records02, fwv02)
}

func TestWriterUseWidth01(t *testing.T) {
	writer := NewWriterWithWidthCalculator(nil, &TextWidthCalculator{
		EastAsianAmbiguousWidth: 2,
	})
	assertWriter(t, writer, records01, fwvUseWidth01)
}

func TestWriterUseWidth03(t *testing.T) {
	writer := NewWriterWithWidthCalculator(nil, &TextWidthCalculator{
		EastAsianAmbiguousWidth: 2,
	})
	assertWriter(t, writer, records03, fwvUseWidth03)
}

func TestWriterUseWidthEaaHalf01(t *testing.T) {
	writer := NewWriterWithWidthCalculator(nil, &TextWidthCalculator{
		EastAsianAmbiguousWidth: 1,
	})
	assertWriter(t, writer, records01, fwvUseWidthEaaHalf01)
}

func TestWriterUseWidthDelimited03(t *testing.T) {
	writer := NewWriterWithWidthCalculator(nil, &TextWidthCalculator{
		EastAsianAmbiguousWidth: 2,
	})
	writer.Delimiter = "|"
	assertWriter(t, writer, records03, fwvUseWidthDelimited03)
}

func TestWriterUseWidthColored03(t *testing.T) {
	writer := NewWriterWithWidthCalculator(nil, &TextWidthCalculator{
		EastAsianAmbiguousWidth: 2,
	})
	writer.Colored = true
	actual := ""
	writer.ForEach(records03, func(line string) error {
		actual += line + "\n"
		return nil
	})
	assertNotEqual(t, actual, fwvUseWidth03)
	// remove color
	re := regexp.MustCompile("\x1b\\[\\d+m")
	noColoredActual := string(re.ReplaceAll([]byte(actual), []byte{}))
	assertEqualForEachLine(t, noColoredActual, fwvUseWidth03)
}

func TestWriterForEachError(t *testing.T) {
	writer := NewWriter(nil)
	err := writer.ForEach([][]string{{"a"}, {"b"}, {"c"}}, func(line string) error {
		if line == "b" {
			return errors.New("TestWriterForEachError")
		} else if line == "c" {
			t.Fail()
		}
		return nil
	})
	if err == nil {
		t.Fatal("err is nil")
	}
	if err.Error() != "TestWriterForEachError" {
		t.Fatal("invalid error")
	}
}
