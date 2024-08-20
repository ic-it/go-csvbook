# csvbook

`csvbook` is a Go library that provides a simple interface for managing and manipulating tabular data in CSV format. It supports reading, writing, and modifying CSV data, both in-memory and through zip archives. It also provides a way to handle metadata for both individual sheets and the entire file.

## Installation

To install `csvbook`, use the following Go command:

```bash
go get github.com/ic-it/go-csvbook
```

## Usage

### Creating a New Sheet

Create a new sheet using the `NewSheet` function:

```go
import "github.com/ic-it/go-csvbook"

sheet := csvbook.NewSheet()
```

### Writing and Reading Data

#### Writing Data

- **Write a Line**

```go
err := sheet.WriteLine([]any{"Column1", "Column2"})
```

- **Write a Cell**

```go
err := sheet.WriteCell("Value", lineNum, colNum)
```

#### Reading Data

- **Read a Line**

```go
line, err := sheet.ReadLine(lineNum)
```

- **Read a Cell**

```go
value, err := sheet.ReadCell(lineNum, colNum)
```

### Managing Metadata

- **Set Metadata**

```go
sheet.SetMeta("key", "value")
```

- **Get Metadata**

```go
value := sheet.GetMeta("key")
```

### Working with Files

You can manage multiple sheets and save them to a zip file or read from one.

- **Create a New File**

```go
file, err := csvbook.NewFile()
```

- **Add a Sheet**

```go
sheet := file.CreateSheet("Sheet1")
```

- **Write to a File**

```go
err := file.WriteZipFile("data.zip")
```

- **Read from a File**

```go
err := file.ReadZipFile("data.zip")
```

- **Set File Metadata**

```go
file.SetMeta("key", "value")
```

- **Get File Metadata**

```go
value := file.GetMeta("key")
```

## Examples

Here is a simple example demonstrating how to use the `csvbook` library:

```go
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

```

## License

`csvbook` is licensed under the MIT License. See the [LICENSE](LICENSE.txt) file for details.