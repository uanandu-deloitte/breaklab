// Package variants implements weighted random selection of a lab's
// breakage variant, with a no-immediate-repeat rule, reproducible picks
// via a seed, and a forced override for tests.
package variants

import (
	"fmt"
	"math/rand/v2"

	"breaklab/internal/catalog"
)

// Pick selects a variant from variantsList.
//
// When seed is nil, the variant named by previous is excluded before a
// weighted random draw, provided previous is non-empty and variantsList has
// two or more entries (the no-immediate-repeat rule). When seed is not nil,
// the draw is deterministic for that seed and runs over the full list,
// ignoring previous: see design.md's "Decisions" section (change
// add-random-variants) for why reproducibility takes priority over
// exclusion in that case.
func Pick(variantsList []catalog.Variant, previous string, seed *int64) (catalog.Variant, error) {
	if err := validateWeights(variantsList); err != nil {
		return catalog.Variant{}, err
	}

	candidates := variantsList
	if seed == nil {
		candidates = excludePrevious(variantsList, previous)
	}

	if seed != nil {
		s := uint64(*seed)
		rng := rand.New(rand.NewPCG(s, s))
		return weightedPick(candidates, rng.IntN)
	}
	return weightedPick(candidates, rand.IntN)
}

// Force returns the variant in variantsList whose id matches forcedID,
// bypassing both exclusion and weighting, for BREAKLAB_FORCE_VARIANT.
func Force(variantsList []catalog.Variant, forcedID string) (catalog.Variant, error) {
	for _, v := range variantsList {
		if v.ID == forcedID {
			return v, nil
		}
	}
	return catalog.Variant{}, fmt.Errorf("BREAKLAB_FORCE_VARIANT %q is not a variant of this lab", forcedID)
}

func validateWeights(variantsList []catalog.Variant) error {
	if len(variantsList) == 0 {
		return fmt.Errorf("no variants to choose from")
	}
	for _, v := range variantsList {
		if v.Weight < 0 {
			return fmt.Errorf("variant %q has a negative weight (%d)", v.ID, v.Weight)
		}
	}
	return nil
}

// excludePrevious drops the variant named by previous from variantsList,
// unless previous is empty, there is only one variant, or previous does not
// name any variant in the list (e.g. stale history from a lab.yaml that has
// since changed) — in all of those cases nothing is excluded.
func excludePrevious(variantsList []catalog.Variant, previous string) []catalog.Variant {
	if previous == "" || len(variantsList) < 2 {
		return variantsList
	}
	filtered := make([]catalog.Variant, 0, len(variantsList)-1)
	for _, v := range variantsList {
		if v.ID != previous {
			filtered = append(filtered, v)
		}
	}
	if len(filtered) == 0 {
		return variantsList
	}
	return filtered
}

// weightedPick draws one variant from candidates using intN(n) to produce a
// uniform value in [0, n). If every candidate has weight 0, it falls back
// to a uniform pick among them instead of never selecting anything.
func weightedPick(candidates []catalog.Variant, intN func(int) int) (catalog.Variant, error) {
	total := 0
	for _, v := range candidates {
		total += v.Weight
	}
	if total == 0 {
		return candidates[intN(len(candidates))], nil
	}

	draw := intN(total)
	for _, v := range candidates {
		if draw < v.Weight {
			return v, nil
		}
		draw -= v.Weight
	}
	// Unreachable: draw < total is guaranteed by intN's contract.
	return candidates[len(candidates)-1], nil
}
