package cfg

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Port int
  // DVS_Addr address of DVS server [localhost:3000]
	DVSAddr string
  // Judge HTTP Address [http://localhost:4000/]
	Judge_Addr string
  // TransactionDB configuration [nproxy:yxorpn@localhost:5432/nproxy]
	TransactionDB_Host string
	TransactionDB_Port int
	TransactionDB_Name string
	TransactionDB_User string
	TransactionDB_Password string
  // Heartbeat - ping/no_command interval [30]
  Heartbeat int
  // Timeout for HTTP connection to NProxy [10]
  Timeout int
  SourceIds []int
  DestinationId int
  MopPpid int
}

func (config *Config) ReadFromJson(configPath string) {
  // Default values
	def := Config{
		Port: 9000,
		DVSAddr: "localhost:3000",
		Judge_Addr: "http://localhost:4000/",
		TransactionDB_Host: "localhost",
		TransactionDB_Port: 5432,
		TransactionDB_Name: "nproxy",
		TransactionDB_User: "nproxy",
		TransactionDB_Password: "yxorpn",
    Heartbeat: 30,
    Timeout: 10,
    SourceIds: []int{1,2,3,4},
	}
	
	
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println("open config: ", err)
	} else {
		if err = json.Unmarshal(file, &def); err != nil {
			log.Println("parse config: ", err)
		}
	}
	*config = def
}
