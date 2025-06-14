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
func NewServeMux(routers []Router, log logger.ILogger) http.Handler {
	router := gin.New()

	router.Use(
		logger.GinMiddleware(log),
		gin.Recovery(),
	)

	for _, r := range routers {
		r.Register(router)
	}

	return router.Handler()
}