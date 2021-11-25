# Asynchronous Event Handline Using Microservices and Kafka

## Overview

- This project demonstrates generation and asynchronous handling of events produced by microservices via a containerised and elastically scalable Kafka cluster.
- The objective is to demonstrate the generation and parallel handling of high volumes of messages for load testing of transactional systems (e.g order matching platforms) in a cloud environment.
- There is a requirement for:

a) High Volume Low Latency Message Generation and 
b) Accurate Test message Generation and 
c) Point in time test message replay 
d) Log collection and monitoring
e) Metrics Tracking and analytics (dashboarding)

## Operating Principle

- Generate high volumes of messages using parallel producers consuming data from storage and passing them to a parrallelised asynchronous queue (e.g kafka)
- Use parallel worker processes to consume the queue messages apply any required logic and make requests to a third party API.
- Log results to an additional parallised queue for analysis, monitoring and alerting

## Design

The design consists of the following components:

- Layer of producers (microservices) to produce requests based on information in a data store (DB or file storage)
- An asynchronous messaging queue to accept and order requests (Kafka) 
- A layer of consumers (microservices) to service requests in the asynch queue and submit to the third party API

## Implementation

(Ongoing)

### Kafka Cluster Deployer

### Kafka Message Schema

### Kafka Schema Registry

### Message Generator Service

### Message Consumer Service
