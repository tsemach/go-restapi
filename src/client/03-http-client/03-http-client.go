package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"

	// "io"
	"net/http"
	"time"
)

type healthStruct struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 10 * time.Second}
	resp, err := client.Get("https://localhost:8080/api")

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
