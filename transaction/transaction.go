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
  var value interface{}
	err := db.QueryRow(fmt.Sprintf("SELECT nextval('%s')", SequenceName)).Scan(&value)
  if err != nil {
    log.Println(err)
    return 0, err
  } else {
    return value.(int), nil
  }
}
