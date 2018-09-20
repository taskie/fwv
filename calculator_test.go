package fwv

import (
	"testing"
)

func assertWidthOfRune(t *testing.T, wcalc WidthCalculator, c rune, expected int) {
	actual := wcalc.CalcWidthOfRune(c)
	if actual != expected {
		t.Errorf("assertWidthOfString: %c (actual: %d, expected: %d)", c, actual, expected)
	}
}

func assertWidthOfString(t *testing.T, wcalc WidthCalculator, s string, expected int) {
	actual := wcalc.CalcWidthOfString(s)
	if actual != expected {
		t.Errorf("assertWidthOfString: %s (actual: %d, expected: %d)", s, actual, expected)
	}
}

func TestSimpleWidthCalculator(t *testing.T) {
	wcalc := &SimpleWidthCalculator{}
	assertWidthOfRune(t, wcalc, 'a', 1)
	assertWidthOfRune(t, wcalc, 'あ', 1)
	assertWidthOfRune(t, wcalc, 'α', 1)
	assertWidthOfRune(t, wcalc, '☺', 1)
	assertWidthOfRune(t, wcalc, 'Å', 1)
	assertWidthOfRune(t, wcalc, '⚡', 1)
	assertWidthOfRune(t, wcalc, '\u200B', 1)
	assertWidthOfRune(t, wcalc, '\U0001F600', 1)
	assertWidthOfRune(t, wcalc, '\U0001F1E6', 1)
	assertWidthOfString(t, wcalc, "abc", 3)
	assertWidthOfString(t, wcalc, "あいう", 3)
	assertWidthOfString(t, wcalc, "αβγ", 3)
	assertWidthOfString(t, wcalc, "Hello, 世界！", 10)
	assertWidthOfString(t, wcalc, "EΩD", 3)
}

func TestTextWidthCalculator(t *testing.T) {
	wcalc := &TextWidthCalculator{
		EastAsianAmbiguousWidth: 2,
	}
	assertWidthOfRune(t, wcalc, 'a', 1)
	assertWidthOfRune(t, wcalc, 'あ', 2)
	assertWidthOfRune(t, wcalc, 'α', 2)
	assertWidthOfRune(t, wcalc, '☺', 1)
	assertWidthOfRune(t, wcalc, 'Å', 1)
	assertWidthOfRune(t, wcalc, '⚡', 2)
	assertWidthOfRune(t, wcalc, '\u200B', 1)
	assertWidthOfRune(t, wcalc, '\U0001F600', 2)
	assertWidthOfRune(t, wcalc, '\U0001F1E6', 1)
	assertWidthOfString(t, wcalc, "abc", 3)
	assertWidthOfString(t, wcalc, "あいう", 6)
	assertWidthOfString(t, wcalc, "αβγ", 6)
	assertWidthOfString(t, wcalc, "Hello, 世界！", 13)
	assertWidthOfString(t, wcalc, "EΩD", 4)
}

func TestSimpleWidthCalculatorEaaHalf(t *testing.T) {
	wcalc := &TextWidthCalculator{
		EastAsianAmbiguousWidth: 1,
	}
	assertWidthOfRune(t, wcalc, 'a', 1)
	assertWidthOfRune(t, wcalc, 'あ', 2)
	assertWidthOfRune(t, wcalc, 'α', 1)
	assertWidthOfRune(t, wcalc, '☺', 1)
	assertWidthOfRune(t, wcalc, 'Å', 1)
	assertWidthOfRune(t, wcalc, '⚡', 2)
	assertWidthOfRune(t, wcalc, '\u200B', 1)
	assertWidthOfRune(t, wcalc, '\U0001F600', 2)
	assertWidthOfRune(t, wcalc, '\U0001F1E6', 1)
	assertWidthOfString(t, wcalc, "abc", 3)
	assertWidthOfString(t, wcalc, "あいう", 6)
	assertWidthOfString(t, wcalc, "αβγ", 3)
	assertWidthOfString(t, wcalc, "Hello, 世界！", 13)
	assertWidthOfString(t, wcalc, "EΩD", 3)
}
