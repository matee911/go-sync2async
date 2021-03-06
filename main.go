package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/matee911/go-sync2async/cfg"
	"github.com/matee911/go-sync2async/db"
	//"github.com/matee911/go-sync2async/judge"
	"github.com/matee911/go-sync2async/logging"
	//"github.com/matee911/go-sync2async/transaction"
	"github.com/matee911/go-sync2async/dvs"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Request struct {
	resultChan    chan string
	TransactionId int
}

var (
	config cfg.Config
	dbConn *sql.DB
)

const (
	userAgent string = "NProxy/0.1a"
)

func parseArguments() (configPath string, dropSequence bool) {
	flag.StringVar(&configPath, "config", "sync2async.json", "Path to configuration file")
	flag.BoolVar(&dropSequence, "dropsequence", false, "Drop TransactionID sequence")
	flag.Parse()
	return
}

func init() {
	var err error
	configPath, dropSequence := parseArguments()
	config.ReadFromJson(configPath)

	// Prepare DB
	dbConn, err := db.Connect(&config)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	prepErr := db.PrepareDb(dbConn, dropSequence)
	if prepErr != nil {
		log.Println(prepErr)
		os.Exit(1)
	} else if (prepErr == nil) && dropSequence {
		os.Exit(0)
	}
}

func licenseHttpHandler(commandsChannel chan dvs.Command) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		defer logging.HttpRequest(time.Now(), req)
		req.ParseForm()

		res.Header().Set("Content-Type", "application/json")
		res.Header().Set("Server", userAgent)

		// NUID (Nagra Unique IDentifier) [int, max 10-digits long]
		address, err := validateAddress(req.PostForm.Get("address"))
		if err != nil {
			http.Error(res, CreateErrorResponse(1, "invalid address", err.Error()).String(), 400)
			return
		}

		// Chipset Type String (please refer to DVS documentation)
		// [text in format: xxxx xxxx xxxx xxxx xx]
		chipset_type_string, err := validateChipset(req.PostForm.Get("chipset_type_string"))
		if err != nil {
			http.Error(res, CreateErrorResponse(1, "invalid chipset", err.Error()).String(), 400)
			return
		}

		// Content identifier [int, max 9-digits long]
		content_id, err := validateContent(req.PostForm.Get("content_id"))
		if err != nil {
			http.Error(res, CreateErrorResponse(1, "invalid content", err.Error()).String(), 400)
			return
		}
		// any binary string data, that will be sent to authorization module
		// for additional authorization (eg. user credentials, token)
		extra := req.PostForm.Get("extra")

		log.Printf("addr: %d chip: %s content: %s extra: %s", address, chipset_type_string, content_id, extra)

		commandsChannel <- dvs.NoCommand()

		//judge.AskForPermission("ala", &config)

		/*
		   errcode - high-level error code
		     400
		     403
		     500
		*/
	}
}

func pingHttpHandler(commandsChannel chan dvs.Command) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		defer logging.HttpRequest(time.Now(), req)
		req.ParseForm()

		res.Header().Set("Content-Type", "application/json")
		res.Header().Set("Server", userAgent)

		commandsChannel <- dvs.NoCommand()
		io.WriteString(res, "Hello")

		//select {
		//case r := <-request.resultChan:
		//	io.WriteString(res, r)
		// TODO(m): ladowanie czasu z konfigu i castowanie
		//case <-time.After(time.Duration(config.Timeout) * time.Second):
		//	io.WriteString(res, "zepsute")
		//}
	}
}

func main() {

	commandsChannel := make(chan dvs.Command, 100)
	master := dvs.Master{SourceIds: config.SourceIds, DestinationId: config.DestinationId, MopPpid: config.MopPpid, Address: config.DVSAddr, CommandsChannel: commandsChannel}
	master.Connect()
	// Close?

	// Heartbeat
	heartbeat := time.NewTicker(time.Second * time.Duration(config.Heartbeat))
	defer heartbeat.Stop()
	go func() {
		for t := range heartbeat.C {
			//Ping(&connection, &config)
			commandsChannel <- dvs.NoCommand()
			log.Printf("Heartbeat: %v", t)
		}
	}()

	// Registered handlers
	http.HandleFunc("/license", licenseHttpHandler(commandsChannel))
	http.HandleFunc("/ping", pingHttpHandler(commandsChannel))

	addr := fmt.Sprintf(":%v", config.Port)
	log.Printf("Listening on port %v", config.Port)
	log.Fatal(http.ListenAndServe(addr, nil))
	os.Exit(1)
}
