## TEVOX

* run with `./start.sh`
* change config in `.env`

that's it...

### the queueing mechanism
If there is a Consumer queue in this service and there are RETRY and DLQ configured then:
* the retry queue is always dead-letterd to the main queue
* after the consumer dectects number of failed attempts (based on the `retry_count` header in the message) the consumer publishes the message to DLQ.
* retry queue has messages with TTL so when the TTL expires the message is requeued back to main queue