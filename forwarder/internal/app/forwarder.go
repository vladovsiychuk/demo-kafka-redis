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
	Interval time.Duration
}

func (f *Forwarder) Run(ctx context.Context) {
	ticker := time.NewTicker(f.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
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
