package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"justtrack.io/tests/backend/forwarder/internal/app"
	"justtrack.io/tests/backend/forwarder/internal/reporter"
	"justtrack.io/tests/backend/forwarder/mocks"
)

func TestForwarder_Run_WithMockery(t *testing.T) {
	partners := []string{"synergyMarketing", "pixelAds", "impressionsMax"}

	mockReporter := mocks.NewClient(t)
	mockProducer := mocks.NewProducer(t)

	mockData := &reporter.Data{
		Clicks:      5,
		Cost:        3.14,
		Date:        "2025-07-01",
		Impressions: 42,
		Installs:    1,
	}

	for _, partner := range partners {
		mockReporter.On("Fetch", mock.Anything, partner, mock.AnythingOfType("string")).Return(mockData, nil).Once()
		mockProducer.On("Send", mock.Anything, partner, mockData).Return(nil).Once()
	}

	f := app.Forwarder{
		Partners: partners,
		Reporter: mockReporter,
		Producer: mockProducer,
		Interval: 10 * time.Millisecond,
	}

	// Cancel after one interval
	ctx, cancel := context.WithTimeout(context.Background(), 11*time.Millisecond)
	defer cancel()
	f.Run(ctx)

	// Assert that the expectations were met
	mockReporter.AssertExpectations(t)
	mockProducer.AssertExpectations(t)

	// Optionally, assert counts
	assert.Equal(t, len(partners), len(mockReporter.Calls), "reporter.Fetch call count")
	assert.Equal(t, len(partners), len(mockProducer.Calls), "producer.Send call count")
}
