#!/bin/bash

if [[ $1 == "build" ]]; then
    docker compose down || true
    docker compose build base
    docker compose up --build
fi


if [[ $1 == "start" ]]; then
    docker compose start
fi

if [[ $1 == "stop" ]]; then
    docker compose stop
fi