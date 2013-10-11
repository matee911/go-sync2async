package main

import (
//	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type Request struct {
	resultChan chan string
	transactionId string
}

func log_request(start time.Time, request *http.Request) {
	log.Printf("\"%s %s\" %s \"%s\" %s",
		request.Method,
		request.URL.Path,
		request.Proto,
		request.UserAgent(),
		time.Since(start),
	)
}

func main() {
	mapping := make(map[string]*Request)
	connection, _ := net.Dial("tcp", "127.0.0.1:9001")
	defer connection.Close()


	http_handler := func(res http.ResponseWriter, req *http.Request) {
		defer log_request(time.Now(), req)
		req.ParseForm()
		transactionId := req.Form.Get("transaction_id")

		request := Request{resultChan: make(chan string), transactionId: transactionId}
		mapping[transactionId] = &request

		go func(request *Request) {
			time.Sleep(2 * time.Second)
			request.resultChan <- request.transactionId
		}(&request)

		res.Header().Set("Content-Type", "text/plain")
		res.Header().Set("Server", "nagra-proxy")
		select {
		case res := <- request.resultChan:
			io.WriteString(res, <-request.resultChan)
		case <-time.After(5 * time.Second):
			io.WriteString(res, "zepsute")
		}
		
	}

	http.HandleFunc("/sync", http_handler)
	log.Printf("Listening on 0.0.0.0:9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
	os.Exit(1)
}
