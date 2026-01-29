package ui

import "strings"

// wordWrap splits text into lines of max width, preserving existing newlines
func wordWrap(text string, width int) []string {
	var lines []string

	paragraphs := strings.Split(text, "\n")

	for _, para := range paragraphs {
		words := strings.Fields(para)
		if len(words) == 0 {
			lines = append(lines, "")
			continue
		}

		var line string
		for _, word := range words {
			if len(line)+len(word)+1 > width {
				lines = append(lines, line)
				line = word
			} else {
				if line != "" {
					line += " "
				}
				line += word
			}
		}
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines
}
