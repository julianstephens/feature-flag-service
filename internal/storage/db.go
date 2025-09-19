package storage

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
)

var (
	dbConn *pgx.Conn
	once   sync.Once
)

func OpenDatabaseConnection(connectionString string) error {
	var err error
	once.Do(func() {
		conn, e := pgx.Connect(context.Background(), connectionString)
		if e != nil {
			err = e
		}
		dbConn = conn
	})
	return err
}

func CloseDatabaseConnection() error {
	if dbConn != nil {
		return dbConn.Close(context.Background())
	}
	return nil
}
