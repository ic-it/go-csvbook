package csvbook

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"sync"
)

func NewSheet() *Sheet {
	return &Sheet{
		lines: make([][]any, 0),
		meta:  make(map[string]any),
	}
}

type Sheet struct {
	lines  [][]any
	maxCol int

	meta map[string]any

	mu sync.Mutex
}

// Get the number of lines in the sheet.
func (s *Sheet) LineNum() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.lines)
}

// If lineNum is not provided, append line to the end of the sheet.
func (s *Sheet) WriteLine(line []any, lineNum ...int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// if line is more than maxCol, update all lines to new maxCol
	if len(line) > s.maxCol {
		for i := range s.lines {
			s.lines[i] = append(s.lines[i], make([]any, len(line)-s.maxCol)...)
		}
		s.maxCol = len(line)
	}

	if len(lineNum) == 0 {
		lineNum = append(lineNum, len(s.lines))
	}
	lineNum_ := lineNum[0]

	if lineNum_ < 0 {
		return errors.Join(fmt.Errorf("line number %d", lineNum_), ErrInvalidLineNum)
	}

	if lineNum_ < len(s.lines) {
		s.lines[lineNum_] = line
	} else {
		for i := len(s.lines); i < lineNum_; i++ {
			s.lines = append(s.lines, make([]any, s.maxCol))
		}
		s.lines = append(s.lines, line)
	}

	return nil
}

// Write Cell
func (s *Sheet) WriteCell(value any, lineNum, colNum int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// if colNum is more than maxCol, update all lines to new maxCol
	if colNum >= s.maxCol {
		for i := range s.lines {
			s.lines[i] = append(s.lines[i], make([]any, colNum-s.maxCol+1)...)
		}
		s.maxCol = colNum + 1
	}

	if lineNum < 0 {
		return errors.Join(fmt.Errorf("line number %d", lineNum), ErrInvalidLineNum)
	}
	if colNum < 0 {
		return errors.Join(fmt.Errorf("col number %d", colNum), ErrorInvalidColNum)
	}

	if lineNum >= len(s.lines) {
		for i := len(s.lines); i < lineNum; i++ {
			s.lines = append(s.lines, make([]any, s.maxCol))
		}
		s.lines = append(s.lines, make([]any, s.maxCol))
	}

	s.lines[lineNum][colNum] = value
	return nil
}

// Get Line Copy
func (s *Sheet) ReadLine(lineNum int) ([]any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if lineNum < 0 || lineNum >= len(s.lines) {
		return nil, errors.Join(fmt.Errorf("line number %d", lineNum), ErrInvalidLineNum)
	}

	line := make([]any, len(s.lines[lineNum]))
	copy(line, s.lines[lineNum])
	return line, nil
}

// Read Cell
func (s *Sheet) ReadCell(lineNum, colNum int) (any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if lineNum < 0 || lineNum >= len(s.lines) {
		return nil, errors.Join(fmt.Errorf("line number %d", lineNum), ErrInvalidLineNum)
	}
	if colNum < 0 || colNum >= len(s.lines[lineNum]) {
		return nil, errors.Join(fmt.Errorf("col number %d", colNum), ErrorInvalidColNum)
	}

	return s.lines[lineNum][colNum], nil
}

func (s *Sheet) SetMeta(key string, value any) {
	s.meta[key] = value
}

func (s *Sheet) GetMeta(key string) any {
	return s.meta[key]
}

func (s *Sheet) write(writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	for _, line := range s.lines {
		// convert []any to []string
		line_ := make([]string, len(line))
		for i, v := range line {
			if v == nil {
				v = ""
			}
			line_[i] = fmt.Sprintf("%v", v)
		}
		if err := csvWriter.Write(line_); err != nil {
			return err
		}
	}

	return nil
}

func (s *Sheet) read(reader io.Reader) error {
	csvReader := csv.NewReader(reader)

	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// convert []string to []any
		line_ := make([]any, len(line))
		for i, v := range line {
			line_[i] = v
		}

		s.lines = append(s.lines, line_)
	}

	return nil
}

func (s *Sheet) muLock() {
	s.mu.Lock()
}

func (s *Sheet) muUnlock() {
	s.mu.Unlock()
}

func (s *Sheet) getMeta() map[string]any {
	return s.meta
}

// updateMeta updates the metadata of the sheet based on the current state of the sheet.
func (s *Sheet) updateMeta() {
	s.meta[_SHEET_META_LINE_NUM] = len(s.lines)
	s.meta[_SHEET_META_COL_NUM] = s.maxCol
}

func (s *Sheet) setMeta(meta map[string]any) {
	s.meta = meta
}
