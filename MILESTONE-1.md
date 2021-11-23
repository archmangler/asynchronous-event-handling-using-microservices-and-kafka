#MILESTONE 1 (Deliverable #1)
============================

*1. Download and unpack Kafka then set the environment variable KAFKA_HOME to the directory of the unpacked location. The easiest way to start is to use Steps 1 and 2 in the “Apache Kafka Quickstart” guide.*

```



```



*2. Define schemas that represent events that will be published and consumed in this e-commerce system. The context of “schema” in this liveProject only refers to a structured, textual representation of the data required to be communicated to the system as an event. Using a JSON document to represent the data is fairly standard. There is no requirement that you install or use the schema registry in Kafka or define an “Avro Schema” for this liveProject. (See the following reference to Kafka in Action, Chapter 11, for more information.)*

        **2.1 Define an event schema representing when an order has been received.**


```
{
    
 "namespace": "org.industrial",
 "type": "record",
 "name": "OrderReceived",
 "fields": [
     {"name": "order_id",  "type": "long",
       "doc":"The Universally unique id that identifies the order"},
     {"name": "event_id",  "type": "long",
     "doc":"The Universally unique event id that identifies this event"},
     {"name": "time", "type": "long",
       "doc":"Time the order was received as UTC milliseconds from the epoch"}
 ]
}
```

        **2.2 Define an event schema representing when an order has been confirmed (is not a duplicate order and can be processed)**

```
{
    
 "namespace": "org.industrial",
 "type": "record",
 "name": "OrderConfirmed",
 "fields": [
     
     {"name": "order_id",  "type": "long",
       "doc":"The Universally unique id that identifies the order"},
     
     {"name": "event_id",  "type": "long",
     "doc":"The Universally unique event id that identifies this event"},
     
     {"name": "time", "type": "long",
       "doc":"Time the order was confirmed as UTC milliseconds from the epoch"}
 ]
}
```

    **2.3. Define an event schema representing when an order has been picked from within a warehouse and packed (is ready to be shipped)**

```
{
    
 "namespace": "org.industrial",
 "type": "record",
 "name": "OrderPicked",
 "fields": [
     {"name": "order_id",  "type": "long",
       "doc":"The Universally unique id that identifies the order"},
     
     {"name": "event_id",  "type": "long",
     "doc":"The Universally unique event id that identifies this event"},

     {"name": "time", "type": "long",
       "doc":"Time the order was confirmed as UTC milliseconds from the epoch"}
 ]
}
```

        **2.4 Define an event schema representing when an email notification needs to be sent out**

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

        **2.5 Define an event schema representing when a consumer is unable to successfully process an event they have received. This “error event should contain the event that could not be processed.**

```
{
    
 "namespace": "org.industrial",
 "type": "record",
 "name": "OrderProcessingError",
 "fields": [
     {"name": "time", "type": "long",
       "doc":"Time the order alert request was generated as UTC milliseconds from the epoch"},
     {"name": "event_id",  "type": "long",
     "doc":"The Universally unique event id that identifies this event"},
     {
        "name": "order_event_id", "type": "long",
        "type":{"type":"enum",
             "name":"order_event_error_type",
             "symbols":["","OrderRecieved","OrderRejected","OrderPicked", "OrderShipped", "OrderDelivered"]},
       "doc":"Type of the error for this order event"
       },
      {
          "name": "order_id",  "type": "long",
          "doc":"The Universally unique id that identifies the related order"
    },
 ]
}

```


*3. Create a Kafka topic that will contain order received events, and verify it exists. The easiest way to do this is to use Step 3 in the “Apache Kafka Quickstart” guide. *

(https://kafka.apache.org/quickstart)















