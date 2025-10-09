package chunker

import (
	"fmt"
	"os"
	"testing"
)

func TestSample(t *testing.T) {
	data, err := os.ReadFile("./sample.md")
	if err != nil {
		t.Fatal(err)
	}

	chunks := Chunk(string(data), 800, 100)
	for _, chunk := range chunks {
		fmt.Println("----------\n", chunk)
	}
}

func TestShort(t *testing.T) {
	text := "word001 word002 word003"
	chunks := Chunk(text, 5, 1)
	expect := []string{
		"word00",
		"01 ",
		"word00",
		"02 ",
		"word00",
		"03",
	}

	if len(chunks) != len(expect) {
		t.Fatalf("expected %d chunks, got %d", len(expect), len(chunks))
	}

	for i, chunk := range chunks {
		if chunk != expect[i] {
			t.Errorf("expected chunk %d to be '%s', got '%s'", i, expect[i], chunk)
		}
	}
}

func TestSimple2(t *testing.T) {
	text := "word9 word10 word11 word12"
	chunks := Chunk(text, 14, 4)
	expect := []string{
		"word9 word10 ",
		"word10 word11 ",
		"word11 word12",
	}

	if len(chunks) != len(expect) {
		t.Fatalf("expected %d chunks, got %d", len(expect), len(chunks))
	}

	for i, chunk := range chunks {
		if chunk != expect[i] {
			t.Errorf("expected chunk %d to be '%s', got '%s'", i, expect[i], chunk)
		}
	}
}

func TestSimple(t *testing.T) {
	text := "word1 word2 word3 word4 word5 word6 word7 word8 word9 word10 word11 word12"
	chunks := Chunk(text, 20, 6)

	expect := []string{
		"word1 word2 word3 ",
		"word3 word4 word5 ",
		"word5 word6 word7 ",
		"word7 word8 word9 ",
		"word9 word10 word11 ",
		"word11 word12",
	}

	if len(chunks) != len(expect) {
		t.Fatalf("expected %d chunks, got %d", len(expect), len(chunks))
	}

	for i, chunk := range chunks {
		if chunk != expect[i] {
			t.Errorf("expected chunk %d to be '%s', got '%s'", i, expect[i], chunk)
		}
	}
}

func TestLongWords(t *testing.T) {
	text := "👋 Lopadotemachoselachogaleokranioleipsanodrimhypotrimmatosilphioparaomelitokatakechymenokichlepikossyphophattoperisteralektryonoptekephalliokigklopeleiolagoiosiraiobaphetraganopterygon 👋"
	expect := []string{
		"👋 Lopadotemachoselachogaleokranioleipsanodrimhyp",
		"anodrimhypotrimmatosilphioparaomelitokatakechymenok",
		"kechymenokichlepikossyphophattoperisteralektryonopt",
		"ektryonoptekephalliokigklopeleiolagoiosiraiobaphetr",
		"aiobaphetraganopterygon 👋",
	}

	chunks := Chunk(text, 50, 10)
	if len(chunks) != len(expect) {
		t.Fatalf("expected %d chunks, got %d", len(expect), len(chunks))
	}

	for i, chunk := range chunks {
		if chunk != expect[i] {
			t.Errorf("expected chunk %d to be '%s', got '%s'", i, expect[i], chunk)
		}
	}
}

func TestChunkMarkdown(t *testing.T) {
	text := `# Sample Product Overview

Welcome to the **Amazing Gadget**! This device helps you stay productive while looking stylish.
Check out the official website for more details: [Amazing Gadget](https://example.com/amazing-gadget).

![Amazing Gadget Photo](https://via.placeholder.com/600x300 "Amazing Gadget in Action")

## Key Features
- **Lightweight & Portable** – only 500 g.
- **Battery Life** – up to 12 hours on a single charge.
- **Connectivity** – Wi-Fi 6 and Bluetooth 5.2 support.
- **Warranty** – 2-year international coverage.

> *Tip:* Combine it with our [accessories pack](https://example.com/accessories) for the best experience.`

	expect := []string{
		`# Sample Product Overview

Welcome to the **Amazing Gadget**! This device helps you stay productive `,
		`you stay productive while looking stylish.
Check out the official website for more details: [Amazing `,
		`Check out the official website for more details: [Amazing `,
		`Gadget](https://example.com/amazing-gadget).

![Amazing Gadget Photo](https://via.placeholder.com/600x300 `,
		`![Amazing Gadget Photo](https://via.placeholder.com/600x300 "Amazing Gadget in Action")

`, `Action")

## Key Features
- **Lightweight & Portable** – only 500 g.
`,
		`- **Battery Life** – up to 12 hours on a single charge.
`,
		`- **Connectivity** – Wi-Fi 6 and Bluetooth 5.2 support.
`, `- **Warranty** – 2-year international coverage.

`,
		`> *Tip:* Combine it with our [accessories pack](https://example.com/accessories) for the best `,
		`for the best experience.`,
	}
	chunks := Chunk(text, 100, 20)
	if len(chunks) != len(expect) {
		t.Fatalf("expected %d chunks, got %d", len(expect), len(chunks))
	}

	for i, chunk := range chunks {
		if chunk != expect[i] {
			t.Errorf("expected chunk %d to be '%s', got '%s'", i, expect[i], chunk)
		}
	}
}

func TestFirstChunk(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		min      int
		max      int
		expected string
	}{
		{
			name:     "simple sentence",
			text:     "This is a simple sentence. This is another sentence.",
			min:      20,
			max:      40,
			expected: "This is a simple sentence.",
		},
		{
			name:     "text shorter than max",
			text:     "This is a short sentence.",
			min:      20,
			max:      40,
			expected: "This is a short sentence.",
		},
		{
			name: "no sentence break in range",
			text: "This is a very long sentence that does not have any sentence breaks within the specified min and max range.",
			min:  20,
			max:  40,

			expected: "This is a very long sentence that does ",
		},
		{
			name:     "first sentence longer than max",
			text:     "This is a very long first sentence that is longer than the max length. This is a second sentence.",
			min:      20,
			max:      40,
			expected: "This is a very long first sentence that ",
		},
		{
			name:     "multiple paragraphs",
			text:     "This is the first paragraph.\n\nThis is the second paragraph.",
			min:      20,
			max:      50,
			expected: "This is the first paragraph.\n\n",
		},
		{
			name:     "multiple paragraphs with links",
			text:     "This is (https://google.com/a/b/c).\n\nThis is the second paragraph.",
			min:      20,
			max:      50,
			expected: "This is (https://google.com/a/b/c).\n\n",
		},
		{
			name:     "multiple paragraphs with links 2",
			text:     "This is (https://google.com/a/b/c).\n\nThis is the second paragraph.",
			min:      5,
			max:      10,
			expected: "This is ",
		},
		{
			name:     "first character is space",
			text:     " This is this a paragraph",
			min:      5,
			max:      10,
			expected: " This is ",
		},
		{
			name:     "first character is dot",
			text:     ".This is this a paragraph",
			min:      5,
			max:      10,
			expected: ".This is ",
		},
		{
			name:     "invalid unicode",
			text:      "Điều này có thể là nguyên nhân khiến mainboard của máy \u0017bị hư hỏng nặng và khả năng cao phải thay thế để có thể tiếp tục sử dụng.Việc sạc pin với nguồn điện không ổn định hay không đủ năng suất có thẻ khiến cho mainboard của thiết bị nói tr\u0017",
			min:     5,
			max:      1000,
			expected: "Điều này có thể là nguyên nhân khiến mainboard của máy bị hư hỏng nặng và khả năng cao phải thay thế để có thể tiếp tục sử dụng.Việc sạc pin với nguồn điện không ổn định hay không đủ năng suất có thẻ khiến cho mainboard của thiết bị nói tr",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//			if tc.name != "no sentence break in range" {
			//return
			//}
			chunk := FirstChunk(tc.text, tc.min, tc.max)
			if chunk != tc.expected {
				t.Errorf("expected '%s', got '%s'", tc.expected, chunk)
			}
		})
	}
}
