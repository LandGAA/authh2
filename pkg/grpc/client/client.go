package client

import (
	"context"
	"fmt"
	pd "github.com/LandGAA/authh2/pkg/grpc/generate"
	"github.com/LandGAA/authh2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Client struct {
	Conn pd.UserServiceClient
}

func NewClient() (Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		logger.Logger.Error("Ошибка подключения клиента к gRPC серверу",
			zap.Error(err))
		return Client{}, fmt.Errorf("ошибка подключения к gRPC серверу: %v", err)
	}

	return Client{
		Conn: pd.NewUserServiceClient(conn),
	}, nil
}
