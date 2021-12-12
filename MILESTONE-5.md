# Objective

*Build three remaining consumers that are needed to create the complete end-to-end solution, leveraging what was learned in Milestone 3. At a high level, these consumers do similar things that the Inventory consumer does. They consume an event from a Kafka topic, publish any errors encountered with processing that event, and then publish a new event to a Kafka topic. The main difference between them is the action each takes when processing an incoming event.*

# Importance to project

*The goal of this milestone is to create the remaining event consumers that are required to complete the system. Once finished, an order can be placed, and communication about the order can traverse through everything that has been built in Milestones 2 through 5.*

# Workflow

## Create a new Notification consumer in Go.

*Create a long-lived subscription to the Notification topic in Kafka.Create functionality that extracts (and logs) the notification information from the relevant event schema.Treat received events in an idempotent manner, meaning any duplicates that are received must not create any side effects within the system. This can be achieved with two potential solutions: configuring the topic to enable idempotence or handling the duplicate checking in the event consumer by tracking the events processed to detect and discard duplicates. Create functionality that publishes an error event containing the received Notification event to the DeadLetterQueue topic in Kafka when the event can’t be processed successfully.Test that the Notification consumer works as expected by posting a Notification event to the Notification topic. Verify that the correct notification information is logged in the Notification consumer. If any errors occurred while processing the Notification event, you should be able to confirm that an event was published to the DeadLetterQueue Kafka topic. The easiest way to verify that an event exists in a topic is to use the command illustrated in Step 5 of the “Apache Kafka Quickstart” guide.*


## Create a new Warehouse consumer in Go.

*Create a long-lived subscription to the OrderConfirmed topic in Kafka.Create functionality that extracts (and logs) the order information from the relevant event schema. Treat received events in an idempotent manner, meaning any duplicates that are received must not create any side effects within the system. This can be achieved with two potential solutions: configuring the topic to enable idempotence or handling the duplicate checking in the event consumer by tracking the events processed to detect and discard duplicates. Create functionality that publishes an error event containing the received OrderConfirmed event to the DeadLetterQueue topic in Kafka when the event can’t be processed successfully. Create functionality that publishes an event notification that the customer’s order is being fulfilled.Test that the Warehouse consumer works as expected by posting an order payload to a running Order service. Verify that the correct order is logged in the Warehouse consumer and a notification event was received by a running the Notification consumer. If any errors occurred while processing the OrderConfirmed event, you should be able to confirm that an event was published to the DeadLetterQueue Kafka topic. The easiest way to verify that an event exists in a topic is to use the command illustrated in Step 5 of the “Apache Kafka Quickstart” guide.*

## Create a new Shipper consumer in Go.

*Create a long-lived subscription to the OrderPickedAndPacked topic in Kafka. Create functionality that extracts (and logs) the order and shipping information from the relevant event schema. Treat received events in an idempotent manner, meaning any duplicates that are received must not create any side effects within the system. This can be achieved with two potential solutions: configuring the topic to enable idempotence or handling the duplicate checking in the event consumer by tracking the events processed in order to detect and discard duplicates. Create functionality that publishes an error event containing the received OrderPickedAndPacked event to the DeadLetterQueue topic in Kafka when the event can’t be processed successfully. Create functionality that publishes an event notification to the customer stating the order is being shipped. Test that the Shipper consumer works as expected by posting an OrderPickedAndPacked event to the OrderPickedAndPacked topic and verify that the correct order is logged in the Shipper consumer and a notification event was received by a running Notification consumer. If any errors occurred while processing the OrderPickedAndPacked event, you should be able to confirm that an event was published to the DeadLetterQueue Kafka topic. The easiest way to verify that an event exists in a topic is to use the command illustrated in Step 5 of the “Apache Kafka Quickstart” guide.*

## Deliverable

- *Notification, which will be subscribed to the Notification topic in Kafka and can publish error events to the DeadLetterQueue topic in Kafka*

- *Warehouse, which will be subscribed to the OrderConfirmed topic in Kafka and can publish events to the Notification topic as well as publish error events to the DeadLetterQueue topic in Kafka*

- *Shipper, which will be subscribed to the OrderPickedAndPacked topic in Kafka and can publish events to the Notification topic as well as publish error events to the DeadLetterQueue topic in Kafka*

[] Check all relevant topics are there:


```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % bin/kafka-topics.sh --list --bootstrap-server localhost:9092
__consumer_offsets

->deadletter-events
->enotification-events

order-error-events
order-received-events

->orderconfirmed-events
->orderpicked-events

quickstart-events
```

[p] Create 3 microservice  copies based on template

```
(base) welcome@Traianos-MacBook-Pro microservices % ls -l
total 0
.
.
.
drwxr-xr-x  4 welcome  staff  128 Dec  9 22:50 notification-service
drwxr-xr-x  3 welcome  staff   96 Dec  9 22:36 shipper-service
drwxr-xr-x  3 welcome  staff   96 Dec  9 22:36 warehouse-service
```

[p] Customise each go template code to microservice task
 [x] notification service
 [x] warehouse-service
 [p] shipper-service

[x] Test message consumption and production

- sample notification:

```
{   
 "namespace": "org.industrial",
 "type": "record",
 "name": "OrderEmailNotification",
 "fields": [
     {"name": "order_id",  "type": "long",
       "doc":"The Universally unique id that identifies the order"},
     {"name": "time", "type": "long",
       "doc":"Time the order alert request was generated as UTC milliseconds from the epoch"},
     {"name": "event_id",  "type": "long",
     "doc":"The Universally unique event id that identifies this event"},
     {"name": "email_type", "type": "long",
        "type":{"type":"enum",
             "name":"email_notification_type",
             "symbols":["OrderConfirmed","OrderRecieved","OrderRejected","OrderPicked", "OrderShipped", "OrderDelivered"]},
       "doc":"Type of the email notification to be sent"}
 ]
}
```

[x] Test error event to deadletter topic

- Notifications:

Produce:

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % bin/kafka-console-producer.sh \
  --topic enotification-events \
   --bootstrap-server localhost:9092
>[{ "namespace": "org.industrial",   "etype": "record",   "name": "OrderEmailNotification","fields": "random stuff" }]
>{ "namespace": "org.industrial", "etype": "record", "name": "OrderEmailNotification", "fields": [  {"name": "order_id", "type": "long",   "doc":"The Universally unique id that identifies the order"},  {"name": "time", "type": "long",   "doc":"Time the order alert request was generated as UTC milliseconds from thee epoch"},  {"name": "event_id", "type": "long",  "doc":"The Universally unique event id that identifies this event"},  {"name": "email_type", "type": "long",   "type":{"type":"enum",    "name":"email_notification_type",    "symbols":["OrderConfirmed","OrderRecieved","OrderRejected","OrderPicked", "OrderShipped"," "OrderDelivered"]},   "doc":"Type of the email notification to be sent"} ] } 
>


```

[x] Received by notification service and registered as an error:


```
(base) welcome@Traianos-MacBook-Pro notification-service % ./notification 
debug> consuming from topic  enotification-events
received:  
DEBUG>   <DEBUG
incorrect message format (not readable json)unexpected end of JSON input
wrote:   to topic  deadletter-events
received:  
DEBUG>   <DEBUG
incorrect message format (not readable json)unexpected end of JSON input
wrote:   to topic  deadletter-events
received:  { "namespace": "org.industrial", "etype": "record", "name": "OrderEmailNotification", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "time", "type": "long", "doc":"Time the order alert request was generated as UTC milliseconds from the epoch"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "email_type", "type": "long", "type":{"type":"enum", "name":"email_notification_type", "symbols":["OrderConfirmed","OrderRecieved","OrderRejected","OrderPicked", "OrderShipped", "OrderDelivered"]}, "doc":"Type of the email notification to be sent"} ] }
DEBUG>  { "namespace": "org.industrial", "etype": "record", "name": "OrderEmailNotification", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "time", "type": "long", "doc":"Time the order alert request was generated as UTC milliseconds from the epoch"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "email_type", "type": "long", "type":{"type":"enum", "name":"email_notification_type", "symbols":["OrderConfirmed","OrderRecieved","OrderRejected","OrderPicked", "OrderShipped", "OrderDelivered"]}, "doc":"Type of the email notification to be sent"} ] } <DEBUG
wrote: { "namespace": "org.industrial", "etype": "record", "name": "OrderEmailNotification", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "time", "type": "long", "doc":"Time the order alert request was generated as UTC milliseconds from the epoch"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "email_type", "type": "long", "type":{"type":"enum", "name":"email_notification_type", "symbols":["OrderConfirmed","OrderRecieved","OrderRejected","OrderPicked", "OrderShipped", "OrderDelivered"]}, "doc":"Type of the email notification to be sent"} ] }  to topic  deadletter-events
received:  
DEBUG>   <DEBUG
incorrect message format (not readable json)unexpected end of JSON input
wrote:   to topic  deadletter-events
```

[x] Produced to the deadletter queue:

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 %  bin/kafka-console-consumer.sh \                            
  --topic enotification-events \                                 
  --from-beginning \
  --bootstrap-server localhost:9092

{ "namespace": "org.industrial", "etype": "record", "name": "OrderEmailNotification", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "time", "type": "long", "doc":"Time the order alert request was generated as UTC milliseconds from the epoch"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "email_type", "type": "long", "type":{"type":"enum", "name":"email_notification_type", "symbols":["OrderConfirmed","OrderRecieved","OrderRejected","OrderPicked", "OrderShipped", "OrderDelivered"]}, "doc":"Type of the email notification to be sent"} ] }

```

[x] Validate end-to-end flow of messages (happy case)


Correctly formatted message successfully produced to the notification queue:

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 %  bin/kafka-console-consumer.sh \
  --topic enotification-events \
  --from-beginning \
  --bootstrap-server localhost:9092

        { "namespace": "org.industrial", "etype": "record", "name": "OrderEmailNotification", "fields": [  {"name": "order_id", "type": "long",   "doc":"The Universally unique id that identifies the order"},  {"name": "time", "type": "long",   "doc":"Time the order alert request was generated as UTC milliseconds from the epoch"},  {"name": "event_id", "type": "long",  "doc":"The Universally unique event id that identifies this event"},  {"name": "email_type", "type": "long",   "type":{"type":"enum",    "name":"email_notification_type",    "symbols":["OrderConfirmed","OrderRecieved","OrderRejected","OrderPicked", "OrderShipped", "OrderDelivered"]},   "doc":"Type of the email notification to be sent"} ] } 

```


- Warehouse Service

[x] Create basic warehouse microservice based on template
[x] Define the event format
[x] Consume from orderConfirmed queue ("orderconfirmed-events")
[x] Publish errors to the DeadletterQueue topic ("deadletter-events")
[x] Publish a notification to the the notification queue that "order is being fulfilled" ("enotification-events") 
[x] Test end to end for error case and happy case

*Event Format:*


```
{ "namespace": "org.industrial", "doctype": "record", "name": "OrderConfirmed", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was confirmed as UTC milliseconds from the epoch"} ] }
```


```
{
  "namespace": "org.industrial",
  "doctype": "record",
  "name": "OrderPicked",
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
```
*Produce Test Message*

```

(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % bin/kafka-console-producer.sh \
  --topic orderconfirmed-events \                                                                                                                                                                                                                                                                                       
   --bootstrap-server localhost:9092

>{ "namespace": "org.industrial", "doctype": "record", "name": "OrderConfirmed", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was confirmed as UTC milliseconds from the epoch"} ] }
```

*Confirm event is consumed from orderconfirmed-events topic*


```
(base) welcome@Traianos-MacBook-Pro warehouse-service % ./warehouse 

received:  { "namespace": "org.industrial", "doctype": "record", "name": "OrderConfirmed", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was confirmed as UTC milliseconds from the epoch"} ] }
DEBUG>  { "namespace": "org.industrial", "doctype": "record", "name": "OrderConfirmed", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was confirmed as UTC milliseconds from the epoch"} ] } <DEBUG
wrote: { "namespace": "org.industrial", "doctype": "record", "name": "OrderConfirmed", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was confirmed as UTC milliseconds from the epoch"} ] }  to topic  enotification-events
```

*Check event is logged in the notifications topic*


```
(base) welcome@Traianos-MacBook-Pro warehouse-service % ./warehouse

debug> consuming from topic  orderconfirmed-events
received:  { "namespace": "org.industrial", "doctype": "record", "name": "OrderConfirmed", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was confirmed as UTC milliseconds from the epoch"} ] }

wrote: { "namespace": "org.industrial", "doctype": "record", "name": "OrderConfirmed", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was confirmed as UTC milliseconds from the epoch"} ] }  to topic  enotification-events
```

*Checking the notifications queue:*

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 %  bin/kafka-console-consumer.sh \
  --topic enotification-events \
  --from-beginning \
  --bootstrap-server localhost:9092


{ "namespace": "org.industrial", "doctype": "record", "name": "OrderConfirmed", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was confirmed as UTC milliseconds from the epoch"} ] }
```

*Produce broken json event message and test for error produce to deadletter*

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % bin/kafka-console-producer.sh \
 --topic orderconfirmed-events \                                                      
 --bootstrap-server localhost:9092

>{ "namespace": org.industrial", "doctype": "record", "name": "OrderPicked", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was confirmed as UTC milliseconds from the epoch"} ] } 
```

*connfirm warehouse service consumes it*

```
(base) welcome@Traianos-MacBook-Pro warehouse-service % ./warehouse 

debug> consuming from topic  orderconfirmed-events
received:  { "namespace": org.industrial", "doctype": "record", "name": "OrderPicked", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was confirmed as UTC milliseconds from the epoch"} ] } 

DEBUG>  { "namespace": org.industrial", "doctype": "record", "name": "OrderPicked", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was confirmed as UTC milliseconds from the epoch"} ] }  <DEBUG
incorrect message format (not readable json)invalid character 'o' looking for beginning of value{ "namespace": org.industrial", "doctype": "record", "name": "OrderPicked", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was confirmed as UTC milliseconds from the epoch"} ] } 

Error with orderconfirmed format: invalid character 'o' looking for beginning of value{ "namespace": org.industrial", "doctype": "record", "name": "OrderPicked", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was confirmed as UTC milliseconds from the epoch"} ] } 

wrote: { "namespace": org.industrial", "doctype": "record", "name": "OrderPicked", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was confirmed as UTC milliseconds from the epoch"} ] }   to topic  deadletter-events
```

*Confirm error result goes to deadletter queue*

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 %  bin/kafka-console-consumer.sh \                                                       
  --topic deadletter-events \
  --from-beginning \
  --bootstrap-server localhost:9092
```

- Shipper Service

[x] Create basic Shipper microservice based on template
[x] Define the shipper event format
[x] Consume from OrderPickedAndPacked  queue ("orderpicked-events")
[x] Publish errors to the DeadletterQueue topic ("deadletter-events")
[x] Publish a notification to the the notification queue that "order is being fulfilled" ("enotification-events") 
[x] Test end to end for error case and happy case

*Produce ordershipped event to orderpicked-events topic*

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % bin/kafka-console-producer.sh \
  --topic orderpicked-events \
   --bootstrap-server localhost:9092
>{ "namespace": "org.industrial", "doctype": "record", "name": "OrderShipped", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was shipped as UTC milliseconds from the epoch"} ] } 
>
```

*Check event is received to the topic*

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 %  bin/kafka-console-consumer.sh \                            
  --topic orderpicked-events \
  --from-beginning \
  --bootstrap-server localhost:9092
{ "namespace": "org.industrial", "doctype": "record", "name": "OrderShipped", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", ""type": "long", "doc":"Time the order was shipped as UTC milliseconds from the epoch"} ] } 
```

*Confirm notification is produced to the enotification-events topic*

```
(base) welcome@Traianos-MacBook-Pro shipper-service % ./shipper          
.
.
.
DEBUG>  { "namespace": "org.industrial","time":113245676757, "doctype": "record", "name": "OrderShipped", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was shipped as UTC milliseconds from the epoch"} ] }  <DEBUG

wrote: { "namespace": "org.industrial","time":113245676757, "doctype": "record", "name": "OrderShipped", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, 

{"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was shipped as UTC milliseconds from the epoch"} ] }   to topic  enotification-events
```

*Confirm error notification is sent to deadletter queue*

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % bin/kafka-console-producer.sh \
  --topic orderpicked-events \
   --bootstrap-server localhost:9092
>{ "namespace": "org.industrial", "doctype": "record", "name": "OrderShipped", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was shipped as UTC milliseconds from the epoch"} ] } 
>{ "namespace": "org.industrial","time":113245676757, "doctype": "record", "name": "OrderShipped", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was shipped as UTC milliseconds from the epoch"} ] } 
>{ "namespace": "org.industrial","time":113245676757, doctype": "record", "name": "OrderShipped", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was shipped as UTC milliseconds from the epoch"} ] }             
>
.
.
.
Error with ordershipped format: invalid character '>' looking for beginning of value>{ "namespace": "org.industrial","time":113245676757, "doctype": "record", "name": "OrderShipped", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was shipped as UTC milliseconds from the epoch"} ] } 
wrote: >{ "namespace": "org.industrial","time":113245676757, "doctype": "record", "name": "OrderShipped", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was shipped as UTC milliseconds from the epoch"} ] }   to topic  deadletter-events
received:  >{ "namespace": "org.industrial","time":113245676757, doctype": "record", "name": "OrderShipped", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was shipped as UTC milliseconds from the epoch"} ] }             
DEBUG>  >{ "namespace": "org.industrial","time":113245676757, doctype": "record", "name": "OrderShipped", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was shipped as UTC milliseconds from the epoch"} ] }              <DEBUG
incorrect message format (not readable json)invalid character '>' looking for beginning of value>{ "namespace": "org.industrial","time":113245676757, doctype": "record", "name": "OrderShipped", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was shipped as UTC milliseconds from the epoch"} ] }             
Error with ordershipped format: invalid character '>' looking for beginning of value>{ "namespace": "org.industrial","time":113245676757, doctype": "record", "name": "OrderShipped", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was shipped as UTC milliseconds from the epoch"} ] }             
wrote: >{ "namespace": "org.industrial","time":113245676757, doctype": "record", "name": "OrderShipped", "fields": [ {"name": "order_id", "type": "long", "doc":"The Universally unique id that identifies the order"}, {"name": "event_id", "type": "long", "doc":"The Universally unique event id that identifies this event"}, {"name": "time", "type": "long", "doc":"Time the order was shipped as UTC milliseconds from the epoch"} ] }               to topic  deadletter-events
```













