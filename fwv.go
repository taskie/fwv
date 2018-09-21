package fwv

import (
	"encoding/csv"
	"io"
)

var (
	Version  = "0.1.0-beta"
	Revision = ""
)

type Application struct {
	Mode                    string
	UseWidth                bool
	EastAsianAmbiguousWidth int
	Whitespaces             string
	UseCRLF                 bool
	Comma                   rune
	CSVComment              rune
	Delimiter               string
	Colored                 bool
}

func NewApplication(mode string) Application {
	return Application{
		Mode:                    mode,
		UseWidth:                true,
		EastAsianAmbiguousWidth: 2,
		Whitespaces:             " ",
		UseCRLF:                 false,
		Comma:                   ',',
	}
}

func (app *Application) ConvertFWVToCSV(r io.Reader, w io.Writer) error {
	csvw := csv.NewWriter(w)
	csvw.UseCRLF = app.UseCRLF
	csvw.Comma = app.Comma
	var reader Reader
	if app.UseWidth {
		reader = NewReaderWithWidthCalculator(r, &TextWidthCalculator{
			EastAsianAmbiguousWidth: app.EastAsianAmbiguousWidth,
		})
	} else {
		reader = NewReader(r)
	}
	reader.SetWhitespaces(app.Whitespaces)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	err = csvw.WriteAll(records)
	return err
}

func (app *Application) ConvertCSVToFWV(r io.Reader, w io.Writer) error {
	csvr := csv.NewReader(r)
	csvr.Comment = app.CSVComment
	csvr.Comma = app.Comma
	records, err := csvr.ReadAll()
	if err != nil {
		return err
	}

	var writer Writer
	if app.UseWidth {
		writer = NewWriterWithWidthCalculator(w, &TextWidthCalculator{
			EastAsianAmbiguousWidth: app.EastAsianAmbiguousWidth,
		})
	} else {
		writer = NewWriter(w)
	}
	writer.UseCRLF = app.UseCRLF
	writer.Delimiter = app.Delimiter
	writer.Colored = app.Colored
	err = writer.WriteAll(records)
	return err
}

func (app *Application) Run(r io.Reader, w io.Writer) error {
	if app.Mode == "f2c" {
		return app.ConvertFWVToCSV(r, w)
	} else {
		return app.ConvertCSVToFWV(r, w)
	}
}
