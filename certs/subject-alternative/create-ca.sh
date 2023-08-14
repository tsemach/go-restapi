openssl \
  req \
  -subj "/C=IL/ST=Jerusalem/O=tsemach.org/OU=R&D/CN=GO RestAPI Root CA" \
  -new \
  -x509 \
  -passout pass:`cat ssl.pass` \
  -keyout ca.key \
  -out ca.crt \
  -days 36500