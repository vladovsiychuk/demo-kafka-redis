package test

import (
	"encoding/json"
	"testing"

	"github/vladovsiychuk/demo-kafka-redis-forwarder/internal/reporter"
)

func TestIntOrString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		input   string
		want    reporter.IntOrString
		wantErr bool
	}{
		{`1`, 1, false},
		{`"2"`, 2, false},
		{`" 3 "`, 3, false},
		{`""`, 0, false},
		{`"abc"`, 0, true},
		{`null`, 0, false},
	}
	for _, tt := range tests {
		var out reporter.IntOrString
		err := json.Unmarshal([]byte(tt.input), &out)
		if (err != nil) != tt.wantErr {
			t.Errorf("Unmarshal(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
		if !tt.wantErr && out != tt.want {
			t.Errorf("Unmarshal(%q) = %v, want %v", tt.input, out, tt.want)
		}
	}
}

func TestFloatOrString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		input   string
		want    reporter.FloatOrString
		wantErr bool
	}{
		{`1.23`, 1.23, false},
		{`"2.34"`, 2.34, false},
		{`" 3.45 "`, 3.45, false},
		{`""`, 0, false},
		{`"abc"`, 0, true},
		{`null`, 0, false},
	}
	for _, tt := range tests {
		var out reporter.FloatOrString
		err := json.Unmarshal([]byte(tt.input), &out)
		if (err != nil) != tt.wantErr {
			t.Errorf("Unmarshal(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
		if !tt.wantErr && (float64(out) != float64(tt.want)) {
			t.Errorf("Unmarshal(%q) = %v, want %v", tt.input, out, tt.want)
		}
	}
}
