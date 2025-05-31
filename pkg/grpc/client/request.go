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

func Check(token string) (*pd.UserResponse, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		logger.Logger.Fatal("Ошибка подключения клиента к gRPC серверу",
			zap.Error(err))
		return &pd.UserResponse{}, fmt.Errorf("Ошибка подключения клиента к gRPC серверу %v", err)
	}
	defer conn.Close()

	ctx := metadata.AppendToOutgoingContext(
		context.Background(),
		"authorization", token,
	)

	client := pd.NewUserServiceClient(conn)
	res, err := client.CheckToken(ctx, &pd.TokenRequest{Access: token})
	if err != nil {
		logger.Logger.Error("Ошибка при вызове CheckToken",
			zap.Error(err))
		return &pd.UserResponse{}, fmt.Errorf("Ошибка при вызове CheckToken: %v", err)
	}

	return res, nil
}

func GetUserByID(id int) (*pd.UserResponse, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		logger.Logger.Fatal("Ошибка подключения клиента к gRPC серверу",
			zap.Error(err))
		return &pd.UserResponse{}, fmt.Errorf("Ошибка подключения клиента к gRPC серверу %v", err)
	}
	defer conn.Close()

	client := pd.NewUserServiceClient(conn)
	user, err := client.GetUserByID(context.TODO(), &pd.IDRequest{Id: int64(id)})
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
