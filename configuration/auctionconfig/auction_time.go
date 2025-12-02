package auctionconfig

import (
	"fullcycle-auction_go/configuration/logger"
	"os"
	"time"
)

const (
	auctionIntervalEnv     = "AUCTION_INTERVAL"
	defaultAuctionInterval = time.Minute * 5
)

func AuctionDuration() time.Duration {
	auctionInterval := os.Getenv(auctionIntervalEnv)
	if auctionInterval == "" {
		return defaultAuctionInterval
	}

	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		logger.Error("Error parsing auction interval, using default value", err)
		return defaultAuctionInterval
	}

	return duration
}

func AuctionEndTime(start time.Time) time.Time {
	return start.Add(AuctionDuration())
}
