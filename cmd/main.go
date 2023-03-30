package main

import (
	"fmt"
	"net/http"

	"github.com/flaviodepaula/sse/package/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() { //thread one

	out := make(chan amqp.Delivery)
	rabbitmqChannel, err := rabbitmq.OpenChannel()

	if err != nil {
		panic(err)
	}

	go rabbitmq.Consume("msgs", rabbitmqChannel, out) // Consume messages from RabbitMQ - thread two

	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		for m := range out {
			fmt.Fprintf(w, "event: message\n") // event name
			fmt.Fprintf(w, "data: %s\n\n", m.Body)
			w.(http.Flusher).Flush()
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "html/index.html")
	})
	http.ListenAndServe(":8080", nil)
}
