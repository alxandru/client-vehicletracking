package http

import (
	"encoding/json"
	"net/http"
	"text/template"

	"github.com/alxandru/client-vehicletracking/pkg/kafka"
	"github.com/gorilla/mux"
)

type PageInfo struct {
	MessagesTotal int
}

func NewHTTPServer(addr string, events *kafka.Events) *http.Server {
	r := mux.NewRouter()
	r.HandleFunc("/", handleGet(events)).Methods(http.MethodGet)
	r.HandleFunc("/events", handleGetEvents(events)).Methods(http.MethodGet)

	fs := http.FileServer(http.Dir("../../static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

func handleGet(events *kafka.Events) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("../../view/index.html")
		tmpl.Execute(w, PageInfo{
			MessagesTotal: len(*events),
		})
	}
}

func handleGetEvents(events *kafka.Events) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// var eventDocs []kafka.EventDocument

		// for _, event := range *events {
		// 	eventDocs = append(eventDocs, kafka.EventDocument{Event: event})
		// }
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusCreated)

		byt := []byte(`{"events":[{"entry":"NE","exit":"NV-Exit","id":7},
						{"entry":"N","exit":"SV-Exit","id":11},
						{"entry":"NV","exit":"SE-Exit","id":20},
						{"entry":"SV","exit":"N-Exit","id":17},
						{"entry":"N","exit":"NE-Exit","id":19},
						{"entry":"SV","exit":"NV-Exit","id":16}]}`)

		// json.NewEncoder(w).Encode(eventDocs)
		json.NewEncoder(w).Encode(string(byt))
	}
}
