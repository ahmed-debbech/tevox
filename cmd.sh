#!/bin/bash

if [[ $1 == "build" ]]; then
    docker compose down || true
    docker compose build base
    docker compose up --build ocr_service rabbitmq
fi


if [[ $1 == "start" ]]; then
    docker compose start ocr_service rabbitmq
fi

if [[ $1 == "stop" ]]; then
    docker compose stop ocr_service rabbitmq
fi