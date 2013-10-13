package main

import (
	"cfg"
	"flag"
	"fmt"
	"io"
	"judge"
	"log"
	"logging"
	"net"
	"net/http"
	"os"
	"time"
)

type Request struct {
	resultChan    chan string
	TransactionId string
}

var (
	config cfg.Config
)

func parseArguments() (config_path string) {
	flag.StringVar(&config_path, "config", "sync2async.json", "Path to configuration file")
	flag.Parse()
	return
}

func init() {
	config_path := parseArguments()
	cfg.LoadConfig(&config, config_path, true)
}

func main() {
	mapping := make(map[string]*Request)
	// TODO: Check connection err
	connection, _ := net.Dial("tcp", config["drm_addr"].(string))
	defer connection.Close()

	http_handler := func(res http.ResponseWriter, req *http.Request) {
		defer logging.HttpRequest(time.Now(), req)
		req.ParseForm()
		transactionId := req.Form.Get("transaction_id")

		// Ask judge...
		// judge.Evidence{TransactionId: transactionId}
		judge.AskForPermission("ala", &config)
		// If OK, hit DRM server
		// in other case, just return some error

		request := Request{resultChan: make(chan string), TransactionId: transactionId}
		mapping[transactionId] = &request

		go func(request *Request) {
			time.Sleep(1 * time.Second)
			request.resultChan <- request.TransactionId
		}(&request)

		res.Header().Set("Content-Type", "text/plain")
		res.Header().Set("Server", "nagra-proxy")
		select {
		case r := <-request.resultChan:
			io.WriteString(res, r)
		case <-time.After(5 * time.Second):
			io.WriteString(res, "zepsute")
		}

	}

	http.HandleFunc("/sync", http_handler)
	addr := fmt.Sprintf(":%v", config["port"])
	log.Printf("Listening on port %v", config["port"])
	log.Fatal(http.ListenAndServe(addr, nil))
	os.Exit(1)
}
