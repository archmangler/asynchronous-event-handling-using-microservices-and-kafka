package main

//Shipper microservice consumes order shipped events from topic `orderpicked-events` and sends "order shipped"
//messages to the notification queue  ("enotification-events") checks for errors and if errors are found in the notification, log them to
//the deadletter queue.
//Needed:
//go get -u github.com/segmentio/kafka-go
//

import (
	"context"
	"encoding/json"
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
	topic0         = "orderpicked-events"
	topic1         = "deadletter-events"
	topic2         = "enotification-events"
	broker1Address = "localhost:9092"
	broker2Address = "localhost:9092"
	broker3Address = "localhost:9092"
)

/*
Sample:
{
  "namespace": "org.industrial",
  "type": "record",
  "name": "OrderShipped",
  "fields": [
    {
      "name": "order_id",
      "type": "long",
      "doc": "The Universally unique id that identifies the order"
    },
    {
      "name": "event_id",
      "type": "long",
      "doc": "The Universally unique event id that identifies this event"
    },
    {
      "name": "time",
      "type": "long",
      "doc": "Time the order was confirmed as UTC milliseconds from the epoch"
    }
  ]
}
*/

type Order struct {
	namespace string `json:"namespace"`
	doctype   string `json:"doctype"`
	name      string `json:"name"`
	fields    string `json:"fields"`
}

//Publish the message to kafka DeadLetter or other topics
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

func consume(ctx context.Context, topic string) (message string) {
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
		message = string(msg.Value)

		fmt.Println("DEBUG> ", message, "<DEBUG")

		err = data_check(message)

		if err != nil {

			//produce to deadletter topic
			fmt.Println("Error with ordershipped format: " + err.Error() + message)
			produce(message, ctx, topic1)

		} else {

			//If all went well ...
			produce(message, ctx, topic2)

		}

	}
}

/*
A dummy error checking routine
*/

func data_check(message string) (err error) {

	data := Order{}

	err = json.Unmarshal([]byte(message), &data)

	if err != nil {
		fmt.Println("incorrect message format (not readable json)" + err.Error() + message)
		return err
	}

	/* Temporary block out due to issues (please troubleshoot further)
	if len(data.namespace) == 0 {
		fmt.Println("Namespace field value: ", data.namespace, " message: ", message)
		err = errors.New("incorrect message format, name field empty")
		return err
	}
	*/

	return nil
}

func main() {

	//The microservice health portal ... note: this scheme doesn't scale per microservice
	//you'd need to have a separate management microservice communicating with n child microservices
	//via a lightweight management port to each new dynamically created child service ...good use case for K8s.
	//because this blocks , we run in a go-routine
	go func() {

		health := newAdminPortal()
		http.HandleFunc("/health", health.handler)
		err := http.ListenAndServe(":8084", nil) //port will have to be dynamically allocated for scalability.
		if err != nil {
			panic(err)
		}

	}()

	//test message to kick off the stream
	//message := "Mary had a little lamb ..."
	//produce a test message
	//produce(message, ctx, topic0)

	//consume incoming order confirmed events from the orderconfirmed  topic/"queue"
	consume(ctx, topic0)

}
