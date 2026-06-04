package com.debbech.tevox.models;

public class Message {

    private String stanzaId;
    private String body;
    private MessageType messageType;

    public Message(String stanzaId, String body, MessageType messageType){
        this.stanzaId = stanzaId;
        this.body = body;
        this.messageType = messageType;
    }
}
