package client

import (
	"context"
	"fmt"
	pd "github.com/LandGAA/authh2/pkg/grpc/generate"
	"github.com/LandGAA/authh2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	pd.UserServiceClient
}

func NewClient() (*Client, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		logger.Logger.Fatal("Ошибка подключения клиента к gRPC серверу",
			zap.Error(err))
		return nil, fmt.Errorf("Ошибка подключения клиента к gRPC серверу %v", err)
	}
	defer conn.Close()
	return &Client{pd.NewUserServiceClient(conn)}, nil
}

func (c *Client) CheckToken(token string) (*pd.UserResponse, error) {
	ctx := metadata.AppendToOutgoingContext(
		context.Background(),
		"authorization", token,
	)

	res, err := c.UserServiceClient.CheckToken(ctx, &pd.TokenRequest{Access: token})
	if err != nil {
		logger.Logger.Error("Ошибка при вызове CheckToken",
			zap.Error(err))
		return &pd.UserResponse{}, fmt.Errorf("Ошибка при вызове CheckToken: %v", err)
	}

	return res, nil
}

func (c *Client) GetUserByID(id int) (*pd.UserResponse, error) {
	user, err := c.UserServiceClient.GetUserByID(context.TODO(), &pd.IDRequest{Id: int64(id)})
	if err != nil {
		logger.Logger.Error("Ошибка при вызове GetUserByID",
			zap.Error(err))
		return &pd.UserResponse{}, fmt.Errorf("Ошибка при вызове GetUserByID: %v", err)
	}

	return &pd.UserResponse{
		Id:    user.Id,
		Role:  user.Role,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}
