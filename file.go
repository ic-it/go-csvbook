package csvbook

import (
	"archive/zip"
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
)

func NewFile(
	meta ...map[string]string,
) (*File, error) {
	file := &File{
		sheets: make(map[string]*Sheet),
	}

	if len(meta) > 0 {
		for k, v := range meta[0] {
			file.meta[k] = v
		}
	} else {
		file.meta = make(map[string]any)
	}

	return file, nil
}

type File struct {
	sheets map[string]*Sheet

	meta map[string]any

	mu sync.Mutex
}

func (f *File) CreateSheet(name string) *Sheet {
	f.mu.Lock()
	defer f.mu.Unlock()

	sheet := NewSheet()
	f.sheets[name] = sheet
	return sheet
}

func (f *File) GetSheet(name string) *Sheet {
	return f.sheets[name]
}

func (f *File) DeleteSheet(name string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.sheets[name].muLock()
	defer f.sheets[name].muUnlock()

	delete(f.sheets, name)
}

func (f *File) Write(writer io.Writer) error {
	// lock all sheets
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, sheet := range f.sheets {
		sheet.muLock()
		defer sheet.muUnlock()
	}

	// create zip writer
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()

	// copy metadata
	meta := make(map[string]any)
	for k, v := range f.meta {
		meta[k] = v
	}
	sheetsMeta := make(map[string]any)
	for name, sheet := range f.sheets {
		sheet.updateMeta()
		sheetsMeta[name] = sheet.getMeta()
	}
	meta[_SHEETS_META] = sheetsMeta
	sheetsNamesMap := make(map[string]string)
	meta[_SHEETS_NAMES_MAP] = sheetsNamesMap

	// write sheets
	for name, sheet := range f.sheets {
		sheetName := escapeFilename(name) + ".csv"
		writer, err := zipWriter.Create(sheetName)
		if err != nil {
			return err
		}
		sheet.write(writer)
		sheetsNamesMap[name] = sheetName
	}

	metaBytes, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	// build metadata
	metaWriter, err := zipWriter.Create("meta.json")
	if err != nil {
		return err
	}

	_, err = metaWriter.Write(metaBytes)
	if err != nil {
		return err
	}

	return nil
}

func (f *File) Read(reader io.ReaderAt, size int64) error {
	// lock all sheets
	f.mu.Lock()
	defer f.mu.Unlock()

	// create zip reader
	zipReader, err := zip.NewReader(reader, size)
	if err != nil {
		return err
	}

	// read metadata
	metaReader, err := zipReader.Open("meta.json")
	if err != nil {
		return err
	}
	defer metaReader.Close()

	metaBytes, err := io.ReadAll(metaReader)
	if err != nil {
		return err
	}

	meta := make(map[string]any)
	err = json.Unmarshal(metaBytes, &meta)
	if err != nil {
		return err
	}

	sheetsMeta, ok := meta[_SHEETS_META].(map[string]any)
	if !ok {
		return ErrInvalidMetadata
	}

	sheetsNamesMap, ok := meta[_SHEETS_NAMES_MAP].(map[string]any)
	if !ok {
		return ErrInvalidMetadata
	}

	// read sheets
	for name, sheetName := range sheetsNamesMap {
		sheetReader, err := zipReader.Open(sheetName.(string))
		if err != nil {
			return err
		}

		sheet := NewSheet()
		err = sheet.read(sheetReader)
		if err != nil {
			return err
		}

		f.sheets[name] = sheet
	}

	// update metadata
	for name, sheet := range f.sheets {
		sheet.updateMeta()
		sheetMeta, ok := sheetsMeta[name].(map[string]any)
		if !ok {
			log.Println("sheet-meta", sheetsMeta[name])
			return ErrInvalidMetadata
		}
		sheet.setMeta(sheetMeta)
	}

	return nil
}

func (f *File) WriteZipFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return f.Write(file)
}

func (f *File) ReadZipFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	size, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	return f.Read(file, size)
}

func (f *File) SetMeta(key string, value any) {
	f.meta[key] = value
}

func (f *File) GetMeta(key string) any {
	return f.meta[key]
}
