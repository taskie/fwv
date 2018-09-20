package fwv

import (
	"bufio"
	"io"
	"strings"
)

type Reader struct {
	UseWidthCalculator bool
	WidthCalculator    WidthCalculator
	whitespaces        string
	whitespaceWidthMap map[rune]int
	underlyingReader   io.Reader
}

func NewReader(r io.Reader) Reader {
	reader := Reader{
		UseWidthCalculator: false,
		WidthCalculator:    &SimpleWidthCalculator{},
		underlyingReader:   r,
	}
	reader.SetWhitespaces(" ")
	return reader
}

func NewReaderWithWidthCalculator(r io.Reader, wcalc WidthCalculator) Reader {
	reader := Reader{
		UseWidthCalculator: true,
		WidthCalculator:    wcalc,
		underlyingReader:   r,
	}
	reader.SetWhitespaces(" ")
	return reader
}

func (reader *Reader) SetWhitespaces(whitespaces string) {
	whitespaceWidthMap := make(map[rune]int)
	for _, c := range whitespaces {
		whitespaceWidthMap[c] = reader.WidthCalculator.CalcWidthOfRune(c)
	}
	reader.whitespaces = whitespaces
	reader.whitespaceWidthMap = whitespaceWidthMap
}

type ColumnSpec struct {
	maxWidth                int
	whitespaceCountByColumn map[int]int
}

func (reader *Reader) makeColumnSpec(lines []string) ColumnSpec {
	maxWidth := -1
	whitespaceCountByColumn := make(map[int]int)
	for _, line := range lines {
		pos := 0
		for _, c := range line {
			var w int
			var ok bool
			if w, ok = reader.whitespaceWidthMap[c]; ok {
				// c is whitespace
				for i := 0; i < w; i++ {
					if pos+i >= maxWidth {
						whitespaceCountByColumn[pos+i] = 1
					} else if _, ok := whitespaceCountByColumn[pos+i]; ok {
						whitespaceCountByColumn[pos+i]++
					}
				}
			} else {
				// c is non-whitespace
				w = reader.WidthCalculator.CalcWidthOfRune(c)
				for i := 0; i < w; i++ {
					delete(whitespaceCountByColumn, pos+i)
				}
			}
			pos += w
			if pos > maxWidth {
				maxWidth = pos
			}
		}
	}
	return ColumnSpec{
		maxWidth:                maxWidth,
		whitespaceCountByColumn: whitespaceCountByColumn,
	}
}

func (reader *Reader) makeColumnRanges(spec ColumnSpec) []IntRange {
	intRanges := make([]IntRange, 0)
	begin := -1
	inRange := false
	for i := 0; i < spec.maxWidth; i++ {
		if spec.whitespaceCountByColumn[i] <= 1 {
			if !inRange {
				begin = i
				inRange = true
			}
		} else if inRange {
			intRanges = append(intRanges, IntRange{
				Begin: begin,
				End:   i,
			})
			inRange = false
		}
	}
	if inRange {
		intRanges = append(intRanges, IntRange{
			Begin: begin,
			End:   spec.maxWidth,
		})
	}
	return intRanges
}

func (reader *Reader) extractCell(
	line string, columnRange IntRange, runeOffset int, widthOffset int,
) (cell string, read int, width int) {
	runes := []rune(line)
	targetRunes := runes[runeOffset:]
	cell = ""
	read = 0
	width = 0
	for _, c := range targetRunes {
		nextWidthOffset := widthOffset + width
		if nextWidthOffset < columnRange.Begin {
			// do nothing
		} else if columnRange.Begin <= nextWidthOffset && nextWidthOffset < columnRange.End {
			cell += string(c)
		} else if columnRange.End <= nextWidthOffset {
			break
		} else {
			panic("unreachable code")
		}
		width += reader.WidthCalculator.CalcWidthOfRune(c)
		read++
	}
	return
}

func (reader *Reader) loadLinesWithWidthCalculator(lines []string, columnRanges []IntRange, handler func(record []string) error) error {
	for _, line := range lines {
		record := make([]string, 0)
		runeOffset := 0
		widthOffset := 0
		for _, columnRange := range columnRanges {
			cell, read, width := reader.extractCell(line, columnRange, runeOffset, widthOffset)
			trimmedCell := strings.Trim(cell, reader.whitespaces)
			record = append(record, trimmedCell)
			runeOffset += read
			widthOffset += width
		}
		err := handler(record)
		if err != nil {
			return err
		}
	}
	return nil
}

func (reader *Reader) loadLinesWithoutWidthCalculator(lines []string, columnRanges []IntRange, handler func(record []string) error) error {
	for _, line := range lines {
		runes := []rune(line)
		l := len(runes)
		record := make([]string, 0)
		for _, columnRange := range columnRanges {
			begin := columnRange.Begin
			if begin > l {
				begin = l
			}
			end := columnRange.End
			if end > l {
				end = l
			}
			record = append(record, strings.Trim(string(runes[begin:end]), reader.whitespaces))
		}
		err := handler(record)
		if err != nil {
			return err
		}
	}
	return nil
}

func (reader *Reader) loadLines(lines []string, handler func(record []string) error) error {
	spec := reader.makeColumnSpec(lines)
	columnRanges := reader.makeColumnRanges(spec)
	var err error
	if reader.UseWidthCalculator {
		err = reader.loadLinesWithWidthCalculator(lines, columnRanges, handler)
	} else {
		err = reader.loadLinesWithoutWidthCalculator(lines, columnRanges, handler)
	}
	return err
}

func (reader *Reader) forEach(handler func(record []string) error) error {
	lines := make([]string, 0)
	scanner := bufio.NewScanner(reader.underlyingReader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return reader.loadLines(lines, handler)
}

func (reader *Reader) ReadAll() ([][]string, error) {
	records := make([][]string, 0)
	err := reader.forEach(func(record []string) error {
		records = append(records, record)
		return nil
	})
	return records, err
}
