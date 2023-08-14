openssl req \
  -new \
  -key server.key \
  -out server.csr \
  -subj "/C=IL/ST=Jerusalem/O=tsemach.org/OU=R&D/CN=localhost" \
  -sha256 \
  -extensions v3_req \
  -reqexts SAN \
  -config <(cat /etc/ssl/openssl.cnf <(printf "[SAN]\nsubjectAltName=DNS:localhost")) \