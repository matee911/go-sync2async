package logging

import (
	"log"
	"net/http"
	"time"
)

func HttpRequest(start time.Time, request *http.Request) {
	log.Printf("\"%s %s %s\" \"%s\" %s",
		request.Method,
		request.RequestURI,
		request.Proto,
		request.UserAgent(),
		time.Since(start),
	)
}
