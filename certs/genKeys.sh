#! /bin/bash

# Use the below lines to generate the keys
openssl req -new -nodes -x509 -out server.pem -keyout server.key -days 365
openssl req -new -nodes -x509 -out client.pem -keyout client.key -days 365