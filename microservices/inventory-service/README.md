Milestone #4
============

#Create a new Inventory consumer in Go.

*Create a long-lived subscription to the OrderReceived ("order-received-events") topic in Kafka.*

[check]

*Create functionality that extracts (and logs) the order information from the relevant event schema.*

[check]

*Treat received events in an idempotent manner, meaning any duplicates that are received must not create any side effects within the system. This can be achieved by two potential solutions: configuring the topic to enable idempotence or handling the duplicate checking in the event consumer by tracking the events processed to detect and discard duplicates.*

*Consider carefully which configuration option the consumer should use to handle Kafka message offset resets. Please refer to the notes below about offsets for more information.*

*Publish an event to the OrderConfirmed Kafka topic when an order has been verified not to be a duplicate.*

* Create functionality that publishes an error event containing the received OrderReceived event to the DeadLetterQueue topic in Kafka when the event can’t be processed successfully. This is the first time you are being asked to publish an error event. You created an error event schema in Milestone 1. If any errors occur when processing the OrderReceived event, create and publish an error event to the DeadLetterQueue topic representing the error received.*

#Test that the Inventory consumer works as expected by posting an order payload to a running Order service. Verify that the correct order was received and logged in the Inventory consumer. You should also be able to verify that an event was published to the OrderConfirmed Kafka topic. The easiest way to verify that an event exists in a topic is to use the command illustrated in Step 5 of the “Apache Kafka Quickstart” guide. If any errors occurred while processing the OrderReceived event, you should be able to confirm that an event was published to th
