package api

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/link33/sidecar/cmd/sidecar/client"
	"github.com/link33/sidecar/internal/repo"
	"github.com/link33/sidecar/model/pb"
	"github.com/sirupsen/logrus"
)

type Server struct {
	router *gin.Engine
	config *repo.Config
	logger logrus.FieldLogger
	ctx    context.Context
	cancel context.CancelFunc
}

type response struct {
	Data []byte `json:"data"`
}

func NewServer(config *repo.Config, logger logrus.FieldLogger) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	return &Server{
		router: router,
		config: config,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (g *Server) Start() error {
	g.router.Use(gin.Recovery())
	v1 := g.router.Group("/v1")
	{
		v1.POST(client.RegisterAppchainUrl, g.registerAppchain)
		v1.POST(client.UpdateAppchainUrl, g.updateAppchain)
		v1.GET(client.GetAppchainUrl, g.getAppchain)
	}

	go func() {
		go func() {
			err := g.router.Run(fmt.Sprintf(":%d", g.config.Port.Http))
			if err != nil {
				panic(err)
			}
		}()
		<-g.ctx.Done()
	}()
	return nil
}

func (g *Server) Stop() error {
	g.cancel()
	g.logger.Infoln("gin service stop")
	return nil
}

func (g *Server) updateAppchain(c *gin.Context) {
	g.sendAppchain(c, pb.Message_APPCHAIN_UPDATE)
}

func (g *Server) registerAppchain(c *gin.Context) {
	g.sendAppchain(c, pb.Message_APPCHAIN_REGISTER)
}

func (g *Server) sendAppchain(c *gin.Context, appchainType pb.Message_Type) {
	panic("implement me")
}

func (g *Server) getAppchain(c *gin.Context) {
	panic("implement me")
}
