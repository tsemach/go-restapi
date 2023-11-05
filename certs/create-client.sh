#!/bin/bash

openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr -subj "/C=IL/ST=Jerusalem/O=tsemach.org/OU=R&D/CN=I-AM-THE-MAN"
openssl x509 -req -in client.csr -passin pass:`cat ssl.pass`  -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 825 -sha256
