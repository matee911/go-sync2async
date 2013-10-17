package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/matee911/go-sync2async/cfg"
	"github.com/matee911/go-sync2async/db"
	//"github.com/matee911/go-sync2async/dvs"
	//"github.com/matee911/go-sync2async/judge"
	"github.com/matee911/go-sync2async/logging"
	//"github.com/matee911/go-sync2async/transaction"
	"encoding/json"
	"errors"
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

func validateAddress(s string) (int, error) {
	if len(s) == 0 {
		return 0, errors.New("address is empty")
	} else if len(s) > 10 {
		return 0, errors.New("address out of range")
	} else if i, err := strconv.Atoi(s); err != nil {
		return i, err
	} else {
		return i, nil
	}
}

func validateChipset(s string) (string, error) {
	if len(s) == 18 {
		return s, errors.New("invalid length of chipset_type_string")
	} else {
		return s, nil
	}
}

func validateContent(s string) (int, error) {
	if len(s) == 0 {
		return 0, errors.New("content is empty")
	} else if len(s) > 9 {
		return 0, errors.New("content out of range")
	} else if i, err := strconv.Atoi(s); err != nil {
		return i, err
	} else {
		return i, nil
	}
}

type ErrorResponse struct {
	Resp ErrRespJSON `json:"resp"`
}

func (r ErrorResponse) String() (s string) {
	body, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(body)
	return
}

type ErrRespJSON struct {
	Status string `json:"status"`
	Ts     int    `json:"ts"`
	//License *LicenseJSON `json:license`
	ErrCode int    `json:"errcode"`
	ErrDesc string `json:"errdesc"`
	ErrText string `json:"err_text"`
}

type LicenseJSON struct {
	Object           string `json:"object"`
	ValidToTimestamp int    `json:"valid_to_timestamp"`
	MetaData         string `json:"metadata"`
}

func licenseHttpHandler(mapping map[int]*Request) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		defer logging.HttpRequest(time.Now(), req)
		req.ParseForm()

		res.Header().Set("Content-Type", "application/json")
		res.Header().Set("Server", userAgent)

		// NUID (Nagra Unique IDentifier) [int, max 10-digits long]
		address, err := validateAddress(req.PostForm.Get("address"))
		if err != nil {
			http.Error(res, err.Error(), 400)
			return
		}

		// Chipset Type String (please refer to DVS documentation)
		// [text in format: xxxx xxxx xxxx xxxx xx]
		chipset_type_string, err := validateChipset(req.PostForm.Get("chipset_type_string"))
		if err != nil {
			http.Error(res, err.Error(), 400)
			return
		}

		// Content identifier [int, max 9-digits long]
		content_id, err := validateContent(req.PostForm.Get("content_id"))
		if err != nil {
			http.Error(res, err.Error(), 400)
			return
		}
		// any binary string data, that will be sent to authorization module
		// for additional authorization (eg. user credentials, token)
		extra := req.PostForm.Get("extra")

		log.Printf("addr: %d chip: %s content: %s extra: %s", address, chipset_type_string, content_id, extra)

		//judge.AskForPermission("ala", &config)

		//request := Request{resultChan: make(chan string), TransactionId: transactionId}
		//mapping[transactionId] = &request
		//go CallDVS(&request)

		/*
		   Response format:
		   {
		     "resp": {
		         "status": [string]: ok | err,
		         "ts" [number]: server time UTC timestamp
		     }
		   }

		   Response contains one object with attribute resp which contains Response object
		   Response object fields:
		   * status [string] = 'err' or 'ok' indicates status of response
		   * ts [int] - current time on server in UTC timestamp format

		   Ok response
		   {
		     "resp": {
		       "status": "ok",
		       "ts": <UTC timestamp>,
		       "license": {
		         "object": "<base64 encoded DVS entitlement response>",
		         "valid_to_timestamp": <UTC timestamp>,
		         "metadata": "<entitlement description from authorization module>"
		       },
		     }
		   }

		   Error response
		   {
		     "resp": {
		       "status": "err",
		       "errcode": <error code>,
		       "errdesc": "<developer message>",
		       "err_text": "<user message>",
		       "ts": <UTC timestamp>
		     }
		   }
		   errcode - high-level error code
		     400
		     403
		     500
		*/

		/*
			select {
			case r := <-request.resultChan:
				io.WriteString(res, r)
			// TODO(m): ladowanie czasu z konfigu i castowanie
			case <-time.After(5 * time.Second):
				io.WriteString(res, "zepsute")
			}
		*/
	}
}

func main() {
	mapping := make(map[int]*Request)
	connection, err := net.Dial("tcp", config.DVS_Addr)
	if err != nil {
		log.Println(err.Error())
	}
	defer connection.Close()

	// Heartbeat
	heartbeat := time.NewTicker(time.Second * time.Duration(config.Heartbeat))
	defer heartbeat.Stop()
	go func() {
		for t := range heartbeat.C {
			//Ping(&connection, &config)
			log.Printf("Heartbeat: %v", t)
		}
	}()

	ping_http_handler := func(res http.ResponseWriter, req *http.Request) {
		defer logging.HttpRequest(time.Now(), req)

		//transactionId := transaction.GetId(dbConn)
		//request := Request{resultChan: make(chan string), TransactionId: transactionId}
		//mapping[transactionId] = &request

		// change to CallDVS
		//go func(request *Request) {
		//	connection.Write(dvs.NoCommand(1, 1, 1, 1))
		//	request.resultChan <- strconv.Itoa(request.TransactionId)
		//}(&request)

		res.Header().Set("Content-Type", "text/plain")
		res.Header().Set("Server", userAgent)
		//select {
		//case r := <-request.resultChan:
		//	io.WriteString(res, r)
		// TODO(m): ladowanie czasu z konfigu i castowanie
		//case <-time.After(time.Duration(config.Timeout) * time.Second):
		//	io.WriteString(res, "zepsute")
		//}
		io.WriteString(res, "io")
	}

	http.HandleFunc("/license", licenseHttpHandler(mapping))
	http.HandleFunc("/ping", ping_http_handler)
	addr := fmt.Sprintf(":%v", config.Port)
	log.Printf("Listening on port %v", config.Port)
	log.Fatal(http.ListenAndServe(addr, nil))
	os.Exit(1)
}
