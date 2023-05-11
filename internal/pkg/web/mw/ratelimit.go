package mw

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

// RateLimit returns middleware to limit requests.
func RateLimit(requestsCount int, d time.Duration) gin.HandlerFunc {
	rl := rate.NewLimiter(rate.Every(d), requestsCount)

	return func(context *gin.Context) {
		if !rl.Allow() {
			context.AbortWithStatus(http.StatusTooManyRequests)
		}
	}
}
