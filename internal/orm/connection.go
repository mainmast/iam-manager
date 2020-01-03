package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
)

var instance *sql.DB
var err error

var maxDBConnections, _ = strconv.Atoi(os.Getenv("IAM_MAX_DB_CONNECTIONS"))
var maxIdleDBConnections, _ = strconv.Atoi(os.Getenv("MAX_IDLE_DB_CONNECTIONS"))

//GetDBInstance ...
func GetDBInstance() (*sql.Conn, error) {

	if instance == nil {

		instance, err = sql.Open("postgres", os.Getenv("IAM_DB_URI"))

		if err != nil {
			return nil, err
		}

		instance.SetMaxOpenConns(maxDBConnections)
		instance.SetMaxIdleConns(maxIdleDBConnections)
	}

	con, err := instance.Conn(context.Background())

	if err != nil {

		fmt.Println("Error getting connection from pool")
		return nil, errors.New("Error getting conneciton from pool")
	}

	return con, nil
}
