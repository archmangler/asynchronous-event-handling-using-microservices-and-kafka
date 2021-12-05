#!/bin/bash

#set this to whatever in the shell environment
#KAFKA_HOME=....

function start_zookeeper() {

     printf "Starting zookeeper ...\n"

     $(nohup ${KAFKA_HOME}/bin/zookeeper-server-start.sh -daemon config/zookeeper.properties > /dev/null 2>&1 &)
     
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
     
     $(nohup ${KAFKA_HOME}/bin/kafka-server-start.sh -daemon config/server.properties > /dev/null 2>&1 &)
     
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

  result=$(${KAFKA_HOME}/bin/kafka-topics.sh --create \
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
