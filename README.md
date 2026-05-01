## TEVOX


* change config in `.env` only
* run with `./cmd.sh build` to build the images and start
* run with `./cmd.sh start` to start the already existing containers

### the 'shared' folder
Is the hub where all services share their work. \
should be in this format: \
`pic_in`: contains entry images to be text (mainly the entrypoint of the pipeline) \
`text_out`: output of the ocr service is set here and also the input of the voice service \
`voice_out`: output of the voice service in wav format.
### The Queueing Logic
If there is a Consumer queue in this service and there are RETRY and DLQ configured then:
* the retry queue is always dead-letterd to the main queue
* after the consumer dectects number of failed attempts (based on the `retry_count` header in the message) the consumer publishes the message to DLQ.
* retry queue has messages with TTL so when the TTL expires the message is requeued back to main queue