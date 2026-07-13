#!/bin/bash

if [[ $1 == "build" ]]; then
    docker compose down || true
    docker compose build base
    docker compose build
    docker compose up --build --no-attach rabbitmq
fi


if [[ $1 == "start" ]]; then
    docker compose start
fi

if [[ $1 == "stop" ]]; then
    docker compose stop
fi