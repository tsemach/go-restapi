package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

func makeResponse(name string, route string) []byte {
	resp := make(map[string]string)
	resp["status"] = "ok"
	resp["message"] = "07-https-server/" + name
	resp["route"] = route
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	return jsonResp
}

func createServer(name string, port int, cert *tls.Certificate) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		jsonResp := makeResponse(name, "/")
		res.Write(jsonResp)
	})

	mux.HandleFunc("/api/create", func(res http.ResponseWriter, req *http.Request) {
		jsonResp := makeResponse(name, "/api/create")
		res.Write(jsonResp)
	})

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", port), // :{port}
		Handler: mux,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*cert},
		},
	}

	return &server
}

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

	// create a WaitGroup
	wg := new(sync.WaitGroup)

	// add two goroutines to `wg` WaitGroup
	wg.Add(2)

	// goroutine to launch a server on port 8080
	go func() {
		server := createServer("ONE", 8080, &cert)
		fmt.Println(server.ListenAndServeTLS("", ""))
		wg.Done()
	}()

	// goroutine to launch a server on port 8081
	go func() {
		server := createServer("TWO", 8081, &cert)
		fmt.Println(server.ListenAndServeTLS("", ""))
		wg.Done()
	}()

	// wait until WaitGroup is done
	wg.Wait()
}
