package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type Msg struct {
	Msg string
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

func index(res http.ResponseWriter, req *http.Request) {
	defer log_request(time.Now(), req)
	res.Header().Set("Content-Type", "text/plain")
	io.WriteString(res, "Dzień dobry")
}

func sync(res http.ResponseWriter, req *http.Request) {
	defer log_request(time.Now(), req)
	res.Header().Set("Content-Type", "application/json")
	io.WriteString(res, "{}")
}

func main() {
	msg := make(chan Msg, 10)

	connection, err := net.Dial("tcp", "", "127.0.0.1:9001")
	defer connection.Close()

	http_handler := func(res http.ResponseWriter, req *http.Request) {
		defer log_request(time.Now(), req)
		msg <- Msg{Msg: "req"}
		res.Header().Set("Content-Type", "application/json")
		io.WriteString(res, "OK")
	}

	go func() {
		for {
			recv_msg := <-msg
			connection.Write("abc")
			log.Printf("%v", recv_msg.Msg)
		}
	}()

	/*    http.HandleFunc("/", index)*/
	http.HandleFunc("/sync", http_handler)

	log.Printf("Listening on 0.0.0.0:9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
	os.Exit(1)
}
