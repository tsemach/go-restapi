package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func JsonPrettyEncode(data interface{}, out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "    ")
	if err := enc.Encode(data); err != nil {
		return err
	}
	return nil
}

func JsonPrettyPrint(j any) {
	var buffer bytes.Buffer

	err := JsonPrettyEncode(j, &buffer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(buffer.String())
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

	mux.HandleFunc("/api/v1/install", func(res http.ResponseWriter, req *http.Request) {
		var install InstallBodyStruct

		json.NewDecoder(req.Body).Decode(&install)

		fmt.Println("[/api/install] req.body:")
		JsonPrettyPrint(install)

		// jsonResp := makeResponse(name, "/api/v1/create")
		resp := make(map[string]string)
		resp["status"] = "ok"
		resp["message"] = "install message is:" + install.Pid + "/" + install.Tid
		resp["route"] = "/api/v1/install"
		jsonResp, _ := json.Marshal(resp)

		res.Write(jsonResp)
	})

	mux.HandleFunc("/api/v1/assign", func(res http.ResponseWriter, req *http.Request) {
		var assign AssignBodyStruct

		json.NewDecoder(req.Body).Decode(&assign)

		fmt.Println("[/api/assign] req.body:")
		JsonPrettyPrint(assign)

		// jsonResp := makeResponse(name, "/api/v1/create")
		resp := make(map[string]string)
		resp["status"] = "ok"
		resp["message"] = "assign message is:" + assign.Pid
		resp["route"] = "/api/v1/assign"
		jsonResp, _ := json.Marshal(resp)

		res.Write(jsonResp)
	})

	return mux
}
