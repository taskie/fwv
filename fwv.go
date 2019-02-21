package fwv

import (
	"encoding/csv"
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
	Mode                    string
	UseWidth                bool
	EastAsianAmbiguousWidth int
	Whitespaces             string
	UseCRLF                 bool
	Comma                   rune
	CSVComment              rune
	Delimiter               string
	Colored                 bool
	ColumnRanges            []IntRange
}

func NewConverter(w io.Writer, r io.Reader, mode string) Converter {
	return Converter{
		Reader:                  r,
		Writer:                  w,
		Mode:                    mode,
		UseWidth:                true,
		EastAsianAmbiguousWidth: 2,
		Whitespaces:             " ",
		UseCRLF:                 false,
		Comma:                   ',',
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

func (c *Converter) ConvertFWVToCSV() error {
	reader := c.fwvReader()
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	csvw := c.csvWriter()
	err = csvw.WriteAll(records)
	return err
}

func (c *Converter) ConvertCSVToFWV() error {
	csvr := c.csvReader()
	records, err := csvr.ReadAll()
	if err != nil {
		return err
	}
	writer := c.fwvWriter()
	err = writer.WriteAll(records)
	return err
}

func (c *Converter) ConvertFWVToFWV() error {
	reader := c.fwvReader()
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	writer := c.fwvWriter()
	err = writer.WriteAll(records)
	return err
}

func (c *Converter) Convert() error {
	switch strings.ToLower(c.Mode) {
	case "f2c":
		return c.ConvertFWVToCSV()
	case "f2f", "shrink":
		return c.ConvertFWVToFWV()
	default:
		return c.ConvertCSVToFWV()
	}
}
