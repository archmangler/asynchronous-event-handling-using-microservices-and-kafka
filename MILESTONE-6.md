# Define and Use KPIs to Evaluate System Performance

Define and implement KPIs and KPI functionality in a microservices based system.

-  An essential part of any business is understanding the performance of the systems that run the business. 
- Defining system-level KPIs is a way to expose information that can provide visibility to the success of technology solutions and can lead to better decision-making capabilities.
- KPIs are quantifiable measurements that can help in evaluating the performance of a system. 
- Publishing these metrics to Kafka can make it possible for different responsible parties to take action on a single KPI metric and can expose information to technology stakeholders to help drive positive change.
- From a technology perspective, metrics can help prove the effectiveness of the KPI. 
- A metric is a single value, typically with attached metadata (name and value pairs of information that can provide additional context). 
- As an example, take tracking how many orders have been received: As a metric, that would be represented by a number. 
- But you could also provide how many products were in the order as metadata to that metric. 
- An example of a KPI that this metric could help prove is "average number of orders per hour.
- A KPI should be able to answer questions like the following: Are the objective and outcome clear? Can it be easily measured? Is it easy to express and explain? Can the outcome be affected by the business? Is there a context to the KPI?

# Objective

Define and publish two key performance indicators (KPIs) that can provide insight into the performance of the system built throughout this liveProject. Now that we have a fully functioning system, we want to capture system-level metrics to ensure that we know how the system is performing.

# Importance to project

Thinking about the measurements that will best evaluate the performance of the system will be critical to finding opportunities for improvement, monitoring system health, and analyzing patterns over time. System-level KPIs can be useful in holding the technology teams accountable for defining what success looks like in a quantifiable way.

# Workflow

## Define and publish the first key performance indicator (KPI).

*Define a KPI to track errors per minute for the Inventory consumer. One way to achieve this would be to create an event that sends a count of the number of errors encountered when processing an event for the Inventory consumer.*

"inventoryErrorsPerMinute" :  A KPI indicating the number of order processing errors encountered per minute during the entire run time of the service (from start to stop/restart). This is derived from the inventoryErrors metric which is a counter incremented by the inventory service on every error event.


[x]
*Define an event schema that represents the KPI defined in the previous step. The minimum information you should consider passing in this event would be the KPI name, KPI metric name and value, as well as a timestamp.*

```
{
    "kpiName" : "inventoryErrorsPerMinute"
    "metricName": "inventoryErrors",
    "metricType": "counter",
    "metricValue": 0,
    "metricTimestamp": 123134565654,
    "metaData": "inventory service error count over lifetime of service run"
}
```

Schema fields definition:

* "kpiName" : "inventoryErrorsPerMinute" : The total errors rolled up every minute
* "metricName": "inventoryErrors", : count of errors encountered per order processing event
* "metricType": "counter", : the type of metric 
* "metricValue": 0, : value of metric for this publishing instance
* "metricTimestamp": 123134565654, : time in UNIX format the event was published
* "metaData": "inventory service error count over lifetime of service run" 

[x]
*Create a new topic in Kafka to store these new KPI events. Use what you learned in Milestone 1. Create a KPI event when an error is detected in the consumer, and publish the event to a the new topic in Kafka that was created in the previous step.*


```
bin/kafka-topics.sh \
  --create --topic order-metrics \
  --replication-factor 1 \
  --partitions 1 \
  --bootstrap-server localhost:9092


base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % bin/kafka-topics.sh --list --bootstrap-server localhost:9092
__consumer_offsets
deadletter-events
enotification-events
order-error-events
order-received-events
orderconfirmed-events
orderpicked-events
quickstart-events

(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % bin/kafka-topics.sh \
  --create --topic order-metrics \
  --replication-factor 1 \
  --partitions 1 \
  --bootstrap-server localhost:9092
Created topic order-metrics.

(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % bin/kafka-topics.sh --list --bootstrap-server localhost:9092
__consumer_offsets
deadletter-events
enotification-events
order-error-events
order-metrics
order-received-events
orderconfirmed-events
orderpicked-events
quickstart-events
```

[?]
*Test the Inventory consumer by submitting a duplicate order to the OrderReceived topic. Ensure that this new KPI event makes it to your new Kafka topic.*

- Metric Generation Code:

The General Schema for the Metric in any microservice would be defined as follows:

```
type Metric struct {
	kpiName         string `json:"name"`
	metricName      string `json:"id"`
	metricType      string `json:"id"`
	metricValue     string `json:"id"`
	metricTimeStamp string `json:"time"`
	metaData        string `json:"data"`
}
```

For inventory service the Metric structure is implemented as follows:

```
type Metric struct {

				kpiName         string `json:"inventoryErrorsPerMinute"`
				metricName      string `json:"inventoryErrors"`
				metricType      string `json:"counter"`
				metricValue     string `json:"metricValue"`
				metricTimeStamp string `json:"metricTimestamp"`
				metaData        string `json:"metaData"`
                }
```

- Test Metric Generation Code

(metrics published to metrics-topic)

```
(base) welcome@Traianos-MacBook-Pro inventory-service % .inventory                                                                                                                                                  

debug> consuming from topic  order-received-events

.
.
.

wrote: inventory Errors Per Minute,inventory Errors,counter,123456,inventory-service-errors  to topic  order-metrics

```

Monitoring the order-metrics topic (we see the error metrics counter increasing with every errored message):

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 %  bin/kafka-console-consumer.sh \                            
  --topic order-metrics \     
  --from-beginning \
  --bootstrap-server localhost:9092

inventory Errors Per Minute,inventory Errors,counter,1,inventory-service-errors

inventory Errors Per Minute,inventory Errors,counter,2,inventory-service-errors

inventory Errors Per Minute,inventory Errors,counter,3,inventory-service-errors
```

## Define and publish the second KPI.

*Define a KPI of your own choosing to track a system-level metric. One suggestion could be to track the average and maximum latency (elapsed time) for processing an order. Check out the last resource in the Resources section for more potential examples.*

- KPI defined as orders-processed-per-minute
- Simple metric of order count over the runtime of the order inventory service


```
type Metric struct {

				kpiName         string `json:"inventoryOrdersPerMinute"`
				metricName      string `json:"inventoryOrders"`
				metricType      string `json:"counter"`
				metricValue     string `json:"metricValue"`
				metricTimeStamp string `json:"metricTimestamp"`
				metaData        string `json:"metaData"`
}
```

- metrics handling modifications to code:

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

			fmt.Println("Publishing Order Metrics to queue", metricsTopic)

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

			fmt.Println("Publishing Error Metrics to queue", metricsTopic)

			publishOrderMetric(orderMetricData, ctx, metricsTopic)

			if err != nil {
				fmt.Println("could not publish metric for this order ..." + err.Error())
			}

		}

	}
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
```



*Test the viability of the KPI definition by writing a justification that addresses the specific questions in the last item under “Notes” below. Define an event schema that represents the KPI you chose to define.*


*In the service that has the most context and information related to the KPI, create and publish the KPI event to a new topic in Kafka.*


```
(base) welcome@Traianos-MacBook-Pro inventory-service % ./inventory
debug> consuming from topic  order-received-events

received:  {"Name":"newOrder", "ID":"78912","Time":"223232113111","Data":"new order", "Eventname":"newOrder"}
.
.
.

Publishing Order Metrics to queue order-metrics
incorrect metric format (not readable json)invalid character '\x00' looking for beginning of value
Metrics data dump: s =  [inventory Orders Per Minute inventory Orders counter 1 inventory-service-orders]
Metrics data dump: s =  [0 0 0 5 0 0 0 27 105 110 118 101 110 116 111 114 121 32 79 114 100 101 114 115 32 80 101 114 32 77 105 110 117 116 101 0 0 0 16 105 110 118 101 110 116 111 114 121 32 79 114 100 101 114 115 0 0 0 7 99 111 117 110 116 101 114 0 0 0 1 49 0 0 0 24 105 110 118 101 110 116 111 114 121 45 115 101 114 118 105 99 101 45 111 114 100 101 114 115]
Metrics data dump: mData =  inventory Orders Per Minute,inventory Orders,counter,1,inventory-service-orders
.
.
.

wrote: inventory Orders Per Minute,inventory Orders,counter,1,inventory-service-orders  to topic  order-metrics

.
.
.

received:  {"Name":"newOrder", "ID":"78912","Time":"223232113111","Data":"new order", "Eventname":"newOrder}
.
.
.

wrote: inventory Errors Per Minute,inventory Errors,counter,1,inventory-service-errors  to topic  order-metrics
```

*Test the service and ensure that the KPI event makes it to the new topic in Kafka.*

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 %  bin/kafka-console-consumer.sh \
  --topic order-metrics \
  --from-beginning \
  --bootstrap-server localhost:9092
inventory Errors Per Minute inventory Errors counter 123456 inventory service errors
inventory Errors Per Minute inventory Errors counter 123456 inventory service errors
inventory Errors Per Minute inventory Errors counter 123456 inventory service errors
inventory Errors Per Minute,inventory Errors,counter,123456,inventory-service-errors
inventory Errors Per Minute,inventory Errors,counter,1,inventory-service-errors
inventory Errors Per Minute,inventory Errors,counter,2,inventory-service-errors
inventory Errors Per Minute,inventory Errors,counter,3,inventory-service-errors


inventory Orders Per Minute,inventory Orders,counter,0,inventory-service-orders
inventory Orders Per Minute,inventory Orders,counter,0,inventory-service-orders
inventory Orders Per Minute,inventory Orders,counter,0,inventory-service-orders
inventory Orders Per Minute,inventory Orders,counter,0,inventory-service-orders
inventory Orders Per Minute,inventory Orders,counter,1,inventory-service-orders
inventory Errors Per Minute,inventory Errors,counter,1,inventory-service-errors


inventory Orders Per Minute,inventory Orders,counter,1,inventory-service-orders
inventory Errors Per Minute,inventory Errors,counter,1,inventory-service-errors

inventory Errors Per Minute,inventory Errors,counter,2,inventory-service-errors
inventory Orders Per Minute,inventory Orders,counter,2,inventory-service-orders
```


**Deliverable**

[x] The deliverable for this milestone is two event schemas representing the two KPIs that were created. 
[x] You will also have modified at least two different Go services and/or consumers to incorporate publishing the defined KPIs as events.

