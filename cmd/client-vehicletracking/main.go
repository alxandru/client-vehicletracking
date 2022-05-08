package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alxandru/client-vehicletracking/pkg/http"
	"github.com/alxandru/client-vehicletracking/pkg/kafka"
)

var (
	kafkaEndpoint      = flag.String("kafkaendpoint", "localhost:9092", "Kafka Endpoint")
	topic              = flag.String("topic", "vehicletraffic", "Kafka Topic")
	groupid            = flag.String("groupid", "vehicletrafficgid", "Consumer Group")
	serverAddress      = flag.String("serveraddress", "192.168.50.10:8080", "Server Address:port")
	totalKafkaMessages int
)

func main() {
	flag.Parse()

	var wg = &sync.WaitGroup{}
	var response = &kafka.Response{}

	consumer := kafka.NewConsumer(*kafkaEndpoint, *topic, *groupid)
	consumer.StartConsumer(func(value []byte, err error) {
		if err != nil {
			fmt.Print("Got error while consuming topic: ", err)
		} else {
			fmt.Printf("Got message %s\n", string(value))
			evDoc := &kafka.EventDocument{}
			if err := json.Unmarshal(value, evDoc); err != nil {
				fmt.Println("Unable to parse json: ", err)
			}
			response.Events = append(response.Events, evDoc)
			totalKafkaMessages += 1
		}
	})

	server := http.NewHTTPServer(*serverAddress, response)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error while listening: ", err)
	}

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
				server.Close()
				return
			}
		}
	}()
	wg.Wait()
}
