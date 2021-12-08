Milestone #4
============

#Create a new Inventory consumer in Go.

*Create a long-lived subscription to the OrderReceived topic in Kafka.*


Implemented using the following func:

```
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

		//Hereafter we consume from the named topic ...
		produce(message, ctx, topic2)

	}
}
```

*Create functionality that extracts (and logs) the order information from the relevant event schema.*

Implemented very lightly by the combined consumer and producer:


```
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

		//Hereafter we consume from the named topic ...
		produce(message, ctx, topic2)

	}
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

```

*Treat received events in an idempotent manner, meaning any duplicates that are received must not create any side effects within the system. This can be achieved by two potential solutions: configuring the topic to enable idempotence or handling the duplicate checking in the event consumer by tracking the events processed to detect and discard duplicates.*

The consumer seems to only consume unique ("new") events as tested below:

- Post the message to the `orders` microservice via REST:

```
(base) welcome@Traianos-MacBook-Pro inventory-service % curl -X POST  http://localhost:8080/orders \
  -H 'Content-Type: application/json' \
  -d '{"Name":"newOrder", "ID":"78916","Time":"223232113115","Data":"new order", "Eventname":"newOrder"}'
```

- Inventory service consumes the new event and writes it to the orderconfirmed-events topic:

```
(base) welcome@Traianos-MacBook-Pro inventory-service % go build inventory.go

(base) welcome@Traianos-MacBook-Pro inventory-service % ./inventory          
debug> consuming from topic  order-received-events

received:  {"Name":"newOrder", "ID":"78916","Time":"223232113115","Data":"new order", "Eventname":"newOrder"}
wrote: {"Name":"newOrder", "ID":"78916","Time":"223232113115","Data":"new order", "Eventname":"newOrder"}  to topic orderconfirmed-events

```

- Tail the "orderconfirmed-events" topic to see if messages consumed from the order-received-events make it through:

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 %  bin/kafka-console-consumer.sh \
  --topic orderconfirmed-events \
  --from-beginning \
  --bootstrap-server localhost:9092

Mary had a little lamb ...
{"Name":"newOrder", "ID":"78914","Time":"223232113113","Data":"new order", "Eventname":"newOrder"}

Mary had a little lamb ...

{"Name":"newOrder", "ID":"78915","Time":"223232113114","Data":"new order", "Eventname":"newOrder"}

^B^[[D{"Name":"newOrder", "ID":"78916","Time":"223232113115","Data":"new order", "Eventname":"newOrder"}

```

*Consider carefully which configuration option the consumer should use to handle Kafka message offset resets. Please refer to the notes below about offsets for more information.*

- Not sure I understand the question ... ????

*Publish an event to the OrderConfirmed Kafka topic when an order has been verified not to be a duplicate.*

Done by publishing via REST to the `orders` microservice , which writes to the orders-received-events topic.

```

(base) welcome@Traianos-MacBook-Pro inventory-service % ./inventory          
debug> consuming from topic  order-received-events

received:  {"Name":"newOrder", "ID":"78916","Time":"223232113115","Data":"new order", "Eventname":"newOrder"}
wrote: {"Name":"newOrder", "ID":"78916","Time":"223232113115","Data":"new order", "Eventname":"newOrder"}  to topic orderconfirmed-events
```

#Create functionality that publishes an error event containing the received OrderReceived event to the DeadLetterQueue topic in Kafka when the event can't be processed successfully. This is the first time you are being asked to publish an error event. You created an error event schema in Milestone 1. If any errors occur when processing the OrderReceived event, create and publish an error event to the DeadLetterQueue topic representing the error received.


[x] Create deadletter topic ("queue")
[x] Create an error checking code block
 [x] check the field values of the message event json content
[x] Create an error handling function
[x] Test error publish to deadletter queue

- Create dead letter queue:

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % bin/kafka-topics.sh --create --topic deadletter-events --replication-factor 1 --partitions 1 --bootstrap-server localhost:9092
Created topic deadletter-events.
.
.
.
.
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 %  bin/kafka-topics.sh --list --bootstrap-server localhost:9092                                                                 
__consumer_offsets
deadletter-events
enotification-events
order-error-events
order-received-events
orderconfirmed-events
orderpicked-events
quickstart-events
```

- Add function to publish failed event messages to dead letter queue:

```

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

```

#Test that the Inventory consumer works as expected by posting an order payload to a running Order service. Verify that the correct order was received and logged in the Inventory consumer. You should also be able to verify that an event was published to the OrderConfirmed Kafka topic. The easiest way to verify that an event exists in a topic is to use the command illustrated in Step 5 of the “Apache Kafka Quickstart” guide. If any errors occurred while processing the OrderReceived event, you should be able to confirm that an event was published to th

- Run the inventory microservice:

```
(base) welcome@Traianos-MacBook-Pro inventory-service % go build inventory.go

(base) welcome@Traianos-MacBook-Pro inventory-service % ./inventory          
debug> consuming from topic  order-received-events
received:  {"Name":"newOrder", "ID":"78923","Time":"223232113122","Data":"", "Eventname":"newOrder"}

DEBUG>  {"Name":"newOrder", "ID":"78923","Time":"223232113122","Data":"", "Eventname":"newOrder"} <DEBUG
wrote: {"Name":"newOrder", "ID":"78923","Time":"223232113122","Data":"", "Eventname":"newOrder"}  to topic  deadletter-events

```

- Post a order with missing data field to `orders` microservice:

```
(base) welcome@Traianos-MacBook-Pro inventory-service % curl -X POST  http://localhost:8080/orders \
  -H 'Content-Type: application/json' \
  -d '{"Name":"newOrder", "ID":"78923","Time":"223232113122","Data":"", "Eventname":"newOrder"}'
(base) welcome@Traianos-MacBook-Pro inventory-service % 
```

- Check the error is raised in `inventory` microservice:

```
DEBUG>  {"Name":"newOrder", "ID":"78923","Time":"223232113122","Data":"", "Eventname":"newOrder"} <DEBUG
wrote: {"Name":"newOrder", "ID":"78923","Time":"223232113122","Data":"", "Eventname":"newOrder"}  to topic  deadletter-events
```

- Check the deadletter queue has a message:


```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 %  bin/kafka-console-consumer.sh \
  --topic deadletter-events \
  --from-beginning \
  --bootstrap-server localhost:9092

{"Name":"newOrder", "ID":"78923","Time":"223232113122","Data":"", "Eventname":"newOrder"}

```



















