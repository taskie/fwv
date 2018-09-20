package fwv

import (
	"golang.org/x/text/width"
	"unicode/utf8"
)

type WidthCalculator interface {
	CalcWidthOfRune(c rune) int
	CalcWidthOfString(s string) int
}

type SimpleWidthCalculator struct{}

func (wcalc *SimpleWidthCalculator) CalcWidthOfRune(c rune) int {
	return 1
}

func (wcalc *SimpleWidthCalculator) CalcWidthOfString(s string) int {
	return utf8.RuneCountInString(s)
}

type TextWidthCalculator struct {
	EastAsianAmbiguousWidth int
}

func (wcalc *TextWidthCalculator) CalcWidthOfRune(c rune) int {
	kind := width.LookupRune(c).Kind()
	switch kind {
	case width.Neutral, width.EastAsianNarrow, width.EastAsianHalfwidth:
		return 1
	case width.EastAsianWide, width.EastAsianFullwidth:
		return 2
	case width.EastAsianAmbiguous:
		return wcalc.EastAsianAmbiguousWidth
	default:
		return 1
	}
}

func (wcalc *TextWidthCalculator) CalcWidthOfString(s string) int {
	w := 0
	for _, c := range s {
		w += wcalc.CalcWidthOfRune(c)
	}
	return w
}
