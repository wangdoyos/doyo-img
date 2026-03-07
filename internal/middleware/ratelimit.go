package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/wangdoyos/doyo-img/internal/config"
)

// ipLimiter 单个 IP 的限流器及最后活跃时间
type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter 基于 IP 的请求限流器
type RateLimiter struct {
	mu       sync.Mutex
	limiters map[string]*ipLimiter
	rate     rate.Limit
	burst    int
}

// NewRateLimiter 创建 IP 限流器实例
func NewRateLimiter(cfg *config.RateLimitConfig) *RateLimiter {
	rl := &RateLimiter{
		limiters: make(map[string]*ipLimiter),
		rate:     rate.Limit(float64(cfg.RequestsPerMinute) / 60.0),
		burst:    cfg.Burst,
	}

	// 启动后台协程，定期清理过期的 IP 记录
	go rl.cleanup()

	return rl
}

// getLimiter 获取指定 IP 的限流器，不存在则创建
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if l, exists := rl.limiters[ip]; exists {
		l.lastSeen = time.Now()
		return l.limiter
	}

	l := rate.NewLimiter(rl.rate, rl.burst)
	rl.limiters[ip] = &ipLimiter{limiter: l, lastSeen: time.Now()}
	return l
}

// cleanup 定期清理超过 10 分钟未活跃的 IP 限流器，防止内存泄漏
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, l := range rl.limiters {
			if time.Since(l.lastSeen) > 10*time.Minute {
				delete(rl.limiters, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Middleware 返回限流中间件，超出限制时返回 429
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"data":    nil,
				"message": "rate limit exceeded, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
