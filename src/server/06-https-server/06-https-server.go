package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mydir)

	certCrt := strings.Join([]string{mydir, "certs/server.crt"}, "/")
	certKey := strings.Join([]string{mydir, "certs/server.key"}, "/")

	fmt.Println("certCrt:", certCrt)
	fmt.Println("certKey:", certKey)

	// generate a `Certificate` struct
	cert, _ := tls.LoadX509KeyPair(certCrt, certKey)

	// create a custom server with `TLSConfig`
	s := &http.Server{
		Addr:    ":8080",
		Handler: nil, // use `http.DefaultServeMux`
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	// handle `/` route
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		resp := make(map[string]string)
		resp["status"] = "ok"
		resp["message"] = "06-https-server/"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		res.Write(jsonResp)

		return
	})

	http.HandleFunc("/api", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		resp := make(map[string]string)
		resp["status"] = "ok"
		resp["message"] = "06-https-server/api"
		
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(jsonResp)

		return
	})

	// run server on port "9000"
	log.Fatal(s.ListenAndServeTLS("", ""))
}
