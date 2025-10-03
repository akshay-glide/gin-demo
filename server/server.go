package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"gin-demo/config"
	"gin-demo/handlers"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ServerHandlerMap struct {
	APIPath string
	Handler handlers.APIHandler
}

type Server struct {
	Config     *config.ConfigAPIServer
	GinApp     *gin.Engine
	Handlers   []*ServerHandlerMap
	HTTPServer *http.Server
}

func NewServerHandlerMap(apipath string, handler handlers.APIHandler) *ServerHandlerMap {
	return &ServerHandlerMap{APIPath: apipath, Handler: handler}
}

func (s *Server) Setup() *Server {
	if s.GinApp == nil {
		log.Fatal().Msg("Server setup incorrectly!")
	}

	for _, each := range s.Handlers {
		group := s.GinApp.Group(each.APIPath)
		each.Handler.RegisterRoutes(group)
	}

	return s
}

func GetServer(config *config.ConfigAPIServer, app *gin.Engine, handlers []*ServerHandlerMap) *Server {
	serv := &Server{
		Config:   config,
		GinApp:   app,
		Handlers: handlers,
	}
	serv.Setup()
	return serv
}

func (s *Server) StartServer() <-chan os.Signal {
	port := ":" + strconv.Itoa(*s.Config.Port)

	s.HTTPServer = &http.Server{
		Addr:    port,
		Handler: s.GinApp,
	}

	// Channel to listen for system signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info().Msgf("Starting server on %s", port)
		if err := s.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	return quit
}

func (s *Server) ShutdownGracefully() {
	log.Info().Msg("Shutting down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.HTTPServer.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	} else {
		log.Info().Msg("Server shutdown completed successfully")
	}
}

func GetGinApplication() *gin.Engine {
	app := gin.New()
	app.Use(gin.Recovery()) // base middleware
	return app
}
