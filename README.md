## TEVOX


* change config in `.env` only
* run with `./cmd.sh build` to build the images and start
* run with `./cmd.sh start` to start the already existing containers


### The Queueing Logic
If there is a Consumer queue in this service and there are RETRY and DLQ configured then:
* the retry queue is always dead-letterd to the main queue
* after the consumer dectects number of failed attempts (based on the `retry_count` header in the message) the consumer publishes the message to DLQ.
* retry queue has messages with TTL so when the TTL expires the message is requeued back to main queue