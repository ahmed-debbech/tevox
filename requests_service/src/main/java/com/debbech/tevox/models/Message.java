package com.debbech.tevox.models;

public class Message implements Comparable<Message>{

    private String stanzaId;
    private String body;
    private MessageType messageType;
    private String fromJid;

    @Override
    public String toString() {
        return "Message{" +
                "stanzaId='" + stanzaId + '\'' +
                ", body='" + body + '\'' +
                ", messageType=" + messageType +
                ", fromJid='" + fromJid + '\'' +
                '}';
    }

    public Message(String stanzaId, String body, MessageType messageType, String fromJid) {
        this.stanzaId = stanzaId;
        this.body = body;
        this.messageType = messageType;
        this.fromJid = fromJid;
    }

    public String getFromJid() {
        return fromJid;
    }

    public void setFromJid(String fromJid) {
        this.fromJid = fromJid;
    }

    @Override
    public int compareTo(Message o) {
        return o.stanzaId.compareTo(this.stanzaId);
    }

    public String getStanzaId() {
        return stanzaId;
    }

    public void setStanzaId(String stanzaId) {
        this.stanzaId = stanzaId;
    }

    public String getBody() {
        return body;
    }

    public void setBody(String body) {
        this.body = body;
    }

    public MessageType getMessageType() {
        return messageType;
    }

    public void setMessageType(MessageType messageType) {
        this.messageType = messageType;
    }
}
