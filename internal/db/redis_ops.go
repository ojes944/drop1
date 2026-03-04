package db

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func SetDriverLocation(driverID string, lat, lng float64) error {
	return Redis.GeoAdd(context.Background(), "drivers", &redis.GeoLocation{
		Name:      driverID,
		Longitude: lng,
		Latitude:  lat,
	}).Err()
}

func GetNearbyDrivers(lat, lng, radius float64) ([]redis.GeoLocation, error) {
	res, err := Redis.GeoRadius(context.Background(), "drivers", lng, lat, &redis.GeoRadiusQuery{
		Radius:    radius,
		Unit:      "m",
		WithCoord: true,
	}).Result()
	return res, err
}

func SetResetToken(userID, token string, expiry time.Duration) error {
	return Redis.Set(context.Background(), "reset:"+userID, token, expiry).Err()
}

func GetResetToken(userID string) (string, error) {
	return Redis.Get(context.Background(), "reset:"+userID).Result()
}
