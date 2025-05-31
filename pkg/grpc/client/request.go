package client

import (
	"context"
	"fmt"
	pd "github.com/LandGAA/authh2/pkg/grpc/generate"
	"github.com/LandGAA/authh2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"time"
)

type Client struct {
	conn   *grpc.ClientConn
	client pd.UserServiceClient
}

func NewClient() (*Client, error) {
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
		return nil, fmt.Errorf("ошибка подключения к gRPC серверу: %v", err)
	}

	return &Client{
		conn:   conn,
		client: pd.NewUserServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) CheckToken(token string) (*pd.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", token)

	res, err := c.client.CheckToken(ctx, &pd.TokenRequest{Access: token})
	if err != nil {
		logger.Logger.Error("Ошибка при вызове CheckToken",
			zap.Error(err))
		return nil, fmt.Errorf("ошибка проверки токена: %v", err)
	}

	return res, nil
}

func (c *Client) GetUserByID(id int) (*pd.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := c.client.GetUserByID(ctx, &pd.IDRequest{Id: int64(id)})
	if err != nil {
		logger.Logger.Error("Ошибка при вызове GetUserByID",
			zap.Int("id", id),
			zap.Error(err))
		return nil, fmt.Errorf("ошибка получения пользователя: %v", err)
	}

	return user, nil
}
