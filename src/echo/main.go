package main

import (
	"flag" // for command line arguments
	"fmt"
	"io/ioutil" // for easy reading of request body
	"log"
	"math/rand"
	"net/http"
	"time" // for sleep & duration
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	sleep := rand.Intn(50) * 100

	log.Printf("%v %v (sleep: %v)", r.Method, r.URL.Path, sleep)
	time.Sleep(time.Duration(sleep) * time.Millisecond)

	fmt.Fprintf(w, "%s", body)
}

func parseArguments() (port int) {
	flag.IntVar(&port, "port", 3000, "Port on which server will listen")
	flag.Parse()
	return
}

func main() {
	rand.Seed(time.Now().UnixNano())

	port := parseArguments()
	addr := fmt.Sprintf(":%v", port)
	handlerFunc := http.HandlerFunc(echoHandler)

	log.Printf("Listening on port %v...", port)
	log.Fatal(http.ListenAndServe(addr, handlerFunc))
}
