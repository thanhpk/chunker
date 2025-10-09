# Chunker

Break text into small chunks

## Installation

To use this library in your Go project, you can use `go get`:

```bash
go get github.com/thanhpk/chunker
```

## Usage

The `chunker` library provides a single function, `Chunk`, which takes a string of text, a chunk size, and a chunk overlap as input, and returns a slice of strings.

```go
package main

import (
	"fmt"
	"github.com/thanhpk/chunker"
)

func main() {
	text := "word1 word2 word3 word4 word5 word6 word7 word8 word9 word10 word11 word12"
	chunkSize := 20
	chunkOverlap := 5 // Overlap is in number of words

	chunks := chunker.Chunk(text, chunkSize, chunkOverlap)
	for i, chunk := range chunks {
		fmt.Printf("Chunk %d: %s", i+1, chunk)
	}
}
```

### Output

```
Chunk 1: word1 word2 word3
Chunk 2: word3 word4 word5
Chunk 3: word5 word6 word7
Chunk 4: word7 word8 word9
Chunk 5: word9 word10 word11
Chunk 6: word11 word12
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
