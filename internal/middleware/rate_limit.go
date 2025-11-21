package middleware

import (
	"net/http"
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
)

func LimitRequests() gin.HandlerFunc {
	// 0.5 req/sec = 5 requests per 10 seconds
	limiter := tollbooth.NewLimiter(0.5, nil)
	limiter.SetTokenBucketExpirationTTL(time.Hour)

	limiter.SetMessage("Too many requests — please slow down.")
	limiter.SetMessageContentType("text/plain; charset=utf-8")

	// Wrap tollbooth handler so we can skip OPTIONS
	tollboothHandler := tollbooth_gin.LimitHandler(limiter)

	return func(c *gin.Context) {
		// Don't rate-limit preflight requests
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		// Apply tollbooth to everything else
		tollboothHandler(c)
	}
}
