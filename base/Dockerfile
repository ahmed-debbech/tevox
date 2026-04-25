# This is the base docker image that will contain all service images
# we need to do this to keep all the dependencies cpu arch agnostic 
# so wherever cpu architecture you run it on, it should run right away

# the common things that are needed almost for every service image in this project are 
# 1 - ubuntu:24.04
# 2 - python and pip
# 3 - go 1.24.4

FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update &&  apt-get install -y g++ \
    autoconf automake libtool \
    pkg-config libpng-dev libjpeg8-dev \
    libtiff5-dev zlib1g-dev libwebpdemux2 \
    libwebp-dev libopenjp2-7-dev libgif-dev \
    libarchive-dev libcurl4-openssl-dev libleptonica-dev wget \
    build-essential --fix-missing


RUN apt-get install -y git

# compiling and installing go 1.4
WORKDIR /go1.4

RUN git clone -b release-branch.go1.4 https://go.googlesource.com/go

WORKDIR /go1.4/go/src

RUN export CGO_ENABLED=0 && ./make.bash

# compiling and instaling go1.17
WORKDIR /go1.17

RUN git clone --depth 1 -b go1.17 https://go.googlesource.com/go 

RUN cd go/src && export GOROOT_BOOTSTRAP=/go1.4/go && ./make.bash

# compiling and instaling go1.21.13
WORKDIR /go1.21.13

RUN git clone --depth 1 -b go1.21.13 https://go.googlesource.com/go 

RUN cd go/src && export GOROOT_BOOTSTRAP=/go1.17/go && ./make.bash

# compiling and instaling go1.22.6
WORKDIR /go1.22.6

RUN git clone --depth 1 -b go1.22.6 https://go.googlesource.com/go 

RUN cd go/src && export GOROOT_BOOTSTRAP=/go1.21.13/go && ./make.bash


# compiling and instaling go1.24.4
WORKDIR /go1.24.4

RUN git clone --depth 1 -b go1.24.4 https://go.googlesource.com/go 

RUN cd go/src && export GOROOT_BOOTSTRAP=/go1.22.6/go && ./make.bash

WORKDIR /

RUN rm -r go1.22.6 go1.21.13 go1.17 go1.4