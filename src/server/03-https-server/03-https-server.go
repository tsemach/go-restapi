package main

import (
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

	// create a custom server
	s := &http.Server{
		Addr:    ":8080",
		Handler: nil, // use `http.DefaultServeMux`
	}

	// handle `/` route
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, "Hello Custom World!")
	})

	// run server on port "9000"
	log.Fatal(s.ListenAndServeTLS(certCrt, certKey))

}
