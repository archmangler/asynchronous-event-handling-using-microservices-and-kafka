# Define and Use KPIs to Evaluate System Performance

Define and implement KPIs and KPI functionality in a microservices based system.

# Objective

Define and publish two key performance indicators (KPIs) that can provide insight into the performance of the system built throughout this liveProject. Now that we have a fully functioning system, we want to capture system-level metrics to ensure that we know how the system is performing.

# Importance to project

Thinking about the measurements that will best evaluate the performance of the system will be critical to finding opportunities for improvement, monitoring system health, and analyzing patterns over time. System-level KPIs can be useful in holding the technology teams accountable for defining what success looks like in a quantifiable way.

# Workflow

## Define and publish the first key performance indicator (KPI).

*Define a KPI to track errors per minute for the Inventory consumer. One way to achieve this would be to create an event that sends a count of the number of errors encountered when processing an event for the Inventory consumer.*

*Define an event schema that represents the KPI defined in the previous step. The minimum information you should consider passing in this event would be the KPI name, KPI metric name and value, as well as a timestamp.*

*Create a new topic in Kafka to store these new KPI events. Use what you learned in Milestone 1. Create a KPI event when an error is detected in the consumer, and publish the event to a the new topic in Kafka that was created in the previous step.*

*Test the Inventory consumer by submitting a duplicate order to the OrderReceived topic. Ensure that this new KPI event makes it to your new Kafka topic.*

## Define and publish the second KPI.

*Define a KPI of your own choosing to track a system-level metric. One suggestion could be to track the average and maximum latency (elapsed time) for processing an order. Check out the last resource in the Resources section for more potential examples.*

*Test the viability of the KPI definition by writing a justification that addresses the specific questions in the last item under “Notes” below. Define an event schema that represents the KPI you chose to define.*

*In the service that has the most context and information related to the KPI, create and publish the KPI event to a new topic in Kafka.*

*Test the service and ensure that the KPI event makes it to the new topic in Kafka.*


**Deliverable**

- The deliverable for this milestone is two event schemas representing the two KPIs that were created. 
- You will also have modified at least two different Go services and/or consumers to incorporate publishing the defined KPIs as events.
- An essential part of any business is understanding the performance of the systems that run the business. 
- Defining system-level KPIs is a way to expose information that can provide visibility to the success of technology solutions and can lead to better decision-making capabilities.
- KPIs are quantifiable measurements that can help in evaluating the performance of a system. 
- Publishing these metrics to Kafka can make it possible for different responsible parties to take action on a single KPI metric and can expose information to technology stakeholders to help drive positive change.
- From a technology perspective, metrics can help prove the effectiveness of the KPI. 
- A metric is a single value, typically with attached metadata (name and value pairs of information that can provide additional context). 
- As an example, take tracking how many orders have been received: As a metric, that would be represented by a number. 
- But you could also provide how many products were in the order as metadata to that metric. 
- An example of a KPI that this metric could help prove is "average number of orders per hour.
- A KPI should be able to answer questions like the following: Are the objective and outcome clear? Can it be easily measured? Is it easy to express and explain? Can the outcome be affected by the business? Is there a context to the KPI?

