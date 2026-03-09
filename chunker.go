// Package chunker provides utilities for splitting text into smaller, manageable chunks.
package chunker

import (
	"strings"

	"unicode"
	"unicode/utf8"
)

// FirstChunk extracts the first chunk of text based on minSize and maxSize.
// It attempts to break at sentence boundaries or spaces, while also handling links.
func FirstChunk(text string, minSize, maxSize int) string {
	if maxSize <= 0 {
		return ""
	}

	if minSize <= 0 {
		minSize = maxSize
	}

	if len(text) <= minSize {
		return SanitizeUnicode(text)
	}

	isLink := false
	link := ""
	candidate := ""
	var fragment string
	
	// Track which byte indices are part of a link.
	isLinkRune := make([]bool, len(text))
	
	lastindex := -1
	for i := range text {
		lastindex = i
	}

	for i, r := range text {
		// Link detection
		if !isLink && r == 'h' && isStartWithLink(text[i:]) {
			isLink = true
			link = ""
		}

		if isLink {
			if r == '\n' || r == ' ' || r == ')' {
				isLink = false
			}
		}

		if isLink && len(link) > 1024 {
			isLink = false
		}

		if isLink {
			if i < len(isLinkRune) {
				isLinkRune[i] = true
			}
			link += string(r)
			
			// If adding this part of the link makes us exceed maxSize, we should stop.
			if len(candidate)+len(fragment)+len(link) > maxSize {
				// We reached the limit while in a link.
				// Commit what we have in fragment to candidate (if it fits).
				if len(candidate)+len(fragment) <= maxSize {
					candidate += fragment
					fragment = ""
				}
				break
			}
			continue
		}

		if !isLink {
			fragment += link + string(r)
			link = ""
		}

		// build fragment
		fragmentEnded := r == '.' || r == '\n' || r == ';' || i == lastindex
		
		if len(fragment)+len(candidate) > maxSize {
			fragmentEnded = true
		}
		
		if fragmentEnded {
			// commit to candidate
			if len(candidate) < minSize || len(candidate)+len(fragment) <= maxSize {
				candidate += fragment
				fragment = ""
				if len(candidate) >= maxSize && i != lastindex {
					continue
				}
				continue
			}
			// too big
			break
		}
	}

	// Second pass: backward search for a good break point
	lasti := len(candidate)
	// Only perform backward search if we haven't reached the end of the text
	if len(candidate) < len(text) {
		for i := range candidate {
			ri := len(candidate) - i - 1
			r := candidate[ri]

			if ri < len(isLinkRune) && isLinkRune[ri] {
				// We are in the middle of a link.
				// Find where the link starts and break BEFORE it.
				for ri > 0 && isLinkRune[ri-1] {
					ri--
				}
				lasti = ri
				break
			}
			
			if r == ' ' || r == '.' || r == '\n' || r == ';' {
				if ri < minSize {
					break
				}
				lasti = ri + 1
				break
			}
		}
	}

	return SanitizeUnicode(candidate[:lasti])
}

// isStartWithLink checks if the given text starts with "http://" or "https://".
func isStartWithLink(text string) bool {
	return strings.HasPrefix(text, "http://") || strings.HasPrefix(text, "https://")
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
	text = SanitizeUnicode(text)
	text = strings.TrimSpace(text)
	
	for len(text) > 0 {
		chunk := FirstChunk(text, chunkSize/2, chunkSize)
		if chunk == "" {
			_, size := utf8.DecodeRuneInString(text)
			text = text[size:]
			text = strings.TrimLeftFunc(text, unicode.IsSpace)
			continue
		}
		chunks = append(chunks, chunk)

		if len(chunk) >= len(text) {
			break
		}

		removeSize := chunkSize - chunkOverlap
		if removeSize <= 0 {
			removeSize = 1
		}
		
		removechunk := FirstChunk(text, removeSize/2, removeSize)
		if len(removechunk) > 0 {
			text = text[len(removechunk):]
		} else {
			_, size := utf8.DecodeRuneInString(text)
			text = text[size:]
		}
		text = strings.TrimLeftFunc(text, unicode.IsSpace)
	}
	return chunks
}

// Substring returns a substring of s from start (inclusive) to end (exclusive)
// based on rune count, not byte count.
func Substring(s string, start int, end int) string {
	if s == "" {
		return ""
	}
	
	var startIdx, endIdx int
	var i int
	foundStart := false
	for byteIdx := range s {
		if i == start {
			startIdx = byteIdx
			foundStart = true
		}
		if i == end {
			endIdx = byteIdx
			return s[startIdx:endIdx]
		}
		i++
	}
	if foundStart {
		return s[startIdx:]
	}
	return ""
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

// isInvalidRune reports whether a rune is invalid or unwanted for text display.
func isInvalidRune(r rune) bool {
	if r == utf8.RuneError {
		return true
	}
	if (r < 0x20 && r != '\n' && r != '\t') || (r >= 0x7f && r < 0xa0) {
		return true
	}
	if (r >= 0xFDD0 && r <= 0xFDEF) || (r&0xFFFF == 0xFFFE) || (r&0xFFFF == 0xFFFF) {
		return true
	}
	if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
		return true
	}
	return false
}

// SanitizeUnicode removes all invalid or unwanted Unicode runes.
func SanitizeUnicode(s string) string {
	needsSanitize := false
	for _, r := range s {
		if isInvalidRune(r) {
			needsSanitize = true
			break
		}
	}
	if !needsSanitize {
		return s
	}

	out := make([]rune, 0, len(s))
	for _, r := range s {
		if !isInvalidRune(r) {
			out = append(out, r)
		}
	}
	return string(out)
}
