package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/matee911/go-sync2async/cfg"
	"github.com/matee911/go-sync2async/db"
	"github.com/matee911/go-sync2async/dvs"
	"github.com/matee911/go-sync2async/judge"
	"github.com/matee911/go-sync2async/logging"
	"github.com/matee911/go-sync2async/transaction"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
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

func Ping(connection *net.Conn, config *cfg.Config) {
	//_, err := connection.Write(dvs.NoCommand(1, 1, 1, 1))
	//if err != nil {
	//	log.Printf("DVS Conn: %v", err)
	//}
}

func CallDVS(request *Request) {
	time.Sleep(1 * time.Second) // Change to real connection
	request.resultChan <- strconv.Itoa(request.TransactionId)
}

func init() {
	var err error
	configPath, dropSequence := parseArguments()
	config.ReadFromJson(configPath)

	// Prepare DB
	if dbConn, err = db.Connect(&config); err != nil {
		os.Exit(1)
	} else {
		var relkind string
		if err := dbConn.QueryRow(db.CheckSeq, transaction.SequenceName).Scan(&relkind); err != nil {
			// no rows? thats actually good
			if err.Error() == "sql: no rows in result set" {
				if _, err := dbConn.Exec(fmt.Sprintf(db.CreateSeq, transaction.SequenceName)); err != nil {
					log.Printf("Last query: %s", db.CreateSeq)
					log.Print(err.Error())
					os.Exit(1)
				}
			} else {
				log.Fatalln(err.Error())
			}
		} else {
			if dropSequence {
				if _, err := dbConn.Exec(fmt.Sprintf(db.DropSeq, transaction.SequenceName)); err != nil {
					log.Printf("Last query: %s", db.CreateSeq)
					log.Print(err.Error())
					os.Exit(1)
				} else {
					log.Print("Dropped.")
					os.Exit(0)
				}
			}
		}
	}
}

func main() {
	mapping := make(map[int]*Request)
	connection, err := net.Dial("tcp", config.DVS_Addr)
	if err != nil {
		log.Println(err.Error())
	}
	defer connection.Close()

	// TODO(m): ladowanie czasu tickera z konfigu i castowanie
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()
	go func() {
		for t := range ticker.C {
			Ping(&connection, &config)
			log.Printf("Tick: %v", t)
		}
	}()

	sync_http_handler := func(res http.ResponseWriter, req *http.Request) {
		defer logging.HttpRequest(time.Now(), req)
		req.ParseForm()
		//transactionId := req.Form.Get("transaction_id")

		judge.AskForPermission("ala", &config)

		transactionId := transaction.GetId(dbConn)
		request := Request{resultChan: make(chan string), TransactionId: transactionId}
		mapping[transactionId] = &request

		go CallDVS(&request)

		res.Header().Set("Content-Type", "text/plain")
		res.Header().Set("Server", userAgent)
		select {
		case r := <-request.resultChan:
			io.WriteString(res, r)
		// TODO(m): ladowanie czasu z konfigu i castowanie
		case <-time.After(5 * time.Second):
			io.WriteString(res, "zepsute")
		}

	}

	ping_http_handler := func(res http.ResponseWriter, req *http.Request) {
		defer logging.HttpRequest(time.Now(), req)

		transactionId := transaction.GetId(dbConn)
		request := Request{resultChan: make(chan string), TransactionId: transactionId}
		mapping[transactionId] = &request

		// change to CallDVS
		go func(request *Request) {
			connection.Write(dvs.NoCommand(1, 1, 1, 1))
			request.resultChan <- strconv.Itoa(request.TransactionId)
		}(&request)

		res.Header().Set("Content-Type", "text/plain")
		res.Header().Set("Server", userAgent)
		select {
		case r := <-request.resultChan:
			io.WriteString(res, r)
		// TODO(m): ladowanie czasu z konfigu i castowanie
		case <-time.After(5 * time.Second):
			io.WriteString(res, "zepsute")
		}

	}

	http.HandleFunc("/sync", sync_http_handler)
	http.HandleFunc("/ping", ping_http_handler)
	addr := fmt.Sprintf(":%v", config.Port)
	log.Printf("Listening on port %v", config.Port)
	log.Fatal(http.ListenAndServe(addr, nil))
	os.Exit(1)
}
