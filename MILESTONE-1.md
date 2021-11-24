# MILESTONE 1 (Deliverable #1)
==============================

*1. Download and unpack Kafka then set the environment variable KAFKA_HOME to the directory of the unpacked location. The easiest way to start is to use Steps 1 and 2 in the “Apache Kafka Quickstart” guide.*


- start kafka services

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % bin/zookeeper-server-start.sh config/zookeeper.properties
.
.
.


(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % bin/kafka-server-start.sh config/server.properties
.
.
.

```

- Injecting sample test events:

```

bin/kafka-console-producer.sh \
  --topic quickstart-events \
   --bootstrap-server localhost:9092

```

- Reading the events:

```
bin/kafka-console-consumer.sh \
  --topic quickstart-events \
  --from-beginning \
  --bootstrap-server localhost:9092
```


*2. Define schemas that represent events that will be published and consumed in this e-commerce system. The context of “schema” in this liveProject only refers to a structured, textual representation of the data required to be communicated to the system as an event. Using a JSON document to represent the data is fairly standard. There is no requirement that you install or use the schema registry in Kafka or define an “Avro Schema” for this liveProject. (See the following reference to Kafka in Action, Chapter 11, for more information.)*

*2.1 Define an event schema representing when an order has been received.*


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

*2.2 Define an event schema representing when an order has been confirmed (is not a duplicate order and can be processed)*

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

*2.3. Define an event schema representing when an order has been picked from within a warehouse and packed (is ready to be shipped)*

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

*2.4 Define an event schema representing when an email notification needs to be sent out*

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

*2.5 Define an event schema representing when a consumer is unable to successfully process an event they have received. This “error event should contain the event that could not be processed.*

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

*3. Create a Kafka topic that will contain order received events, and verify it exists. The easiest way to do this is to use Step 3 in the “Apache Kafka Quickstart” guide.*

(https://kafka.apache.org/quickstart)



- Create a topic

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % bin/kafka-topics.sh --create --topic order-received-events --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092             

Error while executing topic command : Topic 'order-received-events' already exists.
[2021-11-24 00:45:42,265] ERROR org.apache.kafka.common.errors.TopicExistsException: Topic 'order-received-events' already exists.
 (kafka.admin.TopicCommand$)

```

- check topic status:

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % bin/kafka-topics.sh \
  --describe \
  --topic order-received-events \
  --bootstrap-server localhost:9092
Topic: order-received-events    TopicId: owS20NjkSq6AcoxMg1MkKg PartitionCount: 1       ReplicationFactor: 1    Configs: segment.bytes=1073741824
        Topic: order-received-events    Partition: 0    Leader: 0       Replicas: 0     Isr: 0
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % 

```

*4. Modify the topic created in Step 3 by increasing its retention time to three days.*

(documentation is unclear what to use here - lot's of differences based on kafka versions)

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % 
bin/kafka-topics.sh \
  --alter \
  --topic order-received-events  \
  --config retention.ms=4320

Exception in thread "main" java.lang.IllegalArgumentException: --bootstrap-server must be specified
        at kafka.admin.TopicCommand$TopicCommandOptions.checkArgs(TopicCommand.scala:568)
        at kafka.admin.TopicCommand$.main(TopicCommand.scala:48)
        at kafka.admin.TopicCommand.main(TopicCommand.scala)
```

Therefore we modify the server settings after shutting the server down:

```
############################# Log Retention Policy #############################

# The following configurations control the disposal of log segments. The policy can
# be set to delete segments after a period of time, or after a given size has accumulated.
# A segment will be deleted whenever *either* of these criteria are met. Deletion always happens
# from the end of the log.

# The minimum age of a log file to be eligible for deletion due to age
#log.retention.hours=168

log.retention.hours=72

```

*5. Create additional Kafka topics that use the same configuration as the topic created in Step 3.*

*5.1 Create a Kafka topic that will contain order confirmed events, and verify it exists.*

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % bin/kafka-topics.sh --create \
 --topic orderconfirmed-events \
 --replication-factor 1 \
 --partitions 1 \     
 --bootstrap-server localhost:9092
Created topic orderconfirmed-events.
```

*5.2 Create a Kafka topic that will contain order picked and packed events, and verify it exists.*

```

(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % bin/kafka-topics.sh \
  --create --topic orderpicked-events \
  --replication-factor 1 \
  --partitions 1 \
  --bootstrap-server localhost:9092

Created topic orderpicked-events.

```

*5.3 Create a Kafka topic that will contain notification events, and verify it exists.*

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % bin/kafka-topics.sh \
  --create --topic enotification-events \
  --replication-factor 1 \
  --partitions 1 \
  --bootstrap-server localhost:9092

Created topic enotification-events.

```

*5.4 Create a Kafka topic that will contain error events, and verify it exists.*

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % bin/kafka-topics.sh \
  --create --topic order-error-events \
  --replication-factor 1 \
  --partitions 1 \
  --bootstrap-server localhost:9092

Created topic order-error-events.

```

*6.Publish and consume events from the topics created. Create samples of each event schema created in Step 1 and test publishing/consuming them to and from each of their respective topics. The easiest way to do this is to use Steps 4 and 5 in the*

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 %  bin/kafka-console-consumer.sh \
  --topic quickstart-events \
  --from-beginning \
  --bootstrap-server localhost:9092
this is a test message
This is a second message

.
Test message 1
Test message 2
{
 "namespace": "org.industrial",
 "type": "record",
 "name": "OrderReceived",
 "order_id": "0000000001",
 "event_id": "000000001",
 "time": "1637769458"
}

{
 "namespace": "org.industrial",
 "type": "record",
 "name": "OrderConfirmed",
 "order_id": "0000000001",
 "event_id": "000000002",
 "time": "163776998"
}

{
 "namespace": "org.industrial",
 "type": "record",
 "name": "OrderPicked",
 "order_id": "0000000001",
 "event_id": "000000003",
 "time": "163777000"

{
 "namespace": "org.industrial",
 "type": "record",
 "name": "OrderEmailNotification",
 "order_id": "0000000001",
 "event_id": "000000004",
 "time": "163777500"
 "email_type": "OrderConfirmed"
}


```

*7. Deliverable*

**Startup and initialisation Script**


```
#!/bin/bash

function start_zookeeper() {

     printf "Starting zookeeper ...\n"

     $(nohup bin/zookeeper-server-start.sh -daemon config/zookeeper.properties > /dev/null 2>&1 &)
     
     status=$?

     if [[ $status == 0 ]]
     then
        printf "success! zookeeper started, result = $status\n"
     else
        printf "failed to start zookeeper ... $status\n"
     fi
    
     echo $status
}

function start_kafka() {
     
     printf "Starting kafka ...\n"
     
     $(nohup bin/kafka-server-start.sh -daemon config/server.properties > /dev/null 2>&1 &)
     
     status=$?
     
     if [[ $status == 0 ]]
     then
        printf "success! kafka started, result = $status\n"
     else
        printf "failed to start kafka ... $status\n"
     fi
     
     echo $status
}

function create_topics() {

  topic_name=$1

  result=$(bin/kafka-topics.sh --create \
              --topic $topic_name \
              --replication-factor 1 \
              --partitions 1 \
              --bootstrap-server localhost:9092
           )

   echo $result

}

#start zookeeper first
start_result=$(start_zookeeper)

#start kafka now
if [[ $start_result ]]
then
  printf "zookeeper started successfully ,,,\n"
  kafka_result=$(start_kafka)
  if [[ $kafka_result ]]
  then
    printf "kafka started successfully ,,,\n"
  fi
else
  printf "failed to start kafka ... $kafka_result\n"
fi

#create the standard topics
for topic in orderconfirmed-events orderpicked-events enotification-events order-error-events
do
   create_topics $topic
done

exit 0
```

Example run:

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % ./kafkaInit.sh 
zookeeper started successfully ,,,
kafka started successfully ,,,
[2021-11-25 01:27:26,375] ERROR org.apache.kafka.common.errors.TopicExistsException: Topic 'orderconfirmed-events' already exists.
 (kafka.admin.TopicCommand$)
Error while executing topic command : Topic 'orderconfirmed-events' already exists.
[2021-11-25 01:27:27,871] ERROR org.apache.kafka.common.errors.TopicExistsException: Topic 'orderpicked-events' already exists.
 (kafka.admin.TopicCommand$)
Error while executing topic command : Topic 'orderpicked-events' already exists.
[2021-11-25 01:27:29,331] ERROR org.apache.kafka.common.errors.TopicExistsException: Topic 'enotification-events' already exists.
 (kafka.admin.TopicCommand$)
Error while executing topic command : Topic 'enotification-events' already exists.
[2021-11-25 01:27:30,784] ERROR org.apache.kafka.common.errors.TopicExistsException: Topic 'order-error-events' already exists.
 (kafka.admin.TopicCommand$)
Error while executing topic command : Topic 'order-error-events' already exists.
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % 
```







