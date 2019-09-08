package hyperlink

import (
	"net/url"
	"strings"
)

//CreateLinks take a string parameter i.e. stacktrace.
//It will perform operations on that input and convert the input to links
//where there is a tab in the starting of the line and a colon(":") is present in the line.
//It returns the same input as ouput after converting some lines to hyperlink.
func CreateLinks(stack string) string {
	lines := strings.Split(stack, "\n")
	for index, line := range lines {
		if len(line) == 0 || line[0] != '\t' {
			continue
		}

		fileName := ""
		for i, ch := range line {
			if ch == ':' && line[i+1] != '/' {
				fileName = line[1:i]
				break
			}
		}
		var lineNo strings.Builder
		for i := len(fileName) + 2; i < len(line); i++ {
			if line[i] < '0' || line[i] > '9' {
				break
			}
			lineNo.WriteByte(line[i])
		}

		fileURL := url.Values{}
		fileURL.Set("path", fileName)
		fileURL.Set("line", lineNo.String())

		lines[index] = "\t<a href=\"/debug?" + fileURL.Encode() + "\">" + fileName + ":" + lineNo.String() + "</a>" + line[len(fileName)+2+len(lineNo.String()):]
		//lines[index] = "\t<a href=/debug?" + fileURL.Encode() + ">" + fileName + ":" + lineNo.String() + "</a>" + line[len(fileName)+2+len(lineNo.String()):]

	}
	return strings.Join(lines, "\n")
}
