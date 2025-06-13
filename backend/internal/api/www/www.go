package www

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/cairon666/vkr-backend/internal/config"
	"github.com/cairon666/vkr-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Router interface {
	Register(router gin.IRouter)
}

func NewHTTPServer(lc fx.Lifecycle, config *config.Config, handler http.Handler) *http.Server {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", config.WWW.Port),
		Handler:           handler,
		ReadHeaderTimeout: time.Second * 10,
		IdleTimeout:       time.Minute * 10,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			fmt.Println("Starting HTTP server at", srv.Addr)
			go srv.Serve(ln)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return srv
}
func NewServeMux(routers []Router, logger logger.ILogger) http.Handler {
	router := gin.New()

	router.Use(
		CustomLogger(logger),
		gin.Recovery(),
	)

	for _, r := range routers {
		r.Register(router)
	}

	return router.Handler()
}

func CustomLogger(log logger.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		log.Info("WWW",
			logger.String("method", c.Request.Method),
			logger.String("path", c.Request.URL.Path),
			logger.Int("status", status),
			logger.String("latency", latency.String()),
		)
	}
}
