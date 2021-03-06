package judge

import (
	"github.com/matee911/go-sync2async/cfg"
	"io/ioutil"
	"log"
	"net/http"
)

type Evidence struct {
	TransactionId string
}

func AskForPermission(evidence string, config *cfg.Config) bool {
	c := *config
	resp, err := http.Get(c.Judge_Addr)
	if err != nil {
		// handle error
		log.Printf("Error: %v", err)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Body: %v", body)
	}
	defer resp.Body.Close()

	return true
}
