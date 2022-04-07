package app

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
)

type BaseDB struct {
	isPrepared bool
	Connection *pgx.Conn
}

func (db *BaseDB) SetupConnection(dsn string) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return
	}
	db.Connection = conn
	_, err = db.Connection.Exec(context.Background(), "create table shortened_urls (id serial constraint table_name_pk primary key, original_url varchar(2048) not null, user_id int not null, is_deleted boolean default false not null);")
	if err != nil {
		db.isPrepared = false
	}
	_, err = db.Connection.Exec(context.Background(), "create unique index if not exists shortened_urls_original_url_uindex on shortened_urls (original_url);")
	if err != nil {
		db.isPrepared = false
	}
	db.isPrepared = true
}

func (db *BaseDB) IsConnectionEstablished() bool {
	return db.Connection != nil && db.isPrepared
}

func (db *BaseDB) CloseConnection() {
	err := db.Connection.Close(context.Background())
	if err != nil {
		return
	}
}

func (db *BaseDB) Ping() error {
	if db.Connection == nil {
		return errors.New("no connection established")
	}
	err := db.Connection.Ping(context.Background())
	if err != nil {
		return err
	}

	return nil
}
