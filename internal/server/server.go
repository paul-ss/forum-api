package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/paul-ss/forum-api/configs/go"

	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	ctx *Context
	server  *http.Server
}




func New() *Server {
	ctx := NewContext()

	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	r.Static(config.Conf.Web.Static.UrlImg, config.Conf.Web.Static.DirImg)


	return &Server{
		ctx: ctx,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", config.Conf.Web.Server.Address, config.Conf.Web.Server.Port),
			Handler: r,
		},
	}
}

func (s *Server) Run() {
	defer s.ctx.Cleanup()

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	config.Lg("server", "Run").Info("Server listening on " + s.server.Addr)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	config.Lg("server", "Run").Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		config.Lg("server", "Run").Fatal("Server forced to shutdown:", err)
	}
}


