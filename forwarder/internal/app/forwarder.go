package app

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"justtrack.io/tests/backend/forwarder/internal/kafka"
	"justtrack.io/tests/backend/forwarder/internal/reporter"
)

type Forwarder struct {
	Partners []string
	Reporter reporter.Client
	Producer kafka.Producer
	Ticker   Ticker
}

// NewRealTicker returns a Ticker implementation backed by time.Ticker
func NewRealTicker(d time.Duration) Ticker {
	return &RealTicker{ticker: time.NewTicker(d)}
}

type Ticker interface {
	C() <-chan time.Time
	Stop()
}

type RealTicker struct {
	ticker *time.Ticker
}

func (r *RealTicker) C() <-chan time.Time {
	return r.ticker.C
}

func (r *RealTicker) Stop() {
	r.ticker.Stop()
}

func (f *Forwarder) Run(ctx context.Context) {
	defer f.Ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-f.Ticker.C():
			date := t.Format("2006-01-02")
			for _, partner := range f.Partners {
				data, err := f.Reporter.Fetch(ctx, partner, date)
				if err != nil {
					logrus.WithError(err).WithField("partner", partner).Warn("failed to fetch")
					continue
				}
				if err := f.Producer.Send(ctx, partner, data); err != nil {
					logrus.WithError(err).WithField("partner", partner).Warn("failed to send to kafka")
				}
			}
		}
	}
}
