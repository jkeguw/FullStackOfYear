package geoip

import (
	"FullStackOfYear/backend/internal/database"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type LocationInfo struct {
	Country     string    `json:"country"`
	City        string    `json:"city"`
	LastLoginAt time.Time `json:"lastLoginAt"`
}

const LocationCacheTime = 30 * 24 * time.Hour // 30天

// CheckLocation Check whether the IP address location is abnormal
func CheckLocation(ctx context.Context, userID string, ip string) (bool, error) {
	redis := database.RedisClient
	key := fmt.Sprintf("user:location:%s", userID)

	// Get the user's historical login location
	var locations []LocationInfo
	data, err := redis.Get(ctx, key).Bytes()
	if err == nil {
		if err := json.Unmarshal(data, &locations); err != nil {
			return false, err
		}
	}

	// Get the current IP location
	currentLoc, err := getIPLocation(ip)
	if err != nil {
		return false, err
	}

	// If it is the first time logging in, record the location and allow
	if len(locations) == 0 {
		locations = append(locations, LocationInfo{
			Country:     currentLoc.Country,
			City:        currentLoc.City,
			LastLoginAt: time.Now(),
		})
		saveLocations(ctx, userID, locations)
		return true, nil
	}

	// Check if there is a login record for the same country
	for _, loc := range locations {
		if loc.Country == currentLoc.Country {
			return true, nil
		}
	}

	// New location, additional verification required
	return false, nil
}

// Save location information to Redis
func saveLocations(ctx context.Context, userID string, locations []LocationInfo) error {
	redis := database.RedisClient
	key := fmt.Sprintf("user:location:%s", userID)

	data, err := json.Marshal(locations)
	if err != nil {
		return err
	}

	return redis.Set(ctx, key, data, LocationCacheTime).Err()
}

// IPLocation IP geolocation information
type IPLocation struct {
	Country string
	City    string
}

// getIPLocation Get IP geolocation（需要集成MaxMind GeoIP2数据库）
func getIPLocation(ipStr string) (*IPLocation, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address")
	}

	// TODO: 集成MaxMind GeoIP2数据库
	// 这里先返回模拟数据
	return &IPLocation{
		Country: "CN",
		City:    "Beijing",
	}, nil
}
