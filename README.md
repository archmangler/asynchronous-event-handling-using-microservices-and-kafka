# Asynchronous Event Handline Using Microservices and Kafka

## Overview

- This project demonstrates generation and asynchronous handling of events produced by microservices via a containerised and elastically scalable Kafka cluster.
- The objective is to demonstrate the generation and parallel handling of high volumes of messages for load testing of transactional systems (e.g order matching platforms) in a cloud environment.

## Design

The design consists of the following components:

- Layer of producers to produce requests based on information in a data store (DB or file storage)
- An asynchronous messaging queue to accept and order requests (Kafka) 
- A layer of consumers to service requests in the asynch queue and submit to 


