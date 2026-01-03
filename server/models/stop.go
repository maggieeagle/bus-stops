package models

import (
	"database/sql"
)

type Stop struct {
	ID int `db:"stop_id" json:"id"`
	Code string `db:"stop_code" json:"code"`
	Name string `db:"stop_name" json:"name"`
	Latitude string `db:"stop_lat" json:"lat"`
	Longitude string `db:"stop_lon" json:"lon"`
	StopArea string `db:"stop_area" json:"region"`
}

func GetStopList(db *sql.DB, region string) ([]Stop, error) {
    rows, err := db.Query("SELECT stop_id, stop_code, stop_name, stop_lat, stop_lon, stop_area FROM stops WHERE stop_area = $1", region)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var stops []Stop
    for rows.Next() {
        var s Stop
		var code sql.NullString
		var name sql.NullString

        if err := rows.Scan(&s.ID, &code, &name, &s.Latitude, &s.Longitude, &s.StopArea); err != nil {
            return nil, err
        }

		s.Code = code.String
		s.Name = name.String

        stops = append(stops, s)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

	if len(stops) == 0 {
        return nil, sql.ErrNoRows
    }

    return stops, nil
}

func GetStopsAll(db *sql.DB) ([]Stop, error) {
    rows, err := db.Query("SELECT stop_id, stop_code, stop_name, stop_lat, stop_lon, stop_area FROM stops")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var stops []Stop
    for rows.Next() {
        var s Stop
		var code sql.NullString
		var name sql.NullString
		var stopArea sql.NullString

        if err := rows.Scan(&s.ID, &code, &name, &s.Latitude, &s.Longitude, &stopArea); err != nil {
            return nil, err
        }

		s.Code = code.String
		s.Name = name.String
		s.StopArea = stopArea.String

        stops = append(stops, s)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

	if len(stops) == 0 {
        return nil, sql.ErrNoRows
    }

    return stops, nil
}

