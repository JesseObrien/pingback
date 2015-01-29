package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/unrolled/render"
	"net/http"
	"os"
)

var configFileName string
var port string
var networkNodes []NetworkNode

func init() {
	flag.StringVar(&configFileName, "config", ".pingback.conf", "--config=.pingback.conf")
	flag.StringVar(&port, "port", "3001", "--port=80")
	flag.Parse()
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

func handleSplash(writer http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{})

	r.HTML(writer, http.StatusOK, "splash", nil)
}

func main() {

	if err := loadNetworkNodes(); err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleSplash)
	mux.HandleFunc("/ping", handlePingRequest)

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":" + port)
	fmt.Println("Listening on " + port)
}
