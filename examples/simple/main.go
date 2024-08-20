package main

import "github.com/ic-it/go-csvbook"

func main() {
	f, err := csvbook.NewFile()

	if err != nil {
		panic(err)
	}

	sheet := f.CreateSheet("sheet1")

	sheet.WriteLine([]any{"a", "b", "c"})
	sheet.WriteLine([]any{"1", "2", "3"})
	sheet.WriteLine([]any{"4", "5", "6"})
	sheet.WriteLine([]any{"7", "8", "9", 123})

	sheet = f.CreateSheet("sheet2")

	sheet.WriteLine([]any{"a", "b", "c"})
	sheet.WriteLine([]any{"1", "2", "3"})
	sheet.WriteLine([]any{"4", "5", "6"})
	sheet.WriteLine([]any{"7", "8", "9", 123})

	sheet.WriteCell("111", 100, 15)

	err = f.WriteZipFile("test.zip")
	if err != nil {
		panic(err)
	}

	f, _ = csvbook.NewFile()
	err = f.ReadZipFile("test.zip")

	if err != nil {
		panic(err)
	}

	sheet = f.GetSheet("sheet1")

	for i := range sheet.LineNum() {
		line, _ := sheet.ReadLine(i)
		for _, cell := range line {
			print(cell.(string))
			print(" ")
		}
		println()
	}
}
