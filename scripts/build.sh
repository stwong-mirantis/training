#!/usr/bin/env bash

TAG=$1
ALPINE_VER=3.16
GOLANG_VER=1.18.4
ALPINE=alpine:${ALPINE_VER}
GOLANG_ALPINE=golang:${GOLANG_VER}-alpine${ALPINE_VER}

echo "Building..."

docker build --build-arg ALPINE=$ALPINE --build-arg GOLANG_ALPINE=$GOLANG_ALPINE -t "messaging-server:$TAG" .
