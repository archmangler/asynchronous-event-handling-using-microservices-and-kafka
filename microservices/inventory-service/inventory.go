package main

//Inventory service for processing inbound received orders
//Needed:
//go get -u github.com/segmentio/kafka-go
//

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

//produce
var ctx = context.Background()

// the topic and broker address are initialized as constants
const (
	topic1         = "deadletter-queue"
	topic2         = "order-confirmed"
	topic0         = "order-received-events"
	broker1Address = "localhost:9092"
	broker2Address = "localhost:9092"
	broker3Address = "localhost:9092"
)

//Publish the message to kafka DeadLetter or other topics
//call: newOrderHandlers
func produce(message string, ctx context.Context, topic string) {

	i := 0

	// intialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker1Address, broker2Address, broker3Address},
		Topic:   topic,
	})

	// each kafka message has a key and value. The key is used
	// to decide which partition (and consequently, which broker)
	// the message gets published on
	err := w.WriteMessages(ctx, kafka.Message{
		Key: []byte(strconv.Itoa(i)),

		// create an arbitrary message payload for the value
		Value: []byte(message),
	})

	if err != nil {
		panic("could not write message " + err.Error() + "to topic" + topic)
	}

	// log a confirmation once the message is written
	fmt.Println("wrote:", message, " to topic ", topic)

	// sleep for a second
	time.Sleep(time.Second)
}

type healthPortal struct {
	message string
}

func newAdminPortal() *healthPortal {

	message := os.Getenv("SHELL")

	if message == "" {
		panic("can't access local environment")
	}

	return &healthPortal{message: message}
}

func (a healthPortal) handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><h1>Service Health Portal</h1></html>"))
}

func consume(ctx context.Context, topic string) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages

	fmt.Println("debug> consuming from topic ", topic)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker1Address, broker2Address, broker3Address},
		Topic:   topic,
		GroupID: "my-group",
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		fmt.Println("received: ", string(msg.Value))
	}
}

func main() {

	//because this blocks , we run in a go-routine
	go func() {

		health := newAdminPortal()
		http.HandleFunc("/health", health.handler)
		err := http.ListenAndServe(":8081", nil) //port will have to be dynamically allocated for scalability.
		if err != nil {
			panic(err)
		}

	}()

	//test message to kick off the stream
	message := "Mary had a little lamb ..."

	//produce a test message
	produce(message, ctx, topic0)

	//consume incoming order received events from the order received topic/"queue"
	consume(ctx, topic0)

}
