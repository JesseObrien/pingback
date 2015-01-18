package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var configFileName string
var networkNodes []NetworkNode

func init() {
	flag.StringVar(&configFileName, "config", ".pingback.conf", "--config=.pingback.conf")
}

func loadNetworkNodes() error {
	fmt.Println("Loading network nodes from: " + configFileName)
	file, err := os.Open(configFileName)

	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)

	return decoder.Decode(&networkNodes)
}

func main() {

	if err := loadNetworkNodes(); err != nil {
		panic(err)
	}

	http.HandleFunc("/ping", handlePingRequest)
	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
