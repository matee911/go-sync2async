package transaction

import (
	"database/sql"
)

const (
	SequenceName string = "nproxy_transaction_id_seq"
)

func GetId(db *sql.DB) int {
	db.Query("SELECT nextval('"+SequenceName+"')")
	return 1
}