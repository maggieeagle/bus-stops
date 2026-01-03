package handlers

import (
	"server/models"
	"sort"
	"time"
	"strconv"

	"github.com/climech/naturalsort"
	"github.com/kellydunn/golang-geo"
)

func sortRoutes(routes []models.Route) {
	sort.SliceStable(routes, func(i, j int) bool {
		return naturalsort.Compare(routes[i].ShortName, routes[j].ShortName)
	})
}

func findNextArrivals(arrivals []models.Arrival, limit int) []models.Arrival {
	loc, err := time.LoadLocation("Europe/Tallinn")
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now().In(loc) 
	nowStr := now.Format("15:04:05")

	var result []models.Arrival

	// take arrivals after now
	for _, a := range arrivals {
		if a.Time >= nowStr {
			result = append(result, a)
			if len(result) == limit {
				return result
			}
		}
	}

	// wrap around midnight to reach limit if needed
	for _, a := range arrivals {
		if a.Time < nowStr {
			result = append(result, a)
			if len(result) == limit {
				break
			}
		}
	}

	return result
}

func findNearestStop(coords Coordinates, stops []models.Stop) (models.Stop, error) {
    // convert center coordinates to float64
    lat, err := strconv.ParseFloat(coords.Latitude, 64)
    if err != nil {
		return models.Stop{}, err
    }
    lon, err := strconv.ParseFloat(coords.Longitude, 64)
    if err != nil {
		return models.Stop{}, err
    }

    center := geo.NewPoint(lat, lon)
    var minDist float64
    var nearest models.Stop

    for i, s := range stops {
        // convert stop coordinates to float64
        stopLat, err := strconv.ParseFloat(s.Latitude, 64)
        if err != nil {
            continue
        }
        stopLon, err := strconv.ParseFloat(s.Longitude, 64)
        if err != nil {
            continue
        }

        stopPoint := geo.NewPoint(stopLat, stopLon)

        // find the great circle distance
        dist := center.GreatCircleDistance(stopPoint)
        if i == 0 || dist < minDist {
            minDist = dist
            nearest = s
        }
    }

    return nearest, nil
}