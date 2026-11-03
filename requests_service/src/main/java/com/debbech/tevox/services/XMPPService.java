package com.debbech.tevox.services;

import com.debbech.tevox.config.Secrets;

import com.debbech.tevox.models.MessageType;
import org.jivesoftware.smack.*;
import org.jivesoftware.smack.filter.StanzaTypeFilter;
import org.jivesoftware.smack.packet.ExtensionElement;
import org.jivesoftware.smack.packet.Message;
import org.jivesoftware.smack.packet.Stanza;
import org.jivesoftware.smack.tcp.XMPPTCPConnection;
import org.jivesoftware.smack.tcp.XMPPTCPConnectionConfiguration;
import org.jivesoftware.smackx.ping.PingManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;

public class XMPPService {

    private Logger log = LoggerFactory.getLogger(this.getClass());

    private XMPPTCPConnection connection;

    private void maintainConnection(){
        CompletableFuture<Void> disconnected = new CompletableFuture<>();
        connection.addConnectionListener(new ConnectionListener() {
            @Override
            public void connectionClosed() {
                disconnected.completeExceptionally(
                        new RuntimeException("XMPP connection closed")
                );
            }

            @Override
            public void connectionClosedOnError(Exception e) {
                disconnected.completeExceptionally(e);
            }
        });
        try {
            disconnected.get();
        } catch (InterruptedException e) {
            log.error("connection to XMPP server disconnected because {}", e.getMessage());
            throw new RuntimeException(e);
        } catch (ExecutionException e) {
            log.error("connection to XMPP server disconnected because {}", e.getMessage());
            throw new RuntimeException(e);
        }
    }
    public void connect() throws RuntimeException {
        log.info("connecting to xmpp server...");
        try {
            XMPPTCPConnectionConfiguration config = XMPPTCPConnectionConfiguration.builder()
                    .setUsernameAndPassword(Secrets.xmppUsername, Secrets.xmppPassword)
                    .setXmppDomain(Secrets.xmppDomainName)
                    .setHost(Secrets.xmppHost)
                    .build();

            this.connection = new XMPPTCPConnection(config);
            this.connection.connect(); //Establishes a connection to the server
            this.connection.login();
            PingManager pingManager =
                    PingManager.getInstanceFor(connection);
            pingManager.setPingInterval(20);
            log.info("connected to {} and logged in as {} successfully!", Secrets.xmppDomainName, Secrets.xmppUsername );
        } catch (Exception e) {
            log.error("could not connect or login because: {}", e.getMessage());
            throw new RuntimeException(e);
        }
    }

    public void listenToMessages() throws RuntimeException{

        connection.addAsyncStanzaListener(new StanzaListener() {
            public void processStanza(Stanza stanza) {
                if (stanza instanceof Message) {
                    Message message = (Message) stanza;
                    String xml = message.toXML().toString();

                    com.debbech.tevox.models.Message m = null;
                    if(xml.contains("<x xmlns='jabber:x:oob'>")){ //if it is an image
                        log.info("processing Image coming from XMPP server");
                        String id = message.getStanzaId();
                        String body = message.getBody();
                        m = new com.debbech.tevox.models.Message(id, body, MessageType.IMAGE);
                    }else{
                        if((xml.contains("<body>") &&
                                (!xml.substring(xml.indexOf("<body>")+6).startsWith("</bo")))){
                            log.info("processing text message coming from XMPP server");
                            String id = message.getStanzaId();
                            String body = message.getBody();
                            m = new com.debbech.tevox.models.Message(id, body, MessageType.TEXT);
                        }
                    }
                    if(m != null) MessageProcessor.getInstance().appendMessage(m);
                }
            }
        }, StanzaTypeFilter.MESSAGE);

        this.maintainConnection();
    }


}
