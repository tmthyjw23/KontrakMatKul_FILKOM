package middlewares

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ConcurrencyLimiter limits concurrent requests to prevent database lock contention
type ConcurrencyLimiter struct {
	semaphore chan struct{}
	logger    *zap.Logger
}

// NewConcurrencyLimiter creates a new concurrency limiter
// maxConcurrent: maximum number of concurrent requests allowed
func NewConcurrencyLimiter(maxConcurrent int, logger *zap.Logger) *ConcurrencyLimiter {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &ConcurrencyLimiter{
		semaphore: make(chan struct{}, maxConcurrent),
		logger:    logger,
	}
}

// Middleware returns a gin middleware function
func (cl *ConcurrencyLimiter) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Try to acquire semaphore with timeout
		select {
		case cl.semaphore <- struct{}{}:
			defer func() { <-cl.semaphore }()
			ctx.Next()

		case <-ctx.Request.Context().Done():
			// Request context cancelled
			cl.logger.Warn("request cancelled while waiting for semaphore")
			ctx.JSON(http.StatusRequestTimeout, gin.H{
				"code":    http.StatusRequestTimeout,
				"status":  "error",
				"message": "request timeout",
			})
			ctx.Abort()

		default:
			// Semaphore full
			cl.logger.Warn("too many concurrent requests, returning 503")
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"code":    http.StatusServiceUnavailable,
				"status":  "error",
				"message": "server busy, please try again later",
			})
			ctx.Abort()
		}
	}
}

// WriteLock is a simple mutex-based write lock for sequential write operations
type WriteLock struct {
	mu     sync.Mutex
	logger *zap.Logger
}

// NewWriteLock creates a new write lock
func NewWriteLock(logger *zap.Logger) *WriteLock {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &WriteLock{
		logger: logger,
	}
}

// Middleware returns a gin middleware that locks writes
func (wl *WriteLock) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Only lock on write operations
		if ctx.Request.Method == http.MethodPost || ctx.Request.Method == http.MethodPut || ctx.Request.Method == http.MethodDelete {
			wl.mu.Lock()
			defer wl.mu.Unlock()
			wl.logger.Debug("write lock acquired", zap.String("method", ctx.Request.Method), zap.String("path", ctx.Request.URL.Path))
		}

		ctx.Next()
	}
}
