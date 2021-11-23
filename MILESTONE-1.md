MILESTONE 1 (Deliverable #1)

    #1. Download and unpack Kafka then set the environment variable KAFKA_HOME to the directory of the unpacked location. The easiest way to start is to use Steps 1 and 2 in the “Apache Kafka Quickstart” guide.

    #2. Define schemas that represent events that will be published and consumed in this e-commerce system. The context of “schema” in this liveProject only refers to a structured, textual representation of the data required to be communicated to the system as an event. Using a JSON document to represent the data is fairly standard. There is no requirement that you install or use the schema registry in Kafka or define an “Avro Schema” for this liveProject. (See the following reference to Kafka in Action, Chapter 11, for more information.)

        ##2.1 Define an event schema representing when an order has been received.

        


        ##2.2 Define an event schema representing when an order has been confirmed (is not a duplicate order and can be processed).


        ##2.3. Define an event schema representing when an order has been picked from within a warehouse and packed (is ready to be shipped).


        ##2.4 Define an event schema representing when an email notification needs to be sent out.


        ##2.5 Define an event schema representing when a consumer is unable to successfully process an event they have received. This “error event should contain the event that could not be processed.


