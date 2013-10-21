package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/matee911/go-sync2async/cfg"
	"github.com/matee911/go-sync2async/transaction"
	"log"
	"os"
)

const (
	CheckSeq  string = "SELECT c.relkind FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace WHERE c.relname = $1 AND c.relkind = 'S'"
	CreateSeq string = "CREATE SEQUENCE %s INCREMENT 1 START 1"
	DropSeq   string = "DROP SEQUENCE %s"
)

// Connect makes connection to database.
// If environment variable NPROXY_DB_URI is provided,
// it will be used instead of the data from configuration file
// or defaults.
func Connect(config *cfg.Config) (*sql.DB, error) {
	var dbUri string

	// defaults or from file
	c := *config
	host := c.TransactionDB_Host
	port := c.TransactionDB_Port
	name := c.TransactionDB_Name
	user := c.TransactionDB_User
	password := c.TransactionDB_Password

	// Is env variable provided?
	envDbUri := os.Getenv("NPROXY_DB_URI")
	if envDbUri != "" {
		dbUri = envDbUri
	} else {
		dbUri = fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, name)
	}

	log.Printf("Connecting to DB: %s", dbUri)
	db, err := sql.Open("postgres", dbUri)
	if err != nil {
		log.Fatal(err)
	} else {
		// BUG(m): It goes here even if network connection is unavailable.
		// Maybe we should detect it calling DNS Resolver?
		log.Println("Connected to DB")
	}
	return db, err
}

func PrepareDb(dbConn *sql.DB, dropSequence bool) (err error) {
	var relkind string

	if err := dbConn.QueryRow(CheckSeq, transaction.SequenceName).Scan(&relkind); err != nil {
		// no rows? thats actually good
		if err.Error() == "sql: no rows in result set" {
			if _, err := dbConn.Exec(fmt.Sprintf(CreateSeq, transaction.SequenceName)); err != nil {
				log.Printf("Last query: %s", CreateSeq)
				log.Println(err.Error())
				return err
			}
		} else {
			log.Printf("Last query: %s", CheckSeq)
			log.Println(err.Error())
			return err
		}
	} else {
		if dropSequence {
			if _, err := dbConn.Exec(fmt.Sprintf(DropSeq, transaction.SequenceName)); err != nil {
				log.Printf("Last query: %s", CreateSeq)
				log.Println(err.Error())
				return err
			} else {
				log.Println("Dropped.")
			}
		}
	}
	return
}
