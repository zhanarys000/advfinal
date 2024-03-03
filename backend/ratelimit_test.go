package backend

import (
	"adv/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"golang.org/x/time/rate"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimiterMiddleware(t *testing.T) {
	limiter := rate.NewLimiter(1, 1)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	middleware.RateLimiterMiddleware(limiter)(c)
	assert.Equal(t, http.StatusOK, c.Writer.Status())
	middleware.RateLimiterMiddleware(limiter)(c)
	assert.Equal(t, http.StatusTooManyRequests, c.Writer.Status())
}
