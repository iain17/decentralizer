#!/bin/bash -x
openssl genrsa -out fake-ca.key 4096
openssl req -new -x509 \
-subj "/emailAddress=ca@example.com/CN=fake-ca/OU=ops/O=services/L=SomeCity/ST=CA/C=US" \
-days 365 -key fake-ca.key -out fake-ca.crt
openssl genrsa -out fake-server.key 1024
openssl req -new \
-subj "/emailAddress=svr@example.com/CN=fake-svr/OU=ops/O=services/L=SomeCity/ST=CA/C=US" \
-key fake-server.key -out fake-server.csr
openssl x509 -req -days 365 -in fake-server.csr -CA fake-ca.crt -CAkey fake-ca.key \
-extensions v3_server -extfile ./ssl-extensions-x509 \
-set_serial 01 -out fake-server.crt

