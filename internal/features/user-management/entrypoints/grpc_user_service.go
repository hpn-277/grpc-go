package entrypoints

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nguyenphuoc/super-salary-sacrifice/internal/features/user-management/application"
	pb "github.com/nguyenphuoc/super-salary-sacrifice/proto"
)

// UserServiceServer implements the gRPC UserService interface
// This is an ENTRYPOINT in hexagonal architecture (inbound adapter)
type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	userService *application.UserService
}

// NewUserServiceServer creates a new gRPC user service server
func NewUserServiceServer(userService *application.UserService) *UserServiceServer {
	return &UserServiceServer{
		userService: userService,
	}
}

// CreateUser handles the CreateUser gRPC request
func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Map gRPC request to application command
	cmd := application.CreateUserCommand{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	// Execute use case
	user, err := s.userService.CreateUser(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// Map domain user to gRPC response
	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:        user.ID().String(),
			Email:     user.Email().String(),
			FirstName: user.FirstName(),
			LastName:  user.LastName(),
			CreatedAt: timestamppb.New(user.CreatedAt()),
			UpdatedAt: timestamppb.New(user.UpdatedAt()),
		},
	}, nil
}

// GetUser handles the GetUser gRPC request
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// Map gRPC request to application query
	query := application.GetUserQuery{
		UserID: req.UserId,
	}

	// Execute use case
	user, err := s.userService.GetUser(ctx, query)
	if err != nil {
		return nil, err
	}

	// Map domain user to gRPC response
	return &pb.GetUserResponse{
		User: &pb.User{
			Id:        user.ID().String(),
			Email:     user.Email().String(),
			FirstName: user.FirstName(),
			LastName:  user.LastName(),
			CreatedAt: timestamppb.New(user.CreatedAt()),
			UpdatedAt: timestamppb.New(user.UpdatedAt()),
		},
	}, nil
}

// ListUsers handles the ListUsers gRPC request
func (s *UserServiceServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// Map gRPC request to application query
	query := application.ListUsersQuery{
		Offset: int(req.Offset),
		Limit:  int(req.Limit),
	}

	// Execute use case
	users, err := s.userService.ListUsers(ctx, query)
	if err != nil {
		return nil, err
	}

	// Map domain users to gRPC response
	pbUsers := make([]*pb.User, len(users))
	for i, user := range users {
		pbUsers[i] = &pb.User{
			Id:        user.ID().String(),
			Email:     user.Email().String(),
			FirstName: user.FirstName(),
			LastName:  user.LastName(),
			CreatedAt: timestamppb.New(user.CreatedAt()),
			UpdatedAt: timestamppb.New(user.UpdatedAt()),
		}
	}

	return &pb.ListUsersResponse{
		Users: pbUsers,
		Total: int32(len(pbUsers)),
	}, nil
}
