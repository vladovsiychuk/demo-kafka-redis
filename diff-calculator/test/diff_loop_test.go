package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"justtrack.io/tests/backend/diff-calculator/internal/app"
	"justtrack.io/tests/backend/diff-calculator/internal/datastore"
	"justtrack.io/tests/backend/diff-calculator/internal/kafka"
	"justtrack.io/tests/backend/diff-calculator/mocks"
)

// Arrange helpers for the test
func makeEvent(clicks, impressions int) *kafka.Event {
	return &kafka.Event{
		Partner: "pixelAds",
		Data: kafka.EventData{
			Clicks:      clicks,
			Cost:        1.23,
			Date:        "2025-07-01",
			Impressions: impressions,
			Installs:    7,
		},
		SentAt: time.Now(),
	}
}

func TestDiffCalculator_Run_FullLoop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockConsumer := mocks.NewConsumer(t)
	mockProducer := mocks.NewProducer(t)
	mockStore := mocks.NewStore(t)

	// The first event: no previous state, so no diff emitted
	firstEvent := makeEvent(1, 100)
	mockConsumer.On("Next", mock.Anything).Return(firstEvent, true).Once()
	mockStore.On("Get", mock.Anything, "pixelAds", "2025-07-01").Return(nil, nil).Once()
	mockStore.On("Set", mock.Anything, "pixelAds", "2025-07-01", mock.Anything).Return(nil).Once()

	// The second event: clicks changed, should emit a diff
	secondEvent := makeEvent(2, 100)
	prevState := &datastore.State{
		Clicks:      1,
		Cost:        1.23,
		Date:        "2025-07-01",
		Impressions: 100,
		Installs:    7,
	}
	mockConsumer.On("Next", mock.Anything).Return(secondEvent, true).Once()
	mockStore.On("Get", mock.Anything, "pixelAds", "2025-07-01").Return(prevState, nil).Once()
	// We expect a diff to be produced:
	mockProducer.On("Send", mock.Anything, mock.MatchedBy(func(diffEvent *kafka.DiffEvent) bool {
		return diffEvent.Partner == "pixelAds" && diffEvent.Diffs.Clicks != nil && diffEvent.Diffs.Clicks.After == 2
	})).Return(nil).Once()
	mockStore.On("Set", mock.Anything, "pixelAds", "2025-07-01", mock.Anything).Return(nil).Once()

	// After two events, simulate shutdown by returning (nil, false)
	mockConsumer.On("Next", mock.Anything).Return(nil, false).Once()

	// Mocks for Close
	mockProducer.On("Close").Return(nil).Once()
	mockStore.On("Close").Return(nil).Once()
	mockConsumer.On("Close").Return(nil).Once()

	// Create and run the app
	calc := app.DiffCalculator{
		Consumer: mockConsumer,
		Producer: mockProducer,
		Store:    mockStore,
	}
	calc.Run(ctx)

	// Assert everything was called as expected
	mockConsumer.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
	mockStore.AssertExpectations(t)
	// Check Send was only called for the second event
	mockProducer.AssertNumberOfCalls(t, "Send", 1)
	// Also: check the Set was called both times (always update state)
	mockStore.AssertNumberOfCalls(t, "Set", 2)
	// And Next was called three times (2 events + shutdown)
	mockConsumer.AssertNumberOfCalls(t, "Next", 3)
}
