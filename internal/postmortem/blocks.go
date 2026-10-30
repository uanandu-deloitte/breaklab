package postmortem

import "strings"

// BlockKind distinguishes a postmortem segment's rendering treatment.
type BlockKind int

const (
	BlockParagraph BlockKind = iota
	BlockCode
	BlockDiff
)

// DiffLine is one line of a ```diff fenced block, split out so the
// template can color added/removed lines without needing its own
// string-splitting logic.
type DiffLine struct {
	Sign string // "+", "-", or "" for a context line
	Text string // the line with its leading +/- (if any) stripped
}

// Block is one rendered segment of a postmortem section: a paragraph of
// prose, a plain fenced code block, or a ```diff fenced block.
type Block struct {
	Kind BlockKind
	Text string     // BlockParagraph and BlockCode: the rendered text
	Diff []DiffLine // BlockDiff only
}

// IsParagraph, IsCode, and IsDiff let templates branch on a Block's kind
// by name instead of comparing against BlockKind's numeric iota values.
func (b Block) IsParagraph() bool { return b.Kind == BlockParagraph }
func (b Block) IsCode() bool      { return b.Kind == BlockCode }
func (b Block) IsDiff() bool      { return b.Kind == BlockDiff }

// SplitBlocks walks text and returns an ordered list of paragraph and
// fenced-code-block segments, the same paragraph boundaries
// splitParagraphs already used (blank-line-separated), interleaved with
// any fenced code. A fence tagged "diff" (```diff) produces a
// BlockDiff; any other fence (including untagged) produces BlockCode.
// An unterminated fence runs to the end of text. This is deliberately
// not general markdown — postmortems are five fixed sections of prose
// plus optional fenced code, and that's all this renders.
func SplitBlocks(text string) []Block {
	var blocks []Block
	lines := strings.Split(text, "\n")

	var prose []string
	flushProse := func() {
		joined := strings.Join(prose, "\n")
		for _, p := range strings.Split(joined, "\n\n") {
			if p = strings.TrimSpace(p); p != "" {
				blocks = append(blocks, Block{Kind: BlockParagraph, Text: p})
			}
		}
		prose = nil
	}

	for i := 0; i < len(lines); i++ {
		tag, ok := fenceTag(lines[i])
		if !ok {
			prose = append(prose, lines[i])
			continue
		}
		flushProse()
		var code []string
		for i++; i < len(lines) && !isFence(lines[i]); i++ {
			code = append(code, lines[i])
		}
		codeText := strings.Join(code, "\n")
		if tag == "diff" {
			blocks = append(blocks, Block{Kind: BlockDiff, Text: codeText, Diff: splitDiffLines(code)})
		} else {
			blocks = append(blocks, Block{Kind: BlockCode, Text: codeText})
		}
	}
	flushProse()
	return blocks
}

func fenceTag(line string) (tag string, ok bool) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "```") {
		return "", false
	}
	return strings.TrimSpace(strings.TrimPrefix(trimmed, "```")), true
}

func isFence(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "```")
}

func splitDiffLines(lines []string) []DiffLine {
	out := make([]DiffLine, 0, len(lines))
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "+"):
			out = append(out, DiffLine{Sign: "+", Text: strings.TrimPrefix(line, "+")})
		case strings.HasPrefix(line, "-"):
			out = append(out, DiffLine{Sign: "-", Text: strings.TrimPrefix(line, "-")})
		default:
			out = append(out, DiffLine{Text: line})
		}
	}
	return out
}
