package com.debbech.tevox.services;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.charset.StandardCharsets;

public class RabbitMqService {

    private Logger log = LoggerFactory.getLogger(this.getClass());

    private final Queues queues = new Queues();

    public RabbitMqService() throws Exception {
    }

    public void publish(){
        try {
            this.queues.publish("ocr_request_q", new String("yo").getBytes(StandardCharsets.UTF_8), null);
        } catch (Exception e) {
            log.error("##############could not publish message to queue {}", e.getMessage());
            throw new RuntimeException(e);
        }
    }
    public void defineAllQueues(){
        try {
            this.queues.defineAllQueues();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }
}
