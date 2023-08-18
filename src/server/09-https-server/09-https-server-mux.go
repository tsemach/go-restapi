package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func makeResponse(name string, route string) []byte {
	resp := make(map[string]string)
	resp["status"] = "ok"
	resp["message"] = "09-https-server/" + name
	resp["route"] = route
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	return jsonResp
}

func createMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(res http.ResponseWriter, req *http.Request) {
		jsonResp := makeResponse("09-http-server", "/")
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

	return mux
}
