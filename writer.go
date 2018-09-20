package fwv

import (
	"io"
	"strings"
)

type Writer struct {
	WidthCalculator  WidthCalculator
	UseCRLF          bool
	underlyingWriter io.Writer
}

func NewWriter(w io.Writer) Writer {
	return Writer{
		WidthCalculator:  &SimpleWidthCalculator{},
		underlyingWriter: w,
	}
}

func NewWriterWithWidthCalculator(w io.Writer, wcalc WidthCalculator) Writer {
	return Writer{
		WidthCalculator:  wcalc,
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
	maxWidthByColumnIndex := writer.CalcMaxWidthArrayOfColumns(records)
	for _, record := range records {
		line := ""
		first := true
		for j, cell := range record {
			w := maxWidthByColumnIndex[j]
			padLen := w - writer.WidthCalculator.CalcWidthOfString(cell)
			pad := strings.Repeat(" ", padLen)
			if !first {
				line += " "
			}
			line += cell + pad
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
