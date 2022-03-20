package app

import "context"

type DatabaseRouteStorage struct {
	baseDB *BaseDB
}

func (dbStorage *DatabaseRouteStorage) ShortRoute(fullRoute string) (int, error) {
	var id int
	statement := "INSERT INTO shortened_urls (original_url, user_id) VALUES ($1, $2) RETURNING id;"
	row := dbStorage.baseDB.Connection.QueryRow(context.Background(), statement, fullRoute, 0)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dbStorage *DatabaseRouteStorage) GetRouteByID(id int) (string, error) {
	var route string
	statement := "SELECT original_url FROM shortened_urls WHERE id=$1;"
	row := dbStorage.baseDB.Connection.QueryRow(context.Background(), statement, id)
	err := row.Scan(&route)
	if err != nil {
		return "", err
	}

	return route, nil
}
