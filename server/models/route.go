package models

import (
	"database/sql"
    "fmt"
)

type Route struct {
	ID string `db:"route_id" json:"id"`
	ShortName string `db:"route_short_name" json:"short_name"`
	LongName string `db:"route_long_name" json:"long_name"`
}

func GetRouteList(db *sql.DB, stop string, code string) ([]Route, error) {
    var clause string
    var arg string
    if (code != "") {
        clause = `stop_code = $1`
        arg = code
    } else {
        clause = `stop_name = $1`
        arg = stop
    }
    query := fmt.Sprintf(`
    WITH stops_around AS (
    SELECT stop_id
    FROM stops
    WHERE
        %s
    )
    SELECT r.route_id, r.route_short_name, r.route_long_name
    FROM routes r
    WHERE r.route_id IN (
        SELECT t.route_id
        FROM trips t
        INNER JOIN stop_times st ON t.trip_id = st.trip_id
        WHERE st.stop_id IN (SELECT stop_id FROM stops_around)
    );
    `, clause)

    rows, err := db.Query(query, arg)

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var routes []Route
    for rows.Next() {
        var r Route
		var shortName sql.NullString
		var longName sql.NullString

        if err := rows.Scan(&r.ID, &shortName, &longName); err != nil {
            return nil, err
        }

		r.ShortName = shortName.String
		r.LongName = longName.String

        routes = append(routes, r)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

	if len(routes) == 0 {
        return nil, sql.ErrNoRows
    }

    return routes, nil
}
