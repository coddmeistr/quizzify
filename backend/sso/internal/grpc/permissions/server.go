package permissionsgrpc

import (
	"context"
	"errors"

	ssov1 "github.com/coddmeistr/quizzify-online-tests/backend/protos/proto/sso"
	"github.com/coddmeistr/quizzify-online-tests/backend/sso/internal/services/permissions"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Permissions interface {
	AddPermission(ctx context.Context, userID int64, permID int64) error
	RemovePermission(ctx context.Context, userID int64, permID int64) error
}

const (
	emptyValue = 0
)

type serverAPI struct {
	ssov1.UnimplementedPermissionServer
	perm Permissions
}

func Register(gRPC *grpc.Server, perm Permissions) {
	ssov1.RegisterPermissionServer(gRPC, &serverAPI{perm: perm})
}

func (s *serverAPI) AddPermission(ctx context.Context, req *ssov1.AddPermissionRequest) (*ssov1.AddPermissionResponse, error) {
	if err := validateAddPermission(req); err != nil {
		return nil, err
	}

	if err := s.perm.AddPermission(ctx, req.GetUserId(), req.GetPermissionId()); err != nil {
		if errors.Is(err, permissions.ErrPermissionExist) {
			return nil, status.Error(codes.AlreadyExists, "already exist")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.AddPermissionResponse{
		Granted: true,
	}, nil
}

func (s *serverAPI) RemovePermission(ctx context.Context, req *ssov1.RemovePermissionRequest) (*ssov1.RemovePermissionResponse, error) {
	if err := validateRemovePermission(req); err != nil {
		return nil, err
	}

	if err := s.perm.RemovePermission(ctx, req.GetUserId(), req.GetPermissionId()); err != nil {
		if errors.Is(err, permissions.ErrPermissionNotExist) {
			return nil, status.Error(codes.NotFound, "nothing to remove")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RemovePermissionResponse{
		Removed: true,
	}, nil
}

func validateAddPermission(req *ssov1.AddPermissionRequest) error {
	if req.GetPermissionId() == emptyValue {
		return status.Error(codes.InvalidArgument, "permission_id is zero")
	}

	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user_id is zero")
	}

	return nil
}

func validateRemovePermission(req *ssov1.RemovePermissionRequest) error {
	if req.GetPermissionId() == emptyValue {
		return status.Error(codes.InvalidArgument, "permission_id is zero")
	}

	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user_id is zero")
	}

	return nil
}
