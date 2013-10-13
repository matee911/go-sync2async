package judge

import (
	"cfg"
	"io/ioutil"
	"log"
	"net/http"
)

type Evidence struct {
	TransactionId string
}

func AskForPermission(evidence string, config *cfg.Config) bool {
	c := *config
	resp, err := http.Get(c["judge_addr"].(string))
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
