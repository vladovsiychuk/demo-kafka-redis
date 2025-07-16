package kafka

import (
	"time"

	"github/vladovsiychuk/demo-kafka-redis-forwarder/internal/reporter"
)

type Event struct {
	Partner string         `json:"partner"`
	Data    *reporter.Data `json:"data"`
	SentAt  time.Time      `json:"sent_at"`
}
