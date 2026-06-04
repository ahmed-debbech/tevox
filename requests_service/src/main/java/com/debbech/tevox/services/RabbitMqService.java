package com.debbech.tevox.services;

import com.debbech.tevox.config.Secrets;
import com.rabbitmq.client.AMQP;
import com.rabbitmq.client.Channel;
import com.rabbitmq.client.Connection;
import com.rabbitmq.client.ConnectionFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.HashMap;
import java.util.Map;

public class RabbitMqService {

    private static Logger log = LoggerFactory.getLogger(RabbitMqService.class);


    public Connection connection;
    public Channel queueMainCh;

    public RabbitMqService() throws Exception{
        ConnectionFactory factory = new ConnectionFactory();
        factory.setUri(Secrets.rabbitMqUrl);
        factory.setVirtualHost("/");
        try {
            this.connection = factory.newConnection();
        } catch (Exception e) {
            log.error("Failed to connect to RabbitMQ because", e);
            throw e;
        }

        try {
            queueMainCh = connection.createChannel();
        } catch (IOException e) {
            log.error("Failed to open a channel", e);
            throw e;
        }

    }

    public void defineAllQueues() throws Exception {

        try {
            queueMainCh.queueDeclare(
                    "ocr_request_dlq_q",
                    true,
                    false,
                    false,
                    null
            );
        } catch (IOException e) {
            log.error("could not define ocr_request_dlq_q because: {}", e.getMessage());
            throw new Exception(
                    "Could not define "
                            + "ocr_request_dlq_q"
                            + " because: "
                            + e.getMessage()
            );
        }

        Map<String, Object> retryArgs = new HashMap<>();

        retryArgs.put("x-message-ttl", 5000);
        retryArgs.put("x-dead-letter-exchange", "");
        retryArgs.put("x-dead-letter-routing-key", "ocr_request_q");

        try {
            queueMainCh.queueDeclare(
                    "ocr_request_retry_q",
                    true,
                    false,
                    false,
                    retryArgs
            );
        } catch (IOException e) {
            log.error("could not define ocr_request_retry_q because: {}", e.getMessage());
            throw new Exception(
                    "Could not define "
                            + "ocr_request_retry_q"
                            + " because: "
                            + e.getMessage()
            );
        }

        Map<String, Object> mainArgs = new HashMap<>();

        mainArgs.put("x-dead-letter-exchange", "");
        mainArgs.put(
                "x-dead-letter-routing-key",
                "ocr_request_retry_q"
        );

        try {
           queueMainCh.queueDeclare(
                    "ocr_request_q",
                    true,
                    false,
                    false,
                    mainArgs
            );

        } catch (IOException e) {
            log.error("could not define ocr_request_q because: {}", e.getMessage());
            throw new Exception(
                    "Could not define "
                            + "ocr_request_q"
                            + " because: "
                            + e.getMessage()
            );
        }
    }

    public void publish(String queueName, byte[] content, Map<String, Object> headers) throws Exception {

        AMQP.BasicProperties props = new AMQP.BasicProperties.Builder()
                .contentType("text/plain")
                .headers(headers)
                .build();

        try {
            queueMainCh.basicPublish(
                    "",          // exchange
                    queueName,    // routing key (queue name)
                    false,       // mandatory
                    props,
                    content
            );
        } catch (Exception e) {
            log.error("could not publish to {} because: {}", queueName, e.getMessage());
            throw new RuntimeException(
                    "Could not publish to queue " + queueName +
                            " value=" + new String(content, StandardCharsets.UTF_8) +
                            " because " + e.getMessage(),
                    e
            );
        }
    }

    public void ack(long deliveryTag) throws Exception {
        queueMainCh.basicAck(deliveryTag, false);
    }

    public void nack(long deliveryTag) throws Exception {
        queueMainCh.basicNack(deliveryTag, false, false);
    }
}