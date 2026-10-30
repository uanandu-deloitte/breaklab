package postmortem

import "testing"

func TestMeaningfulWordsDropsStopwordsAndShortWords(t *testing.T) {
	got := MeaningfulWords("The cache was stale after a rotation, and this broke auth.")
	want := map[string]bool{"cache": true, "stale": true, "rotation": true, "broke": true, "auth": true}
	if len(got) != len(want) {
		t.Fatalf("got %v, want exactly %v", got, want)
	}
	for _, w := range got {
		if !want[w] {
			t.Errorf("unexpected word %q in %v", w, got)
		}
	}
}

func TestRootCauseMentionsVariantOverlapFound(t *testing.T) {
	hints := []string{"the cache was never invalidated after rotation"}
	rootCause := "The real problem was a stale cache that never got invalidated."
	if !RootCauseMentionsVariant(rootCause, hints) {
		t.Error("expected overlap on 'cache'/'invalidated'/'stale'")
	}
}

func TestRootCauseMentionsVariantNoOverlap(t *testing.T) {
	hints := []string{"the cache was never invalidated after rotation"}
	rootCause := "A trailing semicolon was missing from the nginx config."
	if RootCauseMentionsVariant(rootCause, hints) {
		t.Error("expected no overlap between unrelated sentences")
	}
}

func TestRootCauseMentionsVariantEmptyHints(t *testing.T) {
	if RootCauseMentionsVariant("anything at all here", nil) {
		t.Error("expected no overlap when there are no hints to compare against")
	}
}
