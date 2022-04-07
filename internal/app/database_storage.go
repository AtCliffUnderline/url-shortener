package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"strconv"
)

var ErrRouteAlreadyShortened = errors.New("route already shortened")

type Route struct {
	Route     string
	isDeleted bool
}

type DatabaseRouteStorage struct {
	baseDB *BaseDB
}

func (dbStorage *DatabaseRouteStorage) SaveBatchRoutes(routes []BatchURLShortenerRequest, userID int) ([]BatchURLShortenerURLIDs, error) {
	var result []BatchURLShortenerURLIDs
	tx, err := dbStorage.baseDB.Connection.BeginTx(context.Background(), pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return nil, err
	}

	var maxID int
	statement := "SELECT MAX(id) FROM shortened_urls;"
	row := tx.QueryRow(context.Background(), statement)
	err = row.Scan(&maxID)
	if err != nil {
		err := tx.Rollback(context.Background())
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	vals := []interface{}{}
	statement = "INSERT INTO shortened_urls (id, original_url, user_id) VALUES "
	for cnt, URLToShort := range routes {
		maxID = maxID + 1
		statement += fmt.Sprintf("(%s, %s, %s),", "$"+strconv.Itoa(cnt*3+1), "$"+strconv.Itoa(cnt*3+2), "$"+strconv.Itoa(cnt*3+3))
		vals = append(vals, maxID, URLToShort.URL, userID)
		result = append(result, BatchURLShortenerURLIDs{ID: maxID, CorrelationID: URLToShort.ID, OriginalURL: URLToShort.URL})
	}
	statement = statement[0 : len(statement)-1]
	prepared, err := tx.Prepare(context.Background(), "insert", statement)
	if err != nil {
		tx.Rollback(context.Background())
		return nil, err
	}
	_, err = tx.Exec(context.Background(), prepared.SQL, vals...)
	if err != nil {
		tx.Rollback(context.Background())
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		tx.Rollback(context.Background())
		return nil, err
	}

	return result, nil
}

func (dbStorage *DatabaseRouteStorage) ShortRoute(fullRoute string, userID int) (int, error) {
	dbStorage.baseDB.Connection.Exec(context.Background(), "SELECT setval('the_primary_key_sequence', (SELECT MAX(id) FROM shortened_urls)+1);")
	res, err := dbStorage.isRouteAlreadyPresented(fullRoute)
	if err == nil {
		return res, ErrRouteAlreadyShortened
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return 0, err
	}
	tx, err := dbStorage.baseDB.Connection.BeginTx(context.Background(), pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return 0, err
	}
	var id int
	statement := "INSERT INTO shortened_urls (id, original_url, user_id) VALUES ((SELECT coalesce(MAX(id),1) FROM shortened_urls)+1, $1, $2) RETURNING id;"
	row := tx.QueryRow(context.Background(), statement, fullRoute, userID)
	err = row.Scan(&id)
	if err != nil {
		_ = tx.Rollback(context.Background())
		return 0, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		_ = tx.Rollback(context.Background())
		return 0, err
	}

	return id, nil
}

func (dbStorage *DatabaseRouteStorage) DeleteRouteByIDForUser(routeID int, userID int) error {
	statement := "UPDATE shortened_urls SET is_deleted=TRUE WHERE id=$1 and user_id=$2;"
	_, err := dbStorage.baseDB.Connection.Exec(context.Background(), statement, routeID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (dbStorage *DatabaseRouteStorage) GetRouteByID(id int) (string, error) {
	route := Route{}
	statement := "SELECT original_url, is_deleted FROM shortened_urls WHERE id=$1;"
	row := dbStorage.baseDB.Connection.QueryRow(context.Background(), statement, id)
	err := row.Scan(&route.Route, &route.isDeleted)
	if err != nil {
		return "", err
	}
	if route.isDeleted {
		return "", ErrRouteDeleted
	}

	return route.Route, nil
}

func (dbStorage *DatabaseRouteStorage) isRouteAlreadyPresented(route string) (int, error) {
	var id int
	statement := "SELECT id FROM shortened_urls WHERE original_url=$1;"
	row := dbStorage.baseDB.Connection.QueryRow(context.Background(), statement, route)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
