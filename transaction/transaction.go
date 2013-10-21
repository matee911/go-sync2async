package transaction

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	SequenceName string = "nproxy_transaction_id_seq"
)

func GetId(db *sql.DB) (int, error) {
	var value int
	err := db.QueryRow(fmt.Sprintf("SELECT nextval('%s')", SequenceName)).Scan(&value)
	if err != nil {
		log.Println(err)
	}
	return value, err
}
