package oauth

import (
	"FullStackOfYear/backend/internal/errors"
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

const (
	stateLength    = 32
	stateExpiry    = 10 * time.Minute
	stateKeyPrefix = "oauth:state:"
)

type StateManager struct {
	redis *redis.Client
}

func NewStateManager(rdb *redis.Client) *StateManager {
	return &StateManager{redis: rdb}
}

// GenerateState generates and stores a secure state parameter
func (sm *StateManager) GenerateState(ctx context.Context) (string, error) {
	bytes := make([]byte, stateLength)
	if _, err := rand.Read(bytes); err != nil {
		log.Printf("Failed to generate random state: %v", err)
		return "", errors.NewAppError(errors.InternalError, "Failed to generate state")
	}

	state := base64.URLEncoding.EncodeToString(bytes)
	key := stateKeyPrefix + state

	err := sm.redis.Set(ctx, key, 1, stateExpiry).Err()
	if err != nil {
		log.Printf("Failed to store state in Redis: %v", err)
		return "", errors.NewAppError(errors.InternalError, "Failed to store state")
	}

	return state, nil
}

// ValidateState validates and consumes the state parameter
func (sm *StateManager) ValidateState(ctx context.Context, state string) bool {
	if state == "" {
		log.Print("Empty state parameter received")
		return false
	}

	// Basic format validation
	if len(state) != 44 { // base64 encoded 32 bytes
		log.Printf("Invalid state length: %d", len(state))
		return false
	}

	key := stateKeyPrefix + state
	result := sm.redis.Del(ctx, key).Val()
	if result == 0 {
		log.Printf("State not found or expired: %s", state)
	}
	return result > 0
}
