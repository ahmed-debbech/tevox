package com.debbech.tevox.services;

import com.debbech.tevox.models.DocumentEvent;
import com.debbech.tevox.models.Message;
import com.debbech.tevox.models.MessageType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.Queue;
import java.util.concurrent.PriorityBlockingQueue;

public class MessageProcessor {

    private static Logger log = LoggerFactory.getLogger(MessageProcessor.class);

    private static MessageProcessor instance = null;
    private Queue<Message> internalQueue;
    private Thread eventBuilderThread;
    private DocumentEvent potentialEvent = null;

    private MessageProcessor(){
        this.internalQueue = new PriorityBlockingQueue<>();
        this.eventBuilderThread = new Thread(this::buildTransmittableEvent);
        this.eventBuilderThread.start();
    }

    public static MessageProcessor getInstance(){
        if(instance == null){
            instance = new MessageProcessor();
        }
        return instance;
    }

    public void appendMessage(Message message){
        this.internalQueue.add(message);
        log.info("stanza-id {} - {} / queue-size: {}", message.getStanzaId(), message.getMessageType().toString(), internalQueue.size());
    }

    private String downloadImageAndGetPath(String url){
        return null;
    }

    private void emitEvent(){

    }

    private void buildTransmittableEvent(){
        while(true){
            Message msg = this.internalQueue.poll();
            if(msg == null) continue;

            if(potentialEvent == null){
                this.potentialEvent = new DocumentEvent();
            }

            if((potentialEvent.getTitle() == null) || (potentialEvent.getTitle().isEmpty())){
                if(!msg.getMessageType().equals(MessageType.TEXT)) continue;
                if(!msg.getBody().startsWith("a ")) continue;
                potentialEvent.setTitle(msg.getBody());
                log.info("this is a new event that starts with {}", potentialEvent.getTitle());
            }

            if(msg.getMessageType().equals(MessageType.IMAGE)){
                String path = downloadImageAndGetPath(msg.getBody());
                if(potentialEvent.getImagePaths() == null){
                    potentialEvent.setImagePaths(new ArrayList<>());
                }
                potentialEvent.getImagePaths().add(path);
                log.info("there has been new image added to {}, on img {}", potentialEvent.getTitle(), potentialEvent.getImagePaths().size());
            }

            log.info("eidiejd {}", msg.getBody());
            if(msg.getMessageType().equals(MessageType.TEXT)){
                if(!msg.getBody().equals("b")) continue;
                emitEvent();
                log.info("a new event is about to get transmitted {}", this.potentialEvent);
                this.potentialEvent = null;
            }


            try {
                Thread.sleep(5000);
            } catch (InterruptedException ignored) {}
        }
    }
}
