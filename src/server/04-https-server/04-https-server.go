package main

import (
	"crypto/tls"
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

	certCrt := strings.Join([]string{mydir, "certs/localhost.crt"}, "/")
	certKey := strings.Join([]string{mydir, "certs/localhost.key"}, "/")

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
		fmt.Fprint(res, "[04-https-server] Hello Custom World!")
	})

	http.HandleFunc("/api", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, "[04-https-server] /api called")
	})

	// run server on port "9000"
	log.Fatal(s.ListenAndServeTLS("", ""))

}
