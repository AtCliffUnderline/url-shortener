package app

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
)

type BaseDB struct {
	connection *pgx.Conn
}

func (db *BaseDB) SetupConnection(dsn string) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return
	}
	db.connection = conn
}

func (db *BaseDB) CloseConnection() {
	err := db.connection.Close(context.Background())
	if err != nil {
		return
	}
}

func (db *BaseDB) Ping() error {
	if db.connection == nil {
		return errors.New("no connection established")
	}
	_, err := db.connection.Exec(context.Background(), "SELECT 1")
	if err != nil {
		return err
	}

	return nil
}
