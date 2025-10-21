package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/lekovv/go-web-mvp/service"
)

type TokenBlacklistCleaner struct {
	authService service.AuthServiceInterface
	interval    time.Duration
}

func NewTokenBlacklistCleaner(authService service.AuthServiceInterface, interval time.Duration) *TokenBlacklistCleaner {
	return &TokenBlacklistCleaner{
		authService: authService,
		interval:    interval,
	}
}

func (j *TokenBlacklistCleaner) Start(ctx context.Context) {
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	log.Printf("Token blacklist cleaner started with interval %v", j.interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("Token blacklist cleaner shutting down")
			return
		case <-ticker.C:
			if err := j.authService.DeleteExpiredTokens(ctx); err != nil {
				log.Printf("Token blacklist cleaner failed to delete expired tokens: %v", err)
			} else {
				log.Printf("Token blacklist cleaner successfully deleted expired tokens")
			}
		}
	}
}
