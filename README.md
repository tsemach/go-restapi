### from: https://dev.to/aurelievache/learning-go-by-examples-part-2-create-an-http-rest-api-server-in-go-1cdm
### from: https://dev.to/stack-labs/introduction-to-taskfile-a-makefile-alternative-h92
### from: https://medium.com/rungo/secure-https-servers-in-go-a783008b36da

1. go mod init github.com/scraly/learning-go-by-examples/go-rest-api

2. run the code
````bash
go run src/main.go
````

3. Tasks
This using *Taskfile* https://taskfile.dev 
- go-task --list # list of tasks

4. Using go-swagger: 
- install from: https://github.com/go-swagger/go-swagger/blob/master/docs/install.md
- after installation check: `swagger version`
- check the swagger with: `go-task swagger.validate`

5. create swagger doc in doc/index.html
- run `go-task swagger.doc`

6. creating certificates
  A. openssl req  -new  -newkey rsa:2048  -nodes  -keyout localhost.key  -out localhost.csr
  B. openssl  x509  -req  -days 365  -in localhost.csr  -signkey localhost.key  -out localhost.crt

----
## Create a CA
### Create root certificate authority
````bash
openssl genrsa -des3 -passout pass:`cat ssl.pass` -out ca.key 2048

openssl req -x509 -new -nodes -passin pass:`cat ssl.pass` -key ca.key -sha256 -days 1825 -out ca.crt -subj "/C=IL/ST=Jerusalem/O=tsemach.org/OU=R&D/CN=GO RestAPI Root CA"
````
----

## Create Server Certificate Without SAN (just CN)
### Create server certificate key, ca.key
````bash
openssl genrsa -out server.key 2048

openssl req -new -key server.key -out server.csr -subj "/C=IL/ST=Jerusalem/O=tsemach.org/OU=R&D/CN=localhost"

openssl x509 -req -in server.csr -passin pass:`cat ssl.pass`  -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 825 -sha256 

````
----

----
## Create Server Certificate With SAN
### Create CA
````bash
openssl \
  req \
  -subj "/C=IL/ST=Jerusalem/O=tsemach.org/OU=R&D/CN=GO RestAPI Root CA" \
  -new \
  -x509 \
  -passout pass:`cat ssl.pass` \
  -keyout ca.key \
  -out ca.crt \
  -days 36500
````

### Create Server Key
````bash
openssl genrsa -out server.key 2048
````

### Create Server CSR
````bash
openssl req \
  -new \
  -key server.key \
  -out server.csr \
  -subj "/C=IL/ST=Jerusalem/O=tsemach.org/OU=R&D/CN=localhost" \
  -sha256 \
  -extensions v3_req \
  -reqexts SAN \
  -config <(cat /etc/ssl/openssl.cnf <(printf "[SAN]\nsubjectAltName=DNS:localhost")) \
````
----
````bash
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
````
----
## Create Client Certificate
----
````bash
openssl genrsa -out client.key 2048

openssl req -new -key client.key -out client.csr -subj "/C=IL/ST=Jerusalem/O=tsemach.org/OU=R&D/CN=I-AM-THE-MAN"

openssl x509 -req -in client.csr -passin pass:`cat ssl.pass` -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 825 -sha256 
----
````

## Create Client Certificate
----
````bash
openssl genrsa -des3 -passout:`cat client.pass` -out client.key 2048
openssl req -new -passin pass:`cat client.pass` -key client.key -out client.csr -subj "/C=IL/ST=Tel-Aviv/O=tsemach.org/OU=R&D/CN=I-AM-THE-MAN/GN=Something
openssl x509 -req -days 3650 pass:`cat client.pass` -in client.csr -CA ca.crt -CAKey ca.key -CAcreateserial -out client.crt
````

### Getting subject of a certificate
`openssl x509 -in ca.crt  -noout -subject`
----
"/C=IL/ST=Jerusalem/O=tsemach.org/OU=R&D/CN=00000000-0000-0000-0000-000000000000"

## Running the Server
#### running the server by
1. go run src/server/11-https-server/11-https-server*.go

#### calling health api
curl --cert certs/client.crt --key certs/client.key --cacert certs/ca.crt https://localhost:8080/health
