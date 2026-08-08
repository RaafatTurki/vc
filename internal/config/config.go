package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Address          string
	AllowedOrigins   []string
	StunURLs         []string
	TurnURLs         []string
	TurnSharedSecret string
	TurnTTL          time.Duration
	MaxRoomPeers     int
}

func FromEnv() (Config, error) {
	cfg := Config{
		Address:          envOr("SIGNAL_ADDRESS", ":8080"),
		AllowedOrigins:   splitList(os.Getenv("SIGNAL_ALLOWED_ORIGINS")),
		StunURLs:         splitList(os.Getenv("STUN_URLS")),
		TurnURLs:         splitList(os.Getenv("TURN_URLS")),
		TurnSharedSecret: os.Getenv("TURN_SHARED_SECRET"),
		TurnTTL:          30 * time.Minute,
		MaxRoomPeers:     8,
	}

	var err error
	if value := os.Getenv("TURN_TTL"); value != "" {
		cfg.TurnTTL, err = time.ParseDuration(value)
		if err != nil || cfg.TurnTTL <= 0 {
			return Config{}, fmt.Errorf("TURN_TTL must be a positive duration: %q", value)
		}
	}
	if value := os.Getenv("SIGNAL_MAX_ROOM_PEERS"); value != "" {
		cfg.MaxRoomPeers, err = strconv.Atoi(value)
		if err != nil || cfg.MaxRoomPeers < 2 {
			return Config{}, fmt.Errorf("SIGNAL_MAX_ROOM_PEERS must be at least 2: %q", value)
		}
	}
	if len(cfg.TurnURLs) > 0 && cfg.TurnSharedSecret == "" {
		return Config{}, errors.New("TURN_SHARED_SECRET is required when TURN_URLS is configured")
	}
	return cfg, nil
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func splitList(value string) []string {
	var values []string
	for _, item := range strings.Split(value, ",") {
		if item = strings.TrimSpace(item); item != "" {
			values = append(values, item)
		}
	}
	return values
}
