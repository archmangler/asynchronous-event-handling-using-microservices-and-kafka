package main

//Inventory service for processing inbound received orders
//Needed:
//go get -u github.com/segmentio/kafka-go
//

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

//produce
var ctx = context.Background()

const maxInt32 = 1<<(32-1) - 1

// the topic and broker address are initialized as constants
const (
	topic1         = "deadletter-events"
	topic2         = "orderconfirmed-events"
	topic0         = "order-received-events"
	metricsTopic   = "order-metrics"
	broker1Address = "localhost:9092"
	broker2Address = "localhost:9092"
	broker3Address = "localhost:9092"
)

//DRY - define this in a main library somewhere ...
type Order struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	Time      string `json:"time"`
	Data      string `json:"data"`
	Eventname string `json:"eventname"`
}

/*

* "kpiName" : "inventoryErrorsPerMinute" : The total errors rolled up every minute
* "metricName": "inventoryErrors", : count of errors encountered per order processing event
* "metricType": "counter", : the type of metric
* "metricValue": 0, : value of metric for this publishing instance
* "metricTimestamp": 123134565654, : time in UNIX format the event was published
* "metaData": "inventory service error count over lifetime of service run"

 */

/*
 			metricData.kpiName = ""
			metricData.metricName = ""
			metricData.metricType = ""
			metricData.metricValue = ""
			metricData.metricTimeStamp = ""
			metricData.metaData = ""

*/

//metric publish event structure
type errorMetric struct {
	kpiName         string `json:"inventoryErrorsPerMinute"`
	metricName      string `json:"inventoryErrors"`
	metricType      string `json:"counter"`
	metricValue     string `json:"metricValue"`
	metricTimeStamp string `json:"metricTimestamp"`
	metaData        string `json:"metaData"`
}

type orderMetric struct {
	kpiName         string `json:"inventoryOrdersPerMinute"`
	metricName      string `json:"inventoryOrders"`
	metricType      string `json:"counter"`
	metricValue     string `json:"metricValue"`
	metricTimeStamp string `json:"metricTimestamp"`
	metaData        string `json:"metaData"`
}

//order processing error counter for metrics
var errorCount int = 0
var orderCount int = 0

//Metric structure to hold metric data for duration of service running time
var errorMetricData errorMetric = errorMetric{}
var orderMetricData orderMetric = orderMetric{}

//A little help for encoding []string to []byte
func Encode(s []string) []byte {
	var b []byte
	b = writeLen(b, len(s))
	for _, ss := range s {
		b = writeLen(b, len(ss))
		b = append(b, ss...)
	}
	return b
}

func writeLen(b []byte, l int) []byte {
	if 0 > l || l > maxInt32 {
		panic("writeLen: invalid length")
	}
	var lb [4]byte
	binary.BigEndian.PutUint32(lb[:], uint32(l))
	return append(b, lb[:]...)
}

//order metrics publishing function
func publishOrderMetric(data orderMetric, ctx context.Context, topic string) (err error) {

	d := orderMetric{data.kpiName, data.metricName, data.metricType, data.metricValue, data.metricTimeStamp, data.metaData}

	s := []string{
		d.kpiName,
		d.metricName,
		d.metricType,
		d.metricValue,
		d.metaData,
	}

	//Ugly Hack
	enc_s := Encode(s)
	err = json.Unmarshal([]byte(enc_s), &s)

	//check the message format is wholesome
	if err != nil {
		fmt.Println("incorrect metric format (not readable json)" + err.Error())
	}

	//convert string array to string:
	mData := strings.Join(s, ",")

	fmt.Println("Metrics data dump: s = ", s)
	fmt.Println("Metrics data dump: s = ", enc_s)
	fmt.Println("Metrics data dump: mData = ", mData)

	//publish the metric to kafka topic
	produce(mData, ctx, metricsTopic)

	return nil
}

//metrics publishing function
func publishErrorMetric(data errorMetric, ctx context.Context, topic string) (err error) {

	d := errorMetric{data.kpiName, data.metricName, data.metricType, data.metricValue, data.metricTimeStamp, data.metaData}

	s := []string{
		d.kpiName,
		d.metricName,
		d.metricType,
		d.metricValue,
		d.metaData,
	}

	//ISSUES!!!!
	enc_s := Encode(s)

	err = json.Unmarshal([]byte(enc_s), &s)

	//check the message format is wholesome
	if err != nil {
		fmt.Println("incorrect metric format (not readable json)" + err.Error())
	}

	//convert string array to string:
	mData := strings.Join(s, ",")

	fmt.Println("Metrics data dump: s = ", s)
	fmt.Println("Metrics data dump: s = ", enc_s)
	fmt.Println("Metrics data dump: mData = ", mData)

	//publish the metric to kafka topic
	produce(mData, ctx, metricsTopic)

	return nil
}

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
			produce(message, ctx, topic1)

			//publish metric to the metrics topic
			errorCount += 1

			//Dirty hack, I know ...
			errorCount_s := strconv.Itoa(errorCount)
			errorMetricData.kpiName = "inventory Errors Per Minute"
			errorMetricData.metricName = "inventory Errors"
			errorMetricData.metricType = "counter"
			errorMetricData.metricValue = errorCount_s
			errorMetricData.metricTimeStamp = "121323455454553"
			errorMetricData.metaData = "inventory-service-errors"

			fmt.Println("Publishing Error Metrics to queue", metricsTopic)

			publishErrorMetric(errorMetricData, ctx, metricsTopic)

		} else {

			//If all went well ...
			produce(message, ctx, topic2)

			orderCount_s := strconv.Itoa(orderCount)
			orderMetricData.kpiName = "inventory Orders Per Minute"
			orderMetricData.metricName = "inventory Orders"
			orderMetricData.metricType = "counter"
			orderMetricData.metricValue = orderCount_s
			orderMetricData.metricTimeStamp = "121323455454553"
			orderMetricData.metaData = "inventory-service-orders"

			//publish metric to the metrics topic
			orderCount += 1

			fmt.Println("Publishing Order Metrics to queue", metricsTopic)

			publishOrderMetric(orderMetricData, ctx, metricsTopic)

			if err != nil {
				fmt.Println("could not publish metric for this order ..." + err.Error())
			}

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
		fmt.Println("incorrect message format (not readable json)" + err.Error())
	}

	if len(data.Name) == 0 {
		err = errors.New("incorrect message format, Name field empty")
		return err
	}

	if len(data.ID) == 0 {
		err = errors.New("incorrect message format, ID field empty")
		return err
	}

	if len(data.Data) == 0 {
		err = errors.New("incorrect message format, Data field empty")
		return err
	}

	return nil
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
	//message := "Mary had a little lamb ..."
	//produce a test message
	//produce(message, ctx, topic0)

	//consume incoming order received events from the order received topic/"queue"
	consume(ctx, topic0)

}
