package postmortem

import "testing"

func TestSplitBlocksNoFence(t *testing.T) {
	blocks := SplitBlocks("First paragraph.\n\nSecond paragraph.")
	if len(blocks) != 2 {
		t.Fatalf("got %d blocks, want 2: %+v", len(blocks), blocks)
	}
	if blocks[0].Kind != BlockParagraph || blocks[0].Text != "First paragraph." {
		t.Errorf("blocks[0] = %+v", blocks[0])
	}
	if blocks[1].Kind != BlockParagraph || blocks[1].Text != "Second paragraph." {
		t.Errorf("blocks[1] = %+v", blocks[1])
	}
}

func TestSplitBlocksPlainFence(t *testing.T) {
	text := "Before.\n\n```\nsome code\n```\n\nAfter."
	blocks := SplitBlocks(text)
	if len(blocks) != 3 {
		t.Fatalf("got %d blocks, want 3: %+v", len(blocks), blocks)
	}
	if blocks[0].Kind != BlockParagraph || blocks[0].Text != "Before." {
		t.Errorf("blocks[0] = %+v", blocks[0])
	}
	if blocks[1].Kind != BlockCode || blocks[1].Text != "some code" {
		t.Errorf("blocks[1] = %+v", blocks[1])
	}
	if blocks[2].Kind != BlockParagraph || blocks[2].Text != "After." {
		t.Errorf("blocks[2] = %+v", blocks[2])
	}
}

func TestSplitBlocksDiffFence(t *testing.T) {
	text := "```diff\n- FROM openjdk:8u121-jre\n+ FROM eclipse-temurin:17-jre-alpine\n unrelated context line\n```"
	blocks := SplitBlocks(text)
	if len(blocks) != 1 {
		t.Fatalf("got %d blocks, want 1: %+v", len(blocks), blocks)
	}
	b := blocks[0]
	if b.Kind != BlockDiff {
		t.Fatalf("Kind = %v, want BlockDiff", b.Kind)
	}
	if len(b.Diff) != 3 {
		t.Fatalf("got %d diff lines, want 3: %+v", len(b.Diff), b.Diff)
	}
	if b.Diff[0].Sign != "-" || b.Diff[0].Text != " FROM openjdk:8u121-jre" {
		t.Errorf("Diff[0] = %+v", b.Diff[0])
	}
	if b.Diff[1].Sign != "+" || b.Diff[1].Text != " FROM eclipse-temurin:17-jre-alpine" {
		t.Errorf("Diff[1] = %+v", b.Diff[1])
	}
	if b.Diff[2].Sign != "" {
		t.Errorf("Diff[2].Sign = %q, want empty (context line)", b.Diff[2].Sign)
	}
}

func TestSplitBlocksUnterminatedFence(t *testing.T) {
	text := "Before.\n\n```\nline one\nline two"
	blocks := SplitBlocks(text)
	if len(blocks) != 2 {
		t.Fatalf("got %d blocks, want 2: %+v", len(blocks), blocks)
	}
	if blocks[1].Kind != BlockCode || blocks[1].Text != "line one\nline two" {
		t.Errorf("blocks[1] = %+v", blocks[1])
	}
}

func TestSplitBlocksEmpty(t *testing.T) {
	if blocks := SplitBlocks(""); len(blocks) != 0 {
		t.Errorf("got %d blocks, want 0: %+v", len(blocks), blocks)
	}
}
