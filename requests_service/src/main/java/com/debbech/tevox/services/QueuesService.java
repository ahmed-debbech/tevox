package com.debbech.tevox.services;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.charset.StandardCharsets;

public class QueuesService {

    private Logger log = LoggerFactory.getLogger(this.getClass());

    private final RabbitMqService queues;
    private static QueuesService instance = null;

    private QueuesService() {
        try {
            this.queues = new RabbitMqService();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    public static QueuesService getInstance(){
        if(instance == null){
            instance = new QueuesService();
        }
        return instance;
    }

    public void publish(String text){
        try {
            this.queues.publish("ocr_request_q", text.getBytes(StandardCharsets.UTF_8), null);
        } catch (Exception e) {
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
