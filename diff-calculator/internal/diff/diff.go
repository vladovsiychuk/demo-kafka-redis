package diff

import (
	"justtrack.io/tests/backend/diff-calculator/internal/datastore"
	"justtrack.io/tests/backend/diff-calculator/internal/kafka"
)

// Calculate compares previous and current state, and returns FieldDiffs with only changed fields.
func Calculate(prev, curr datastore.State) kafka.FieldDiffs {
	diffs := kafka.FieldDiffs{}

	if prev.Clicks != curr.Clicks {
		diffs.Clicks = &kafka.FieldChange{Before: prev.Clicks, After: curr.Clicks}
	}
	if prev.Cost != curr.Cost {
		diffs.Cost = &kafka.FieldChange{Before: prev.Cost, After: curr.Cost}
	}
	if prev.Impressions != curr.Impressions {
		diffs.Impressions = &kafka.FieldChange{Before: prev.Impressions, After: curr.Impressions}
	}
	if prev.Installs != curr.Installs {
		diffs.Installs = &kafka.FieldChange{Before: prev.Installs, After: curr.Installs}
	}
	return diffs
}
