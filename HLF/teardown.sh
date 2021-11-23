#!/bin/bash
#set -e



docker-compose -f docker-compose.yml down

# Shut down the Docker containers for the system tests.
docker-compose -f docker-compose.yml kill && docker-compose -f docker-compose.yml down

# remove the local state
rm -f ~/.hfc-key-store/*

chmod -R 0755 ./crypto-config
rm -fr config/*
rm -fr crypto-config/*

rm -fr js/wallet/*

# remove chaincode docker images
docker rm $(docker ps -aq)

# Your system is now clean
