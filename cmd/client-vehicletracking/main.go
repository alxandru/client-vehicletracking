package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alxandru/client-vehicletracking/pkg/kafka"
)

var (
	kafkaEndpoint = flag.String("kafkaendpoint", "localhost:9092", "Kafka Endpoint")
	topic         = flag.String("topic", "vehicletraffic", "Kafka Topic")
	groupid       = flag.String("groupid", "vehicletrafficgid", "Consumer Group")
)

func main() {
	flag.Parse()
	var wg = &sync.WaitGroup{}

	consumer := kafka.NewConsumer(*kafkaEndpoint, *topic, *groupid)
	consumer.StartConsumer(func(value string, err error) {
		if err != nil {
			fmt.Print("Got error while consuming topic: ", err)
		} else {
			fmt.Printf("Got message %s\n", value)
		}
	})
	wg.Add(1)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt)
	go func() {
		defer wg.Done()
		for {
			s := <-sigChan
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt:
				consumer.StopConsumer()
				return
			}
		}
	}()
	wg.Wait()
}
