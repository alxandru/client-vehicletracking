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
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

func handleGet(events *kafka.Events) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("../../static/index.html")
		tmpl.Execute(w, PageInfo{
			MessagesTotal: len(*events),
		})
	}
}

func handleGetEvents(events *kafka.Events) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var eventDocs []kafka.EventDocument

		for _, event := range *events {
			eventDocs = append(eventDocs, kafka.EventDocument{Event: event})
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(eventDocs)
	}
}
