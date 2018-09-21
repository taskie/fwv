package fwv

import (
	"github.com/fatih/color"
	"io"
	"strings"
)

var defaultColorOrder = []string{"green", "yellow", "blue", "magenta", "cyan"}

var colorStringToFgMap = map[string]color.Attribute{
	"BLACK":   color.FgBlack,
	"RED":     color.FgRed,
	"GREEN":   color.FgGreen,
	"YELLOW":  color.FgYellow,
	"BLUE":    color.FgBlue,
	"MAGENTA": color.FgMagenta,
	"CYAN":    color.FgCyan,
	"WHITE":   color.FgWhite,
}

func colorStringToFunc(colorString string) func(...interface{}) string {
	if v, ok := colorStringToFgMap[strings.ToUpper(colorString)]; ok {
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
		colorOrder:       defaultColorOrder,
		underlyingWriter: w,
	}
}

func (writer *Writer) CalcMaxWidthArrayOfColumns(records [][]string) []int {
	maxWidthByColumnIndex := make([]int, 0)
	for _, row := range records {
		for j, cell := range row {
			w := writer.WidthCalculator.CalcWidthOfString(cell)
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

func (writer *Writer) ForEach(records [][]string, handler func(line string) error) error {
	oldNoColor := color.NoColor
	defer func() { color.NoColor = oldNoColor }()
	color.NoColor = !writer.Colored

	maxWidthByColumnIndex := writer.CalcMaxWidthArrayOfColumns(records)
	for _, record := range records {
		line := ""
		first := true
		for j, cell := range record {
			w := maxWidthByColumnIndex[j]
			padLen := w - writer.WidthCalculator.CalcWidthOfString(cell)
			pad := strings.Repeat(" ", padLen)
			if !first {
				if writer.Delimiter != "" {
					line += writer.Delimiter
				} else {
					line += " "
				}
			}
			if writer.Colored {
				colorString := writer.colorOrder[j%len(writer.colorOrder)]
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

func (writer *Writer) WriteAll(records [][]string) error {
	br := "\n"
	if writer.UseCRLF {
		br = "\r\n"
	}
	err := writer.ForEach(records, func(line string) error {
		_, err := writer.underlyingWriter.Write([]byte(line + br))
		return err
	})
	return err
}
