// Package chunker provides utilities for splitting text into smaller, manageable chunks.
package chunker

import (
	"strings"

	"unicode/utf8"
)

// TYPELINK is a constant used to mark a rune as part of a link.
const TYPELINK = 100000

// TYPEEND is a constant used to mark the end of a chunk.
const TYPEEND = -1

// isStartWithLink checks if the given text starts with "http://" or "https://".
func isStartWithLink(text string) bool {
	return strings.HasPrefix(text, "http://") || strings.HasPrefix(text, "https://")
}

// sentenceSplitCharacterM is a map of characters that can be used to split sentences.
var sentenceSplitCharacterM = map[string]bool{
	"\n": true,
	".":  true,
	";":  true,
	" ":  true,
}

// FirstChunk extracts the first chunk of text based on minSize and maxSize.
// It attempts to break at sentence boundaries or spaces, while also handling links.
func FirstChunk(text string, minSize, maxSize int) string {
	if maxSize == 0 {
		return ""
	}

	if minSize == 0 {
		minSize = maxSize
	}

	if len(text) <= minSize {
		return text
	}

	candidate := ""
	var link string
	isLink := false

	lastindex := LastRuneByteIndex(text)
	runetypemap := map[int]int{}
	var fragment string
	for i, r := range text {
		// build link
		if !isLink && isStartWithLink(text[i:]) {
			isLink = true
			link = ""
		}

		if isLink {
			if r == '\n' || r == ' ' || r == ')' {
				isLink = false
			}
		}

		if isLink {
			// max link is 1024
			if len(link) > 1024 {
				isLink = false
			}
		}

		if isLink {
			runetypemap[i] = TYPELINK
			link += string(r)
			continue
		}

		runetypemap[i] = 1 // normal
		if !isLink {
			fragment += link + string(r)
			link = ""
		}

		// build fragment
		fragmentEnded := r == '.' || r == '\n' || r == ';' || i == lastindex
		if len(fragment)+len(candidate) > maxSize { // add and extra rune
			fragmentEnded = true
		}
		if fragmentEnded {
			// commit to candidate
			if len(candidate) < minSize || len(candidate)+len(fragment) < maxSize {
				//  keep do it
				candidate += fragment
				fragment = ""
				continue
			}

			// too big
			break
		}
	}
	runetypemap[lastindex] = TYPEEND // mark end
	// runetypemap[len(text)] = TYPEEND // mark end

	lasti := len(candidate)
	for i := range candidate {
		ri := len(candidate) - i - 1
		r := candidate[ri]

		if runetypemap[ri] == TYPEEND {
			break
		}
		if runetypemap[ri] == TYPELINK {
			continue // do not break in middle of the link
		}
		if r == ' ' || r == '.' || r == '\n' {
			if ri < minSize {
				break
			}
			lasti = ri + 1
			break
		}
	}

	return candidate[:lasti]
}

// Chunk splits the given text into chunks of a specified size with a given overlap.
// chunkSize: maximum number of characters for each chunk.
// chunkOverlap: the number of characters to overlap between consecutive chunks.
func Chunk(text string, chunkSize int, chunkOverlap int) []string {
	if chunkSize < 2 {
		chunkSize = 2
	}

	if chunkOverlap > chunkSize-1 {
		chunkOverlap = chunkSize - 1
	}

	chunks := []string{}
	text = Substring(text, 0, 1000_000) // work with top 1M characters
	for len(text) > 0 {
		chunk := FirstChunk(text, chunkSize/2, chunkSize)
		chunks = append(chunks, chunk)
		if len(chunk) < 2 {
			text = text[len(chunk):]
			continue
		}

		if len(chunk) >= len(text) {
			text = ""
			continue
		}
		removechunk := FirstChunk(text, (chunkSize-chunkOverlap)/2, (chunkSize - chunkOverlap))
		if len(removechunk) > 0 {
			text = text[len(removechunk):]
			continue
		}
		// make sure we alway move forward
		text = text[len(chunk):]
	}
	return chunks
}

// Substring returns a substring of s from start (inclusive) to end (exclusive)
// based on rune count, not byte count.
func Substring(s string, start int, end int) string {
	if s == "" {
		return ""
	}

	if start == 0 && end >= len(s) {
		return s
	}

	start_str_idx := 0
	i := 0
	for j := range s {
		if i == start {
			start_str_idx = j
		}
		if i == end {
			return s[start_str_idx:j]
		}
		i++
	}
	return s[start_str_idx:]
}

// LastRuneByteIndex returns the byte index of the last rune in a string.
// It returns -1 if the string is empty.
func LastRuneByteIndex(s string) int {
	if s == "" {
		return -1
	}
	_, size := utf8.DecodeLastRuneInString(s)
	if size <= 0 {
		return -1
	}
	return len(s) - size
}
