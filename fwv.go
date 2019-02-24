package fwv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

var (
	Version  = "0.1.0-beta"
	Revision = ""
)

type Converter struct {
	Reader                  io.Reader
	Writer                  io.Writer
	FromType                string
	ToType                  string
	UseWidth                bool
	EastAsianAmbiguousWidth int
	Whitespaces             string
	NoTrim                  bool
	UseCRLF                 bool
	Comma                   rune
	CSVComment              rune
	Delimiter               string
	Colored                 bool
	ColumnRanges            []IntRange
	ShowColumnRanges        bool
}

func NewConverter(w io.Writer, r io.Reader, fromType string, toType string) *Converter {
	return &Converter{
		Reader:                  r,
		Writer:                  w,
		FromType:                fromType,
		ToType:                  toType,
		UseWidth:                true,
		EastAsianAmbiguousWidth: 2,
		Whitespaces:             " ",
		UseCRLF:                 false,
		Comma:                   ',',
		Delimiter:               " ",
	}
}

func (c *Converter) fwvReader() *Reader {
	var reader Reader
	if c.UseWidth {
		reader = NewReaderWithWidthCalculator(c.Reader, &TextWidthCalculator{
			EastAsianAmbiguousWidth: c.EastAsianAmbiguousWidth,
		})
	} else {
		reader = NewReader(c.Reader)
	}
	reader.SetWhitespaces(c.Whitespaces)
	reader.ColumnRanges = c.ColumnRanges
	reader.NoTrim = c.NoTrim
	return &reader
}

func (c *Converter) fwvWriter() *Writer {
	var writer Writer
	if c.UseWidth {
		writer = NewWriterWithWidthCalculator(c.Writer, &TextWidthCalculator{
			EastAsianAmbiguousWidth: c.EastAsianAmbiguousWidth,
		})
	} else {
		writer = NewWriter(c.Writer)
	}
	writer.UseCRLF = c.UseCRLF
	writer.Delimiter = c.Delimiter
	writer.Colored = c.Colored
	return &writer
}

func (c *Converter) csvReader() *csv.Reader {
	csvr := csv.NewReader(c.Reader)
	csvr.Comment = c.CSVComment
	csvr.Comma = c.Comma
	return csvr
}

func (c *Converter) csvWriter() *csv.Writer {
	csvw := csv.NewWriter(c.Writer)
	csvw.UseCRLF = c.UseCRLF
	csvw.Comma = c.Comma
	return csvw
}

func (c *Converter) Convert() error {
	var records [][]string
	var readInfo *ReadInfo
	var err error
	switch strings.ToLower(c.FromType) {
	case "csv":
		records, err = c.csvReader().ReadAll()
	default:
		records, readInfo, err = c.fwvReader().ReadAllInfo()
	}
	if err != nil {
		return err
	}
	if c.ShowColumnRanges {
		if readInfo == nil {
			return fmt.Errorf("can't show column ranges")
		}
		ranges := make([]string, 0)
		for _, cr := range readInfo.ColumnRanges {
			ranges = append(ranges, fmt.Sprintf("%d:%d", cr.Begin, cr.End))
		}
		records = [][]string{ranges}
	}
	switch strings.ToLower(c.ToType) {
	case "csv":
		err = c.csvWriter().WriteAll(records)
	default:
		err = c.fwvWriter().WriteAll(records)
	}
	return err
}
