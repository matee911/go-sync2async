package main

import (
	"flag"
	"fmt"
	"github.com/matee911/go-sync2async/cfg"
	"github.com/matee911/go-sync2async/dvs"
	"github.com/matee911/go-sync2async/judge"
	"github.com/matee911/go-sync2async/logging"
	"io"
	"log"
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
	connection, err := net.Dial("tcp", config["drm_addr"].(string))
	if err != nil {
		log.Println(err.Error())
	}
	defer connection.Close()

	sync_http_handler := func(res http.ResponseWriter, req *http.Request) {
		defer logging.HttpRequest(time.Now(), req)
		req.ParseForm()
		transactionId := req.Form.Get("transaction_id")

		judge.AskForPermission("ala", &config)

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

	ping_http_handler := func(res http.ResponseWriter, req *http.Request) {
		defer logging.HttpRequest(time.Now(), req)
		req.ParseForm()
		transactionId := req.Form.Get("transaction_id")

		request := Request{resultChan: make(chan string), TransactionId: transactionId}
		mapping[transactionId] = &request

		go func(request *Request) {
			header := dvs.RootHeader(1, dvs.CmdTypeOther, 1, 1, 1)
			msg := dvs.NoCommand()
			connection.Write(dvs.DeviceIO(fmt.Sprint(header, msg)))
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

	http.HandleFunc("/sync", sync_http_handler)
	http.HandleFunc("/ping", ping_http_handler)
	addr := fmt.Sprintf(":%v", config["port"])
	log.Printf("Listening on port %v", config["port"])
	log.Fatal(http.ListenAndServe(addr, nil))
	os.Exit(1)
}
