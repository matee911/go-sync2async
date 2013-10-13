package cfg

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config map[string]interface{}

func LoadConfig(config *Config, config_path string, fail bool) {
	file, err := ioutil.ReadFile(config_path)
	if err != nil {
		log.Println("open config: ", err)
		if fail {
			os.Exit(1)
		}
	}

	var temp Config
	if err = json.Unmarshal(file, &temp); err != nil {
		log.Println("parse config: ", err)
		if fail {
			os.Exit(1)
		}
	}
	log.Printf("temp %v", temp)
	*config = temp
}
