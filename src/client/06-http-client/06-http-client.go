package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type healthStruct struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type createStrucRequestType struct {
	Pid string `json:"pid"`
	Tid string `json:"tid"`
}

type createStructResponseType struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Route   string `json:"route"`
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func main() {
	// read in the cert file
	caCert, err := os.ReadFile("/home/tsemach/projects/go-restapi/certs/ca.crt")
	if err != nil {
		log.Fatalf("Failed to append %q to RootCAs: %v", "ca.crt", err)
	}

	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
		log.Println("No certs appended, using system certs only")
	}

	tr := &http.Transport{
		// TLSClientConfig: &tls.Config{InsecureSkipVerify: true, RootCAs: caCertPool},
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false, RootCAs: rootCAs, ServerName: "localhost"},
	}

	bodyBody := &createStrucRequestType{
		Pid: "abc",
		Tid: "xyz",
	}

	r, w := io.Pipe()
	go func() {
		json.NewEncoder(w).Encode(bodyBody)
		w.Close()
	}()

	client := &http.Client{Transport: tr, Timeout: 10 * time.Second}
	resp, err := client.Post("https://localhost:8080/api/v1/create", "application/json", r)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	var health healthStruct

	err = json.NewDecoder(resp.Body).Decode(&health)

	fmt.Println("health:", health)
	fmt.Println("health.status:", health.Status)
	fmt.Println("health.status:", health.Message)
}
