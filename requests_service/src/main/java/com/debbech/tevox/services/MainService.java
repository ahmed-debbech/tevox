package com.debbech.tevox.services;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


public class MainService {
    private static Logger log = LoggerFactory.getLogger(MainService.class);

    public static String execute1() {
        XMPPService xmppService = new XMPPService();
        xmppService.connect();

        xmppService.listenToMessages();
        return null;
    }

    public static String execute2(){
        log.info("connecting to rabbitmq...");

        try {
            QueuesService rms = new QueuesService();
            rms.defineAllQueues();
            log.info("connected to rabbitmq successfully");
            Thread.sleep(5000);

            rms.publish();
        } catch (Exception e) {
            log.error("error defining queues {}", e.getStackTrace());
            throw new RuntimeException(e);
        }
        return null;
    }
}
