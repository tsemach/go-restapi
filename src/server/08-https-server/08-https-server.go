package main

import (
	// "crypto/tls"
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
	resp["message"] = "08-https-server/" + name
	resp["route"] = route
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	return jsonResp
}

func createServer(name string, port int) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(res http.ResponseWriter, req *http.Request) {
		jsonResp := makeResponse(name, "/")
		res.Write(jsonResp)
	})

	mux.HandleFunc("/api/v1/create", func(res http.ResponseWriter, req *http.Request) {
		var create struct {
			Pid string `json:"pid"`
			Tid string `json:"tid"`
		}
		json.NewDecoder(req.Body).Decode(&create)

		fmt.Println("[/api/create] req.body:", create)

		// jsonResp := makeResponse(name, "/api/v1/create")
		resp := make(map[string]string)
		resp["status"] = "ok"
		resp["message"] = "create message is:" + create.Pid + "/" + create.Tid
		resp["route"] = "/api/v1/create"
		jsonResp, _ := json.Marshal(resp)

		res.Write(jsonResp)
	})

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", port), // :{port}
		Handler: mux,
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

	// create a WaitGroup
	wg := new(sync.WaitGroup)

	// add two goroutines to `wg` WaitGroup
	wg.Add(1)

	// goroutine to launch a server on port 8080
	go func() {
		server := createServer("server", 8080)
		fmt.Println(server.ListenAndServeTLS(certCrt, certKey))
		wg.Done()
	}()

	// wait until WaitGroup is done
	wg.Wait()
}
