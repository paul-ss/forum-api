package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/paul-ss/forum-api/configs/go"
	delivery2 "github.com/paul-ss/forum-api/internal/app/forum/delivery"
	delivery3 "github.com/paul-ss/forum-api/internal/app/post/delivery"
	delivery4 "github.com/paul-ss/forum-api/internal/app/thread/delivery"
	"github.com/paul-ss/forum-api/internal/app/user/delivery"

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

	delivery.CreateUserDelivery(ctx.db.DbPool, r)
	delivery2.CreateForumDelivery(ctx.db.DbPool, r)
	delivery3.CreatePostDelivery(ctx.db.DbPool, r)
	delivery4.CreateThreadDelivery(ctx.db.DbPool, r)

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


