package main

import (
	"fmt"
	"log"
	_ "log"
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

	// handle `/` route
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, "Hello World!")
	})

	// run server on port "9000"
	log.Fatal(http.ListenAndServeTLS(":8080", certCrt, certKey, nil))
}
