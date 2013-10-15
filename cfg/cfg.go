package cfg

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

//type Config map[string]interface{}
type Config struct {
	Port int
	Timeout int
	DVS_Addr string
	Judge_Addr string
	TransactionDB_Host string
	TransactionDB_Port int
	TransactionDB_Name string
	TransactionDB_User string
	TransactionDB_Password string
}

func (config *Config) ReadFromJson(config_path string) {
	def := Config{
		Port: 9000,
		Timeout: 5,
		DVS_Addr: "localhost:3000",
		Judge_Addr: "http://localhost:4000/",
		TransactionDB_Host: "localhost",
		TransactionDB_Port: 5432,
		TransactionDB_Name: "nproxy",
		TransactionDB_User: "nproxy",
		TransactionDB_Password: "yxorpn",
	}
	
	
	file, err := ioutil.ReadFile(config_path)
	if err != nil {
		log.Println("open config: ", err)
	} else {
		if err = json.Unmarshal(file, &def); err != nil {
			log.Println("parse config: ", err)
		}
	}
	*config = def
}
