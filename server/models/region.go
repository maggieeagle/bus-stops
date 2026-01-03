package models

import (
	"database/sql"
)

func GetRegionList(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT DISTINCT stop_area FROM stops WHERE stop_area IS NOT NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var regions []string

	for rows.Next() {
		var region string
		if err := rows.Scan(&region); err != nil {
			return nil, err
		}
		regions = append(regions, region)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return regions, nil
}