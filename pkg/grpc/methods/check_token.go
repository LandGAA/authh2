package methods

import (
	"context"
	"github.com/LandGAA/authh2/internal/usecase"
	pd "github.com/LandGAA/authh2/pkg/grpc/generate"
	"github.com/LandGAA/authh2/pkg/jwt"
	"github.com/LandGAA/authh2/pkg/logger"
	"strconv"
)

type UserServiceServer struct {
	pd.UnimplementedUserServiceServer
	uu usecase.UseCase
}

func (s *UserServiceServer) CheckToken(ctx context.Context, req *pd.TokenRequest) (*pd.UserResponse, error) {
	logger.Logger.Info("Получен запрос на обновление токена от Forum")

	claims, err := jwt.ValidateToken(req.Access)
	if err != nil {
		return nil, err
	}

	return &pd.UserResponse{
		Id:    strconv.Itoa(claims.ID),
		Role:  claims.Role,
		Email: claims.Email,
	}, nil
}

func (s *UserServiceServer) GetUserByID(ctx context.Context, req *pd.IDRequest) (*pd.UserResponse, error) {
	user, err := s.uu.GetUserByID(int(req.Id))
	if err != nil {
		return nil, err
	}
	return &pd.UserResponse{
		Id:    strconv.Itoa(user.ID),
		Role:  user.Role,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}
