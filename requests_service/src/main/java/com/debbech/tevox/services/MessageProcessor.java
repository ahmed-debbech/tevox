package com.debbech.tevox.services;

import com.debbech.tevox.models.DocumentEvent;
import com.debbech.tevox.models.Message;
import com.debbech.tevox.models.MessageType;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.Response;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tools.jackson.databind.ObjectMapper;
import tools.jackson.databind.ObjectWriter;

import java.io.FileNotFoundException;
import java.io.FileOutputStream;
import java.io.IOException;
import java.time.Instant;
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
        Instant instant = Instant.now();
        long timeStampSeconds = instant.getEpochSecond();
        String outputFile = "/output/" + timeStampSeconds;

        OkHttpClient client = new OkHttpClient();

        Request request = new Request.Builder()
                .url(url)
                .build();

        try (Response response = client.newCall(request).execute()) {
            if (!response.isSuccessful()) {
                throw new IOException("Failed: " + response);
            }

            try (FileOutputStream fos = new FileOutputStream(outputFile)) {
                fos.write(response.body().bytes());
            } catch (FileNotFoundException e) {
                throw new RuntimeException(e);
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        } catch (IOException e) {
            log.error("an error occurred when downloading image... because {}", e.getMessage());
            return null;
        }

        return outputFile;
    }

    private void emitEvent(){
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        String json = ow.writeValueAsString(this.potentialEvent);
        log.info("event {}", json);
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
                potentialEvent.setTitle(msg.getBody().substring(2));
                log.info("this is a new event that starts with {}", potentialEvent.getTitle());
            }

            if(msg.getMessageType().equals(MessageType.IMAGE)){
                log.info("there is a new image added to {}", potentialEvent.getTitle());
                String path = downloadImageAndGetPath(msg.getBody());
                if(path == null) continue;

                if(potentialEvent.getImagePaths() == null){
                    potentialEvent.setImagePaths(new ArrayList<>());
                }
                potentialEvent.getImagePaths().add(path);
            }

            if(msg.getMessageType().equals(MessageType.TEXT)){
                if(!msg.getBody().equals("b")) continue;
                log.info("a new event is about to get transmitted {}", this.potentialEvent);
                emitEvent();
                this.potentialEvent = null;
            }


            try {
                Thread.sleep(5000);
            } catch (InterruptedException ignored) {}
        }
    }
}
