package server

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSimpleTokenBucketLimiter(t *testing.T) {
	limiter := newLimiter(20, 2)

	require.True(t, limiter.allow("192.0.2.1"))
	require.True(t, limiter.allow("192.0.2.1"))
	require.False(t, limiter.allow("192.0.2.1"))

	time.Sleep(60 * time.Millisecond)
	require.True(t, limiter.allow("192.0.2.1"))
}

func TestSimpleTokenBucketLimiterDisabled(t *testing.T) {
	for _, tc := range []struct {
		name  string
		rps   float64
		burst int
	}{
		{name: "zero rps", rps: 0, burst: 1},
		{name: "zero burst", rps: 1, burst: 0},
		{name: "negative rps", rps: -1, burst: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			limiter := newLimiter(tc.rps, tc.burst)
			require.True(t, limiter.allow("192.0.2.1"))
		})
	}
}

func TestClientIP(t *testing.T) {
	tests := []struct {
		name       string
		forwarded  string
		realIP     string
		remoteAddr string
		want       string
	}{
		{name: "forwarded chain", forwarded: " 192.0.2.10, 198.51.100.2", remoteAddr: "10.0.0.1:1234", want: "192.0.2.10"},
		{name: "real ip", realIP: " 192.0.2.11 ", remoteAddr: "10.0.0.1:1234", want: "192.0.2.11"},
		{name: "ipv6 peer", remoteAddr: "[2001:db8::1]:1234", want: "2001:db8::1"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/", nil)
			request.RemoteAddr = test.remoteAddr
			request.Header.Set("X-Forwarded-For", test.forwarded)
			request.Header.Set("X-Real-IP", test.realIP)
			require.Equal(t, test.want, clientIP(request))
		})
	}
}
