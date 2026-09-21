package server

import (
	v1 "api-gateway/api/helloworld/v1"
	"api-gateway/internal/conf"
	"api-gateway/internal/service"

	"io"
	"net"
	nethttp "net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

// simpleTokenBucketLimiter provides per-IP rate limiting using a token bucket.
type simpleTokenBucketLimiter struct {
	rps     float64
	burst   int
	mu      sync.Mutex
	buckets map[string]*tokenBucket
}

const maxLimiterBuckets = 10000

type tokenBucket struct {
	tokens   float64
	lastUsed time.Time
}

func newLimiter(rps float64, burst int) *simpleTokenBucketLimiter {
	return &simpleTokenBucketLimiter{
		rps:     rps,
		burst:   burst,
		buckets: make(map[string]*tokenBucket),
	}
}

func (l *simpleTokenBucketLimiter) allow(ip string) bool {
	if l == nil || l.rps <= 0 || l.burst <= 0 {
		return true
	}
	now := time.Now().UTC()
	l.mu.Lock()
	defer l.mu.Unlock()
	b, ok := l.buckets[ip]
	if !ok {
		if len(l.buckets) >= maxLimiterBuckets {
			l.evictOldest()
		}
		l.buckets[ip] = &tokenBucket{tokens: float64(l.burst - 1), lastUsed: now}
		return true
	}
	elapsed := now.Sub(b.lastUsed).Seconds()
	b.tokens = minFloat(float64(l.burst), b.tokens+elapsed*l.rps)
	b.lastUsed = now
	if b.tokens >= 1 {
		b.tokens -= 1
		return true
	}
	return false
}

func (l *simpleTokenBucketLimiter) evictOldest() {
	var oldestIP string
	var oldest time.Time
	for ip, bucket := range l.buckets {
		if oldestIP == "" || bucket.lastUsed.Before(oldest) {
			oldestIP = ip
			oldest = bucket.lastUsed
		}
	}
	if oldestIP != "" {
		delete(l.buckets, oldestIP)
	}
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *kratoshttp.Server {
	var opts = []kratoshttp.ServerOption{
		kratoshttp.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, kratoshttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, kratoshttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, kratoshttp.Timeout(c.Http.Timeout.AsDuration()))
	}
	// Configure rate limiter via env vars (defaults: 10 rps, burst 20).
	rps := 10.0
	burst := 20
	if v := os.Getenv("RATE_LIMIT_RPS"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f >= 0 {
			rps = f
		}
	}
	if v := os.Getenv("RATE_LIMIT_BURST"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			burst = n
		}
	}
	limiter := newLimiter(rps, burst)

	srv := kratoshttp.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)

	// Serve Swagger UI static files and openapi.yaml
	srv.HandleFunc("/swagger/openapi.yaml", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		nethttp.ServeFile(w, r, "swagger-ui/openapi.yaml")
	})
	srv.HandleFunc("/swagger/", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if r.URL.Path == "/swagger/" || r.URL.Path == "/swagger/index.html" {
			nethttp.ServeFile(w, r, "swagger-ui/index.html")
			return
		}
		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
	})

	// Health check endpoint
	srv.HandleFunc("/health", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	// Auth service health check endpoint
	srv.HandleFunc("/api/auth/health", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		target := "http://auth-service:8081/health"
		req, err := nethttp.NewRequest(r.Method, target, r.Body)
		if err != nil {
			nethttp.Error(w, "Bad Gateway", nethttp.StatusBadGateway)
			return
		}
		req.Header = r.Header
		resp, err := nethttp.DefaultClient.Do(req)
		if err != nil {
			nethttp.Error(w, "Bad Gateway", nethttp.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	// User service health check endpoint
	srv.HandleFunc("/api/users/health", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		target := "http://user-service:8082/health"
		req, err := nethttp.NewRequest(r.Method, target, r.Body)
		if err != nil {
			nethttp.Error(w, "Bad Gateway", nethttp.StatusBadGateway)
			return
		}
		req.Header = r.Header
		resp, err := nethttp.DefaultClient.Do(req)
		if err != nil {
			nethttp.Error(w, "Bad Gateway", nethttp.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	// Proxy /api/auth/* to auth-service
	srv.HandlePrefix("/api/auth/", nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if !limiter.allow(clientIP(r)) {
			writeRateLimitResponse(w, rps)
			return
		}
		target := "http://auth-service:8081" + r.URL.Path
		req, err := nethttp.NewRequest(r.Method, target, r.Body)
		if err != nil {
			nethttp.Error(w, "Bad Gateway", nethttp.StatusBadGateway)
			return
		}
		req.Header = r.Header
		resp, err := nethttp.DefaultClient.Do(req)
		if err != nil {
			nethttp.Error(w, "Bad Gateway", nethttp.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}))

	// Proxy /api/users/* to user-service
	srv.HandlePrefix("/api/users/", nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if !limiter.allow(clientIP(r)) {
			writeRateLimitResponse(w, rps)
			return
		}
		target := "http://user-service:8082" + r.URL.Path
		req, err := nethttp.NewRequest(r.Method, target, r.Body)
		if err != nil {
			nethttp.Error(w, "Bad Gateway", nethttp.StatusBadGateway)
			return
		}
		req.Header = r.Header
		resp, err := nethttp.DefaultClient.Do(req)
		if err != nil {
			nethttp.Error(w, "Bad Gateway", nethttp.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}))

	return srv
}

func writeRateLimitResponse(w nethttp.ResponseWriter, rps float64) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if rps > 0 {
		retryAfter := int((1 / rps) + 0.999999)
		if retryAfter < 1 {
			retryAfter = 1
		}
		w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
	}
	w.WriteHeader(nethttp.StatusTooManyRequests)
	_, _ = w.Write([]byte("rate limit exceeded"))
}

// clientIP extracts the first forwarded client IP, falling back to the peer address.
func clientIP(r *nethttp.Request) string {
	for _, forwarded := range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		if ip := strings.TrimSpace(forwarded); ip != "" {
			return ip
		}
	}
	if ip := strings.TrimSpace(r.Header.Get("X-Real-IP")); ip != "" {
		return ip
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}
