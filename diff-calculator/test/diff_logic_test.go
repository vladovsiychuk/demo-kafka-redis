package test

import (
	"testing"

	"github/vladovsiychuk/demo-kafkaredis-diff/internal/datastore"
	"github/vladovsiychuk/demo-kafkaredis-diff/internal/diff"
)

func TestCalculateDiff(t *testing.T) {
	cases := []struct {
		name      string
		prev      datastore.State
		curr      datastore.State
		wantDiffs map[string]bool // which fields are expected to change
	}{
		{
			name: "all fields same",
			prev: datastore.State{Clicks: 1, Cost: 2.0, Impressions: 3, Installs: 4},
			curr: datastore.State{Clicks: 1, Cost: 2.0, Impressions: 3, Installs: 4},
			wantDiffs: map[string]bool{
				"Clicks": false, "Cost": false, "Impressions": false, "Installs": false,
			},
		},
		{
			name: "one field changed",
			prev: datastore.State{Clicks: 1, Cost: 2.0, Impressions: 3, Installs: 4},
			curr: datastore.State{Clicks: 2, Cost: 2.0, Impressions: 3, Installs: 4},
			wantDiffs: map[string]bool{
				"Clicks": true, "Cost": false, "Impressions": false, "Installs": false,
			},
		},
		{
			name: "multiple fields changed",
			prev: datastore.State{Clicks: 1, Cost: 2.0, Impressions: 3, Installs: 4},
			curr: datastore.State{Clicks: 2, Cost: 1.5, Impressions: 5, Installs: 4},
			wantDiffs: map[string]bool{
				"Clicks": true, "Cost": true, "Impressions": true, "Installs": false,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			diffs := diff.Calculate(c.prev, c.curr)
			// Check each field
			if (diffs.Clicks != nil) != c.wantDiffs["Clicks"] {
				t.Errorf("Clicks diff: got %v, want %v", diffs.Clicks != nil, c.wantDiffs["Clicks"])
			}
			if (diffs.Cost != nil) != c.wantDiffs["Cost"] {
				t.Errorf("Cost diff: got %v, want %v", diffs.Cost != nil, c.wantDiffs["Cost"])
			}
			if (diffs.Impressions != nil) != c.wantDiffs["Impressions"] {
				t.Errorf("Impressions diff: got %v, want %v", diffs.Impressions != nil, c.wantDiffs["Impressions"])
			}
			if (diffs.Installs != nil) != c.wantDiffs["Installs"] {
				t.Errorf("Installs diff: got %v, want %v", diffs.Installs != nil, c.wantDiffs["Installs"])
			}

			hasChanges := c.wantDiffs["Clicks"] || c.wantDiffs["Cost"] || c.wantDiffs["Impressions"] || c.wantDiffs["Installs"]
			if diffs.HasChanges() != hasChanges {
				t.Errorf("HasChanges: got %v, want %v", diffs.HasChanges(), hasChanges)
			}
		})
	}
}
