package server

import (
	v1 "api-gateway/api/helloworld/v1"
	"api-gateway/internal/conf"
	"api-gateway/internal/service"

	"io"
	nethttp "net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

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
	srv := kratoshttp.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)

	// Serve Swagger UI static files and openapi.yaml
	srv.HandleFunc("/swagger/", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if r.URL.Path == "/swagger/" || r.URL.Path == "/swagger/index.html" {
			nethttp.ServeFile(w, r, "swagger-ui/index.html")
			return
		}
		if r.URL.Path == "/swagger/openapi.yaml" {
			nethttp.ServeFile(w, r, "openapi.yaml")
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

	// Proxy /auth/* to auth-service
	srv.HandleFunc("/auth/", func(w nethttp.ResponseWriter, r *nethttp.Request) {
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
	})

	// Proxy /user/* to user-service
	srv.HandleFunc("/user/", func(w nethttp.ResponseWriter, r *nethttp.Request) {
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
	})

	return srv
}
