package test

import (
	"context"
	"sync"
	"testing"
	"time"

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

	var wg sync.WaitGroup
	wg.Add(len(partners) * 2) // each partner gets Fetch and Send

	for _, partner := range partners {
		mockReporter.On("Fetch", mock.Anything, partner, mock.AnythingOfType("string")).
			Return(mockData, nil).
			Run(func(args mock.Arguments) { wg.Done() }).
			Once()
		mockProducer.On("Send", mock.Anything, partner, mockData).
			Return(nil).
			Run(func(args mock.Arguments) { wg.Done() }).
			Once()
	}

	c := make(chan time.Time, 1)
	mockTicker := &mocks.Ticker{}
	mockTicker.On("C").Return((<-chan time.Time)(c))
	mockTicker.On("Stop").Return()

	f := app.Forwarder{
		Partners: partners,
		Reporter: mockReporter,
		Producer: mockProducer,
		Ticker:   mockTicker,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() {
		f.Run(ctx)
		close(done)
	}()

	c <- time.Now() // trigger tick

	wg.Wait() // wait until all Fetch and Send calls are done

	cancel() // stop the forwarder loop
	<-done   // wait for it to finish

	mockReporter.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}
