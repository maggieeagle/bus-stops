package models

import (
	"database/sql"
	"fmt"
)

type Arrival struct {
	Time string `db:"arrival_time" json:"time"`
	Direction string `db:"direction_code" json:"direction"`
}

func GetArrivalList(db *sql.DB, stop string, code string, route string) ([]Arrival, error) {
	var clause string
    var args []interface{}
    if (code != "") {
        clause = `stop_code = $1`
        args = append(args, code)
    } else {
        clause = `stop_name = $1`
        args = append(args, stop)
    }
	args = append(args, route)
    query := fmt.Sprintf(`
	WITH matching_stops AS (
		SELECT stop_id
		FROM stops
		WHERE
			%s
	)
	SELECT DISTINCT
		st.arrival_time,
		t.direction_code
	FROM stop_times st
	JOIN trips t ON st.trip_id = t.trip_id
	JOIN routes r ON t.route_id = r.route_id
	WHERE
		st.stop_id IN (SELECT stop_id FROM matching_stops)
		AND r.route_short_name = $2
	ORDER BY
		st.arrival_time;
    `, clause)

	rows, err := db.Query(query, args...)

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var arrivals []Arrival
	for rows.Next() {
		var arrival Arrival
		if err := rows.Scan(&arrival.Time, &arrival.Direction); err != nil {
			return nil, err
		}

		arrivals = append(arrivals, arrival)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(arrivals) == 0 {
		return nil, sql.ErrNoRows
	}

	return arrivals, nil
}
