package http

import (
	"context"
	"log"
	v1_0_0 "main/src/http/controllers/v1.0.0"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	srv    *http.Server

	// Controllers
	// v1.0.0
	apiController  *v1_0_0.ApiController
	userController *v1_0_0.UserController

	wg             *sync.WaitGroup
	wsShutdownChan chan struct{}
	isReady        bool
}

func (s *Server) Init(addr string, wg *sync.WaitGroup, isDebug bool, wsShutdownChan chan struct{}) {
	s.wsShutdownChan = wsShutdownChan

	if !isDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	// get new router in init function (not in constructor) because git.SetMode global function
	// and gin.SetMode don't affect to router structure
	s.router = gin.New()
	s.wg = wg
	s.registerValidators()
	s.registerGlobalMiddlewares()
	s.initRoutes()

	s.srv = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}
}

func (s *Server) Run(ctx context.Context) error {
	httpShutdownCh := make(chan struct{})

	go func() {
		<-ctx.Done()

		log.Println("webServer.shutdown: init")

		graceCtx, graceCancel := context.WithTimeout(ctx, 1*time.Second)
		defer graceCancel()

		if err := s.srv.Shutdown(graceCtx); err != nil {
			log.Println(err)
		}

		httpShutdownCh <- struct{}{}
	}()

	go func() {
		s.wg.Add(1)
		defer s.wg.Done()

		s.isReady = true
		err := s.srv.ListenAndServe()
		if err != http.ErrServerClosed {
			panic(err)
		}

		<-httpShutdownCh
		s.isReady = false

		log.Println("webServer.shutdown: complete")
		close(s.wsShutdownChan)
	}()

	return nil
}

func (s *Server) IsReady() bool {
	return s.isReady
}

func (s *Server) registerValidators() {
}

func (s *Server) registerGlobalMiddlewares() {

}

func (s *Server) initRoutes() {
	v100 := s.router.Group("/v1.0.0/")
	{
		v100.GET("/", s.apiController.Version)
		v100.GET("/users/:id", s.userController.GetUser)
	}
}

func NewServer(
	apiController *v1_0_0.ApiController,
	userController *v1_0_0.UserController,
) *Server {
	return &Server{
		apiController:  apiController,
		userController: userController,
	}
}
