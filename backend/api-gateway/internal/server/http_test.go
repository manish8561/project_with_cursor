package server

import (
	"testing"
	"time"

	"api-gateway/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
)

func TestNewHTTPServer(t *testing.T) {
	cfg := &conf.Server{Http: &conf.Server_HTTP{Network: "tcp", Addr: ":0", Timeout: durationpb.New(5 * time.Second)}}
	server := NewHTTPServer(cfg, nil, log.DefaultLogger)
	require.NotNil(t, server)
}

func TestNewHTTPServer_UsesOptionalHTTPConfig(t *testing.T) {
	cfg := &conf.Server{Http: &conf.Server_HTTP{Network: "tcp"}}
	server := NewHTTPServer(cfg, nil, log.DefaultLogger)
	require.NotNil(t, server)
}

func TestNewHTTPServer_EmptyConfigStillConstructs(t *testing.T) {
	server := NewHTTPServer(&conf.Server{}, nil, log.DefaultLogger)
	require.NotNil(t, server)
}

// Note: Rate limiter and clientIP tests have been moved to the router module
// See internal/router/router_test.go for those tests
