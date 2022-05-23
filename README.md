# Vehicle Tracking Client

This application consumes the Kafka events sent by the [Vehicle Tracking](https://github.com/alxandru/vehicle_tracking_deepstream) application and offers the posibility to access them via an http API. A Javascript client perdiodically queries the API and displays the information using a chord diagram.

The Kafka consumer was written in Go using the [kafka-go](ithub.com/segmentio/kafka-go) library.

The api web server was built with [gorilla/mux](github.com/gorilla/mux) router.

For the chord diagram the D3.js library was used.

## Table of contents

* [Requirements](#requirements)
* [Usage](#usage)

<a name="requirements"></a>

## Requirements

* [Golang](https://go.dev/dl/) (go1.18.1 version was used)
* [Vehicle Tracking Application](https://github.com/alxandru/vehicle_tracking_deepstream) (running and processing a video stream)

<a name="Usage"></a>

## Usage

Download the repo:

```bash
$ git clone https://github.com/alxandru/client-vehicletracking.git
$ cd client-vehicletracking/
```

Build the project:

```bash
$ cd cmd/client-vehicletracking
$ go build
```

Run the application indicating the Kafka endpoint and the Kafka topic. They have to be the same as in Vehicle Tracking Application:


```bash
$ ./client-vehicletracking -kafkaendpoint=<Host:port> -topic=vehicletraffic
```

Once the consumer starts to read events from the `vehicletraffic` topic the application start to output what it reads:

```bash
Reading Message
Got message {"event":{"entry":"N", "exit":"SV-Exit", "id":11}}
Reading Message
Got message {"event":{"entry":"NV", "exit":"SE-Exit", "id":20}}
Reading Message
Got message {"event":{"entry":"SV", "exit":"N-Exit", "id":17}}
Reading Message
Got message {"event":{"entry":"N", "exit":"NE-Exit", "id":19}}
...
```

Open a browser connecting to `http://127.0.0.1:8080/` to see how the chord diagram is built as the events are consumed:

![Alt Text](https://media.giphy.com/media/Yd0vSIHyUsqrdu8KDh/giphy.gif)



