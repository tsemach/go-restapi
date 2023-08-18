package main

import (
	"crypto/tls"
	"crypto/x509"

	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var host = "localhost"
var port = 8080
var caCert = "certs/ca.crt"
var certOpt = tls.RequireAndVerifyClientCert

// func makeResponse(name string, route string) []byte {
// 	resp := make(map[string]string)
// 	resp["status"] = "ok"
// 	resp["message"] = "08-https-server/" + name
// 	resp["route"] = route
// 	jsonResp, err := json.Marshal(resp)
// 	if err != nil {
// 		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
// 	}

// 	return jsonResp
// }

// func createMux() *http.ServeMux {
// 	mux := http.NewServeMux()

// 	mux.HandleFunc("/health", func(res http.ResponseWriter, req *http.Request) {
// 		jsonResp := makeResponse("09-http-server", "/")
// 		res.Write(jsonResp)
// 	})

// 	mux.HandleFunc("/api/v1/create", func(res http.ResponseWriter, req *http.Request) {
// 		var create struct {
// 			Pid string `json:"pid"`
// 			Tid string `json:"tid"`
// 		}
// 		json.NewDecoder(req.Body).Decode(&create)

// 		fmt.Println("[/api/create] req.body:", create)

// 		// jsonResp := makeResponse(name, "/api/v1/create")
// 		resp := make(map[string]string)
// 		resp["status"] = "ok"
// 		resp["message"] = "create message is:" + create.Pid + "/" + create.Tid
// 		resp["route"] = "/api/v1/create"
// 		jsonResp, _ := json.Marshal(resp)

// 		res.Write(jsonResp)
// 	})

// 	return mux
// }

func getTLSConfig(host, caFullPath string, certOpt tls.ClientAuthType) *tls.Config {
	var caCert []byte
	var err error
	var caCertPool *x509.CertPool

	if certOpt > tls.RequestClientCert {
		caCert, err = os.ReadFile(caFullPath)
		if err != nil {
			log.Fatal("ERROR opening cert file", caFullPath, ", error ", err)
		}
		caCertPool = x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
	}

	return &tls.Config{
		ServerName: host,
		// ClientAuth: tls.NoClientCert,				// Client certificate will not be requested and it is not required
		// ClientAuth: tls.RequestClientCert,			// Client certificate will be requested, but it is not required
		// ClientAuth: tls.RequireAnyClientCert,		// Client certificate is required, but any client certificate is acceptable
		// ClientAuth: tls.VerifyClientCertIfGiven,		// Client certificate will be requested and if present must be in the server's Certificate Pool
		// ClientAuth: tls.RequireAndVerifyClientCert,	// Client certificate will be required and must be present in the server's Certificate Pool
		ClientAuth: certOpt,
		ClientCAs:  caCertPool,
		MinVersion: tls.VersionTLS12, // TLS versions below 1.2 are considered insecure - see https://www.rfc-editor.org/rfc/rfc7525.txt for details
	}
}

func createServer() *http.Server {
	mux := createMux()
	tls := getTLSConfig(host, caCert, tls.ClientAuthType(certOpt))

	return &http.Server{
		Addr:      fmt.Sprintf(":%v", port),
		Handler:   mux,
		TLSConfig: tls,

		ReadTimeout:  20 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func main() {
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mydir)

	certCrt := strings.Join([]string{mydir, "certs/server.crt"}, "/")
	certKey := strings.Join([]string{mydir, "certs/server.key"}, "/")

	// create a WaitGroup
	wg := new(sync.WaitGroup)

	// add two goroutines to `wg` WaitGroup
	wg.Add(1)

	// goroutine to launch a server on port 8080
	go func() {
		server := createServer()
		fmt.Println(server.ListenAndServeTLS(certCrt, certKey))
		wg.Done()
	}()

	// wait until WaitGroup is done
	wg.Wait()
}
