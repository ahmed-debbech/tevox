package com.debbech.tevox.services;

import com.debbech.tevox.models.Message;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Queue;
import java.util.concurrent.PriorityBlockingQueue;

public class MessageProcessor {

    private static Logger log = LoggerFactory.getLogger(MessageProcessor.class);

    private static MessageProcessor instance = null;
    private Queue<Message> internalQueue;

    private MessageProcessor(){
        this.internalQueue = new PriorityBlockingQueue<>();
    }

    public static MessageProcessor getInstance(){
        if(instance == null){
            instance = new MessageProcessor();
        }
        return instance;
    }

    public void appendMessage(Message message){
        this.internalQueue.add(message); log.info("{} / {}", internalQueue.peek().getStanzaId(), internalQueue.size());
    }

}
