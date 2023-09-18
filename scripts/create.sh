curl -k --key certs/client.key --cert certs/client.crt  --cacert certs/ca.crt https://localhost:8080/api/v1/create -d '{"pid": "abc", "tid": "xyz"}'
