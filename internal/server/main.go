package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/klog"

	"github.com/joyrex2001/kubedock/internal/config"
	"github.com/joyrex2001/kubedock/internal/kubernetes"
	"github.com/joyrex2001/kubedock/internal/model"
	"github.com/joyrex2001/kubedock/internal/server/httputil"
	"github.com/joyrex2001/kubedock/internal/server/routes"
)

// Server is the API server.
type Server struct {
	db *model.Database
}

// New will instantiate a Server object.
func New(db *model.Database) *Server {
	return &Server{db: db}
}

// Run will initialize the http api server and configure all available
// routers.
func (s *Server) Run() error {
	if !klog.V(2) {
		gin.SetMode(gin.ReleaseMode)
	}

	router := s.getGinEngine()
	err := s.setUpRoutes(router)
	if err != nil {
		return err
	}

	socket := viper.GetString("server.socket")
	if socket == "" {
		port := viper.GetString("server.listen-addr")
		if viper.GetBool("server.enable-tls") {
			cert := viper.GetString("server.cert-file")
			key := viper.GetString("server.key-file")
			router.RunTLS(port, cert, key)
		} else {
			router.Run(port)
		}
	} else {
		router.RunUnix(socket)
	}

	return nil
}

// getGinEngine will return a gin.Engine router and configure the
// appropriate middleware.
func (s *Server) getGinEngine() *gin.Engine {
	router := gin.New()
	router.Use(httputil.VersionAliasMiddleware(router))
	router.Use(gin.Logger())
	router.Use(httputil.RequestLoggerMiddleware())
	router.Use(httputil.ResponseLoggerMiddleware())
	router.Use(gin.Recovery())
	return router
}

// setUpRoutes will configure the routes for the server.
func (s *Server) setUpRoutes(router *gin.Engine) error {
	cfg, err := config.GetKubernetes()
	if err != nil {
		return err
	}

	cli, err := clientset.NewForConfig(cfg)
	if err != nil {
		return err
	}

	kube := kubernetes.New(cfg, cli, viper.GetString("kubernetes.namespace"))
	routes.New(router, s.db, kube)

	return nil
}
