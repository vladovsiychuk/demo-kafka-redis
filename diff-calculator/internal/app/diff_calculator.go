package app

import (
	"context"
	"time"

	"github/vladovsiychuk/demo-kafkaredis-diff/internal/datastore"
	"github/vladovsiychuk/demo-kafkaredis-diff/internal/diff"
	"github/vladovsiychuk/demo-kafkaredis-diff/internal/kafka"

	"github.com/sirupsen/logrus"
)

type DiffCalculator struct {
	Consumer kafka.Consumer
	Producer kafka.Producer
	Store    datastore.Store
}

func (dc *DiffCalculator) Run(ctx context.Context) {
	defer dc.Producer.Close()
	defer dc.Store.Close()
	defer dc.Consumer.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			event, ok := dc.Consumer.Next(ctx)
			if !ok || event == nil {
				return
			}

			prev, err := dc.Store.Get(ctx, event.Partner, event.Data.Date)
			if err != nil {
				logrus.WithError(err).WithFields(logrus.Fields{
					"partner": event.Partner,
					"date":    event.Data.Date,
				}).Warn("failed to get previous state")
				continue
			}

			curr := &datastore.State{
				Clicks:      event.Data.Clicks,
				Cost:        event.Data.Cost,
				Date:        event.Data.Date,
				Impressions: event.Data.Impressions,
				Installs:    event.Data.Installs,
			}

			if prev != nil {
				diffFields := diff.Calculate(*prev, *curr)
				if diffFields.HasChanges() {
					diffEvent := &kafka.DiffEvent{
						Partner:   event.Partner,
						Date:      event.Data.Date,
						Diffs:     diffFields,
						Timestamp: time.Now().UTC(),
					}
					if err := dc.Producer.Send(ctx, diffEvent); err != nil {
						logrus.WithError(err).WithFields(logrus.Fields{
							"partner": event.Partner,
							"date":    event.Data.Date,
						}).Warn("failed to send diff to kafka")
					}
				}
			}

			if err := dc.Store.Set(ctx, event.Partner, event.Data.Date, curr); err != nil {
				logrus.WithError(err).WithFields(logrus.Fields{
					"partner": event.Partner,
					"date":    event.Data.Date,
				}).Warn("failed to update state in store")
			}
		}
	}
}
