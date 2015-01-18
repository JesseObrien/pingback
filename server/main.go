package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
)

var port string
var DebugMode bool

func init() {
	flag.StringVar(&port, "port", ":7733", "The port your server can run on. Default 7733.")
	flag.BoolVar(&DebugMode, "debugMode", false, "Turn debug on or off.")
}

func DebugPrint(a ...interface{}) {
	if !DebugMode {
		fmt.Println(a)
	}
}

type PingRequest struct {
	Host string
}

type PingResponse struct {
	Host   string
	Status string
}

func ResolveHost(host string) PingResponse {
	DebugPrint("Resolving host:" + host)
	resp := PingResponse{}
	getResp, err := http.Head(host)

	defer getResp.Body.Close()

	if err != nil {
		resp.Status = "failed"
	} else {
		resp.Status = getResp.Status
	}

	return resp
}

func handlePingRequest(conn net.Conn) {
	// Send some shit
	defer conn.Close()

	var r PingRequest
	decoder := json.NewDecoder(conn)

	if err := decoder.Decode(&r); err != nil {
		panic(err)
	}

	DebugPrint("Request to resolve host received.")

	resp := ResolveHost(r.Host)

	DebugPrint("Host resolved.")

	encoder := json.NewEncoder(conn)
	encoder.Encode(resp)
	DebugPrint("Response sent.")
}

func main() {

	server, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	DebugPrint("Starting server on port: " + port)

	for {
		conn, err := server.Accept()

		if err != nil {
			panic(err)
		}

		go handlePingRequest(conn)
	}
}
