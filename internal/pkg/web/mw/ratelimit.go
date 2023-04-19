package mw

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

func RateLimitPerSec(requests int) gin.HandlerFunc {
	rl := rate.NewLimiter(rate.Every(time.Second), requests)

	return func(context *gin.Context) {
		if !rl.Allow() {
			context.AbortWithStatus(http.StatusTooManyRequests)
		}
	}
}
