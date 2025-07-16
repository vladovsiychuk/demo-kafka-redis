package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github/vladovsiychuk/demo-kafkaredis-diff/internal/app"
	"github/vladovsiychuk/demo-kafkaredis-diff/internal/datastore"
	"github/vladovsiychuk/demo-kafkaredis-diff/internal/kafka"

	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

const (
	KAFKA_BOOTSTRAP_SERVER = "localhost:9092"
	EVENTS_TOPIC           = "events"
	DIFFS_TOPIC            = "diffs"
	REDIS_ADDR             = "localhost:6379"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := kafka.CreateTopicIfNotExists(KAFKA_BOOTSTRAP_SERVER, DIFFS_TOPIC, 1); err != nil {
		logrus.WithError(err).Fatal("failed to create Kafka topic")
	}

	consumer := kafka.NewConsumer(KAFKA_BOOTSTRAP_SERVER, EVENTS_TOPIC, "diff-calculator-group")
	producer := kafka.NewProducer(KAFKA_BOOTSTRAP_SERVER, DIFFS_TOPIC)
	defer producer.Close()

	store := datastore.NewRedisStore(REDIS_ADDR)

	calculator := app.DiffCalculator{
		Consumer: consumer,
		Producer: producer,
		Store:    store,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		logrus.Info("starting diff-calculator")
		calculator.Run(ctx)
		logrus.Info("diff-calculator stopped")
	}()

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, unix.SIGTERM, unix.SIGINT)
	<-sig
	cancel()

	wg.Wait()
	logrus.Info("service exited")
}
