package middleware

import (
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
)

func LimitRequests() gin.HandlerFunc {
	// Allow 5 requests per 10 seconds (0.5 req/sec)
	limiter := tollbooth.NewLimiter(0.5, nil)
	limiter.SetTokenBucketExpirationTTL(time.Hour)

	limiter.SetMessage("Too many requests — please slow down.")
	limiter.SetMessageContentType("text/plain; charset=utf-8")

	return tollbooth_gin.LimitHandler(limiter)
}
