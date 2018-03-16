#!/bin/sh

cp ./Dockerfile ./Dockerfile_arm
sed -i 's/ubuntu:/armv7\/armhf-ubuntu:/' docker/Dockerfile_arm
