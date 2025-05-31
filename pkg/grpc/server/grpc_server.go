package server

import (
	pd "github.com/LandGAA/authh2/pkg/grpc/generate"
	"github.com/LandGAA/authh2/pkg/grpc/methods"
	"github.com/LandGAA/authh2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func Run() {
	ls, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Logger.Fatal("Ошибка запуска gRPC сервера!",
			zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	pd.RegisterUserServiceServer(grpcServer, &methods.UserServiceServer{})

	logger.Logger.Info("gRPC сервер запущен!")
	if err := grpcServer.Serve(ls); err != nil {
		logger.Logger.Fatal("Ошибка запуска gRPC сервера!",
			zap.Error(err))
	}
}
