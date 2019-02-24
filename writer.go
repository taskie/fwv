package fwv

import (
	"io"
	"strings"

	"github.com/fatih/color"
)

var defaultColorOrder = []string{"red", "green", "blue", "yellow", "magenta", "cyan"}

var colorStringToFgMap = map[string]color.Attribute{
	"black":   color.FgBlack,
	"red":     color.FgRed,
	"green":   color.FgGreen,
	"yellow":  color.FgYellow,
	"blue":    color.FgBlue,
	"magenta": color.FgMagenta,
	"cyan":    color.FgCyan,
	"white":   color.FgWhite,
}

func colorStringToFunc(colorString string) func(...interface{}) string {
	if v, ok := colorStringToFgMap[strings.ToLower(colorString)]; ok {
		return color.New(v).SprintFunc()
	} else {
		return color.New(color.FgBlack).SprintFunc()
	}
}

type Writer struct {
	WidthCalculator  WidthCalculator
	UseCRLF          bool
	Colored          bool
	Delimiter        string
	underlyingWriter io.Writer
	colorOrder       []string
}

func NewWriter(w io.Writer) Writer {
	return NewWriterWithWidthCalculator(w, &SimpleWidthCalculator{})
}

func NewWriterWithWidthCalculator(w io.Writer, wcalc WidthCalculator) Writer {
	return Writer{
		WidthCalculator:  wcalc,
		Delimiter:        " ",
		colorOrder:       defaultColorOrder,
		underlyingWriter: w,
	}
}

func (w *Writer) CalcMaxWidthArrayOfColumns(records [][]string) []int {
	maxWidthByColumnIndex := make([]int, 0)
	for _, row := range records {
		for j, cell := range row {
			w := w.WidthCalculator.CalcWidthOfString(cell)
			for j >= len(maxWidthByColumnIndex) {
				maxWidthByColumnIndex = append(maxWidthByColumnIndex, w)
			}
			if w > maxWidthByColumnIndex[j] {
				maxWidthByColumnIndex[j] = w
			}
		}
	}
	return maxWidthByColumnIndex
}

func (w *Writer) ForEach(records [][]string, handler func(line string) error) error {
	oldNoColor := color.NoColor
	defer func() { color.NoColor = oldNoColor }()
	color.NoColor = !w.Colored

	maxWidthByColumnIndex := w.CalcMaxWidthArrayOfColumns(records)
	for _, record := range records {
		line := ""
		first := true
		for j, cell := range record {
			width := maxWidthByColumnIndex[j]
			padLen := width - w.WidthCalculator.CalcWidthOfString(cell)
			pad := strings.Repeat(" ", padLen)
			if !first {
				line += w.Delimiter
			}
			if w.Colored {
				colorString := w.colorOrder[j%len(w.colorOrder)]
				colorFunc := colorStringToFunc(colorString)
				line += colorFunc(cell) + pad
			} else {
				line += cell + pad
			}
			first = false
		}
		err := handler(line)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) WriteAll(records [][]string) error {
	br := "\n"
	if w.UseCRLF {
		br = "\r\n"
	}
	err := w.ForEach(records, func(line string) error {
		_, err := w.underlyingWriter.Write([]byte(line + br))
		return err
	})
	return err
}
