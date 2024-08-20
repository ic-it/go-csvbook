package csvbook

import "strings"

func escapeFilename(name string) string {
	newName := strings.Builder{}
	for _, c := range name {
		if c == '/' {
			newName.WriteString("_slash_")
		} else {
			newName.WriteRune(c)
		}
	}
	return newName.String()
}
