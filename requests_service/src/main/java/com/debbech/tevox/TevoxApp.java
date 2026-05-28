package com.debbech.tevox;

import com.debbech.tevox.services.MainService;
import io.github.resilience4j.circuitbreaker.CircuitBreaker;
import io.github.resilience4j.decorators.Decorators;
import io.github.resilience4j.retry.Retry;
import io.github.resilience4j.retry.RetryConfig;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.function.Supplier;

class TevoxApp {

    private static Logger log = LoggerFactory.getLogger(TevoxApp.class);

    public static void main(String[] args) {
        log.info("Requests service has started!");

        Thread thread1 = new Thread(() -> {
            CircuitBreaker circuitBreaker = CircuitBreaker.ofDefaults("XMPPService");
            RetryConfig rc = new RetryConfig.Builder<>()
                    .maxAttempts(9999999)
                    .waitDuration(Duration.ofSeconds(5))
                    .retryExceptions(RuntimeException.class)
                    .build();

            Retry retry = Retry.of("XMPPService", rc);

            Supplier<String> supplier = MainService::execute1;
            Supplier<String> decoratedSupplier = Decorators.ofSupplier(supplier)
                    .withCircuitBreaker(circuitBreaker)
                    .withRetry(retry)
                    .decorate();

            decoratedSupplier.get();
        });

        Thread thread2 = new Thread(() -> {
            CircuitBreaker circuitBreaker1 = CircuitBreaker.ofDefaults("RabbitMqService");
            RetryConfig rc1 = new RetryConfig.Builder<>()
                    .maxAttempts(9999999)
                    .waitDuration(Duration.ofSeconds(5))
                    .retryExceptions(RuntimeException.class)
                    .build();

            Retry retry1 = Retry.of("RabbitMqService", rc1);

            Supplier<String> supplier1 = MainService::execute2;
            Supplier<String> decoratedSupplier1 = Decorators.ofSupplier(supplier1)
                    .withCircuitBreaker(circuitBreaker1)
                    .withRetry(retry1)
                    .decorate();

            decoratedSupplier1.get();
        });


        thread1.start();
        thread2.start();



    }
}