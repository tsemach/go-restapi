openssl \
  x509 \
  -req \
  -days 36500 \
  -in server.csr \
  -CA ca.crt \
  -CAkey ca.key \
  -CAcreateserial \
  -out server.crt \
  -extfile <(cat /etc/ssl/openssl.cnf <(printf "[SAN]\nsubjectAltName=DNS:localhost")) \
  -extensions SAN \
  -passin pass:`cat ssl.pass`