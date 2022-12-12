package external

import (
	"grpc/app/config"
	"grpc/app/proto/imagestorage"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	imagestorage.ImageStorageServer
	log  *logrus.Logger
	conf *config.Config
}

func NewGRPCServer(u imagestorage.ImageStorageServer, log *logrus.Logger, conf *config.Config) *GRPCServer {
	return &GRPCServer{
		ImageStorageServer: u,
		log:                log,
		conf:               conf,
	}
}

func (e *GRPCServer) Run() {
	server := grpc.NewServer(grpc.MaxConcurrentStreams(100))

	imagestorage.RegisterImageStorageServer(server, e.ImageStorageServer)

	reflection.Register(server)

	listener, err := net.Listen("tcp", e.conf.ServerIP())
	if err != nil {
		e.log.Warnln("не удалось запустить tcp слушатель:", err)
		return
	}

	e.log.Infof("grpc сервер запущен на :8080")

	if err := server.Serve(listener); err != nil {
		e.log.Fatalln("не удалось запустить сервер:", err)
	}

}
