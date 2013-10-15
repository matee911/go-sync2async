package db

import (
	"github.com/matee911/go-sync2async/cfg"
	"fmt"
	"log"
	_ "github.com/lib/pq"
	"database/sql"
)

const (
	CheckSeq string = "SELECT c.relkind FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace WHERE c.relname = $1 AND c.relkind = 'S'"
	CreateSeq string = "CREATE SEQUENCE %s INCREMENT 1 START 1"
	DropSeq string = "DROP SEQUENCE %s"
)

func Connect(config *cfg.Config) (*sql.DB, error) {
	c := *config
	host := c.TransactionDB_Host
	port := c.TransactionDB_Port
	name := c.TransactionDB_Name
	user := c.TransactionDB_User
	password := c.TransactionDB_Password
	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, name)
	log.Print(dbURI)
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}