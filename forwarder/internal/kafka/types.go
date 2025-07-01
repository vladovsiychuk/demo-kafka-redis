package kafka

import (
	"time"

	"justtrack.io/tests/backend/forwarder/internal/reporter"
)

type Event struct {
	Partner string         `json:"partner"`
	Data    *reporter.Data `json:"data"`
	SentAt  time.Time      `json:"sent_at"`
}
