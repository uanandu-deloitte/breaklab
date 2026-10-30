package variants

import (
	"strings"
	"testing"

	"breaklab/internal/catalog"
)

func TestPickWeightedDistribution(t *testing.T) {
	variantsList := []catalog.Variant{
		{ID: "a", Weight: 1},
		{ID: "b", Weight: 2},
		{ID: "c", Weight: 3},
	}
	const n = 10000
	counts := map[string]int{}
	for i := 0; i < n; i++ {
		v, err := Pick(variantsList, "", nil)
		if err != nil {
			t.Fatal(err)
		}
		counts[v.ID]++
	}

	// Observed share should land within 20% (relative) of weight/total for
	// each variant; 10k draws makes that comfortably stable.
	total := 6.0
	for _, v := range variantsList {
		want := float64(v.Weight) / total
		got := float64(counts[v.ID]) / float64(n)
		if diff := got - want; diff < -0.2*want-0.02 || diff > 0.2*want+0.02 {
			t.Errorf("variant %q: observed share %.3f, want ~%.3f", v.ID, got, want)
		}
	}
}

func TestPickNoImmediateRepeat(t *testing.T) {
	variantsList := []catalog.Variant{
		{ID: "a", Weight: 1},
		{ID: "b", Weight: 1},
		{ID: "c", Weight: 1},
	}
	for _, previous := range []string{"a", "b", "c"} {
		for i := 0; i < 1000; i++ {
			v, err := Pick(variantsList, previous, nil)
			if err != nil {
				t.Fatal(err)
			}
			if v.ID == previous {
				t.Fatalf("picked the previous variant %q", previous)
			}
		}
	}
}

func TestPickNoRepeatSkippedForOneVariant(t *testing.T) {
	variantsList := []catalog.Variant{{ID: "only", Weight: 1}}
	v, err := Pick(variantsList, "only", nil)
	if err != nil || v.ID != "only" {
		t.Fatalf("v = %+v, err = %v", v, err)
	}
}

func TestPickStalePreviousIgnored(t *testing.T) {
	variantsList := []catalog.Variant{{ID: "a", Weight: 1}, {ID: "b", Weight: 1}}
	// "previous" no longer exists in the lab's variant list; nothing should
	// be excluded, so both variants remain reachable.
	seen := map[string]bool{}
	for i := 0; i < 200; i++ {
		v, err := Pick(variantsList, "removed-variant", nil)
		if err != nil {
			t.Fatal(err)
		}
		seen[v.ID] = true
	}
	if !seen["a"] || !seen["b"] {
		t.Errorf("seen = %v, want both a and b reachable", seen)
	}
}

func TestPickWeightZeroFallback(t *testing.T) {
	variantsList := []catalog.Variant{
		{ID: "previous", Weight: 1},
		{ID: "disabled-one", Weight: 0},
		{ID: "disabled-two", Weight: 0},
	}
	seen := map[string]bool{}
	for i := 0; i < 200; i++ {
		v, err := Pick(variantsList, "previous", nil)
		if err != nil {
			t.Fatal(err)
		}
		if v.ID == "previous" {
			t.Fatalf("picked excluded previous variant")
		}
		seen[v.ID] = true
	}
	if !seen["disabled-one"] || !seen["disabled-two"] {
		t.Errorf("seen = %v, want a uniform pick over both weight-0 candidates", seen)
	}
}

func TestPickSeedIsReproducible(t *testing.T) {
	variantsList := []catalog.Variant{
		{ID: "a", Weight: 1},
		{ID: "b", Weight: 1},
		{ID: "c", Weight: 1},
	}
	seed := int64(42)

	v1, err := Pick(variantsList, "", &seed)
	if err != nil {
		t.Fatal(err)
	}
	v2, err := Pick(variantsList, "", &seed)
	if err != nil {
		t.Fatal(err)
	}
	if v1.ID != v2.ID {
		t.Fatalf("same seed picked %q then %q", v1.ID, v2.ID)
	}

	// Simulate two successive `start --seed 42` calls where the first call's
	// pick has already been recorded as "previous" by the time the second
	// one runs: the seed must still bypass that exclusion (design.md
	// Decision 4) so the two calls agree.
	v3, err := Pick(variantsList, v1.ID, &seed)
	if err != nil {
		t.Fatal(err)
	}
	if v3.ID != v1.ID {
		t.Fatalf("seeded pick changed when previous was set: got %q, want %q", v3.ID, v1.ID)
	}
}

func TestPickSeedDifferentSeedsCanDiffer(t *testing.T) {
	variantsList := []catalog.Variant{
		{ID: "a", Weight: 1},
		{ID: "b", Weight: 1},
		{ID: "c", Weight: 1},
		{ID: "d", Weight: 1},
		{ID: "e", Weight: 1},
	}
	seen := map[string]bool{}
	for seed := int64(0); seed < 50; seed++ {
		v, err := Pick(variantsList, "", &seed)
		if err != nil {
			t.Fatal(err)
		}
		seen[v.ID] = true
	}
	if len(seen) < 2 {
		t.Errorf("50 different seeds all picked the same variant: %v", seen)
	}
}

func TestPickRejectsNegativeWeight(t *testing.T) {
	variantsList := []catalog.Variant{{ID: "a", Weight: -1}, {ID: "b", Weight: 1}}
	if _, err := Pick(variantsList, "", nil); err == nil || !strings.Contains(err.Error(), "negative weight") {
		t.Fatalf("err = %v", err)
	}
}

func TestPickRejectsEmptyList(t *testing.T) {
	if _, err := Pick(nil, "", nil); err == nil {
		t.Fatal("expected error for empty variant list")
	}
}

func TestForce(t *testing.T) {
	variantsList := []catalog.Variant{{ID: "a", Weight: 1}, {ID: "b", Weight: 1}}

	v, err := Force(variantsList, "b")
	if err != nil || v.ID != "b" {
		t.Fatalf("v = %+v, err = %v", v, err)
	}

	if _, err := Force(variantsList, "nope"); err == nil || !strings.Contains(err.Error(), "nope") {
		t.Fatalf("err = %v", err)
	}
}
