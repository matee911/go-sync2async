package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
	"encoding/json"
	"io/ioutil"
	"logging"
	"judge"
)

type Request struct {
	resultChan    chan string
	transactionId string
}

type Config map[string]interface{}
var (
	config Config
)

func parseArguments() (config_path string) {
	flag.StringVar(&config_path, "config", "sync2async.json", "Path to configuration file")
	flag.Parse()
	return
}

func loadConfig(config_path string, fail bool) {
	file, err := ioutil.ReadFile(config_path)
	if err != nil {
		log.Println("open config: ", err)
		if fail {
			os.Exit(1)
		}
	}

	var temp Config
	if err = json.Unmarshal(file, &temp); err != nil {
		log.Println("parse config: ", err)
		if fail {
			os.Exit(1)
		}
	}
	
	config = temp
}


func init() {
	config_path := parseArguments()
	loadConfig(config_path, true)
}

func main() {
	mapping := make(map[string]*Request)
	connection, _ := net.Dial("tcp", config["drm_addr"].(string))
	defer connection.Close()

	http_handler := func(res http.ResponseWriter, req *http.Request) {
		defer logging.HttpRequest(time.Now(), req)
		req.ParseForm()
		transactionId := req.Form.Get("transaction_id")

		// Ask judge...
		judge.AskForPermission()
		// If OK, hit DRM server
		// in other case, just return some error

		request := Request{resultChan: make(chan string), transactionId: transactionId}
		mapping[transactionId] = &request

		go func(request *Request) {
			time.Sleep(1 * time.Second)
			request.resultChan <- request.transactionId
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
