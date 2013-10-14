package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"github.com/matee911/go-sync2async/logging"
	"net/http"
	"strconv"
	"time"
)

var port int

func log_request(start time.Time, request *http.Request) {
	log.Printf("\"%s %s\" %s \"%s\" %s",
		request.Method,
		request.URL.RequestURI(),
		request.Proto,
		request.UserAgent(),
		time.Since(start),
	)
}

func judgeHandler(res http.ResponseWriter, req *http.Request) {
	defer logging.HttpRequest(time.Now(), req)
	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Server", "judge/0.1")
	req.ParseForm()
	transactionId, _ := strconv.Atoi(req.Form.Get("transaction_id"))
	if transactionId%2 == 0 {
		io.WriteString(res, "OK")
	} else {
		io.WriteString(res, "FAIL")
	}
}

func parseArguments() {
	flag.IntVar(&port, "port", 3000, "Port on which server will listen")
	flag.Parse()
}

func main() {
	parseArguments()
	handlerFunc := http.HandlerFunc(judgeHandler)
	addr := fmt.Sprintf(":%v", port)
	log.Printf("Listening on port %v...", port)
	log.Fatal(http.ListenAndServe(addr, handlerFunc))
}
