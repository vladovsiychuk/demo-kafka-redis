package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"justtrack.io/tests/backend/forwarder/internal/app"
	"justtrack.io/tests/backend/forwarder/internal/kafka"
	"justtrack.io/tests/backend/forwarder/internal/reporter"
)

const (
	KAFKA_BOOTSTRAP_SERVER = "localhost:9092"
	REPORTER_BASE_URL      = "http://localhost:8030"
	KAFKA_TOPIC            = "events"
)

var partners = []string{"synergyMarketing", "pixelAds", "impressionsMax"}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := kafka.CreateTopicIfNotExists(KAFKA_BOOTSTRAP_SERVER, KAFKA_TOPIC, 1); err != nil {
		logrus.WithError(err).Fatal("failed to create Kafka topic")
	}

	reporterClient := reporter.NewClient(REPORTER_BASE_URL)
	kafkaProducer := kafka.NewProducer(KAFKA_BOOTSTRAP_SERVER, KAFKA_TOPIC)
	defer kafkaProducer.Close()

	forwarder := app.Forwarder{
		Partners: partners,
		Reporter: reporterClient,
		Producer: kafkaProducer,
		Ticker:   app.NewRealTicker(time.Second),
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		logrus.Info("starting forwarder")
		forwarder.Run(ctx)
		logrus.Info("forwarder stopped")
	}()

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, unix.SIGTERM, unix.SIGINT)
	<-sig
	cancel()

	wg.Wait()
	logrus.Info("service exited")
}
