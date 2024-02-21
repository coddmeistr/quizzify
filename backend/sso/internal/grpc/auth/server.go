package authgrpc

import (
	"context"
	"errors"

	ssov1 "github.com/maxik12233/quizzify-online-tests/backend/protos/gen/go/sso"
	"github.com/maxik12233/quizzify-online-tests/backend/sso/internal/services/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, login string, email string, password string, appID int) (token string, err error)
	Register(ctx context.Context, login string, email string, password string) (userID uint64, err error)
	IsAdmin(ctx context.Context, userID uint64) (bool, error)
}

const (
	emptyValue = 0
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.Register(ctx, req.GetLogin(), req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) { // Not sure if we should take errors directly from Auth
			return nil, status.Error(codes.AlreadyExists, "User already exists")
		}
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: int64(userID),
	}, nil
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetLogin(), req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) { // Not sure if we should take errors directly from Auth
			return nil, status.Error(codes.InvalidArgument, "Invalid credentials or user not exists")
		}
		if errors.Is(err, auth.ErrAppNotFound) { // Not sure if we should take errors directly from Auth
			return nil, status.Error(codes.NotFound, "App not found")
		}
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, uint64(req.GetUserId()))
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) { // Not sure if we should take errors directly from Auth
			return nil, status.Error(codes.NotFound, "User with given ID not exist")
		}
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetLogin() == "" {
		return status.Error(codes.InvalidArgument, "login is empty")
	}

	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is empty")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is empty")
	}

	return nil
}

func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetLogin() == "" && req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "login and email are empty")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is empty")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is zero")
	}

	return nil
}

func validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user_id is zero")
	}

	return nil
}
