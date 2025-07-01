package kafka

import (
	"time"
)

type Event struct {
	Partner string    `json:"partner"`
	Data    EventData `json:"data"`
	SentAt  time.Time `json:"sent_at"`
}

type EventData struct {
	Clicks      int     `json:"clicks"`
	Cost        float64 `json:"cost"`
	Date        string  `json:"date"`
	Impressions int     `json:"impressions"`
	Installs    int     `json:"installs"`
}

type DiffEvent struct {
	Partner   string     `json:"partner"`
	Date      string     `json:"date"`
	Diffs     FieldDiffs `json:"diffs"`
	Timestamp time.Time  `json:"timestamp"`
}

type FieldDiffs struct {
	Clicks      *FieldChange `json:"clicks,omitempty"`
	Cost        *FieldChange `json:"cost,omitempty"`
	Impressions *FieldChange `json:"impressions,omitempty"`
	Installs    *FieldChange `json:"installs,omitempty"`
}

type FieldChange struct {
	Before interface{} `json:"before"`
	After  interface{} `json:"after"`
}

func (d FieldDiffs) HasChanges() bool {
	return d.Clicks != nil || d.Cost != nil || d.Impressions != nil || d.Installs != nil
}
