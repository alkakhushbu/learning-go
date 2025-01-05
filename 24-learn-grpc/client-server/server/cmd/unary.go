package main

import (
	"context"
	"fmt"
	pb "server/gen/proto/v1"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NewUser struct {
	Id    int64
	Name  string   `json:"name" validate:"required"`        // Name of the user, a required field
	Email string   `json:"email" validate:"required,email"` // Email ID of the user, should be in valid email format
	Roles []string `json:"roles" validate:"required"`       // Roles assigned to the user, a required field
}

func (u userService) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	nu := req.GetUser()
	if nu == nil {
		return nil, status.Error(codes.InvalidArgument, "user is nil")
	}

	// Parsing the received request to our User struct
	newUser := NewUser{
		Id:    1,
		Name:  nu.Name,
		Email: nu.Email,
		Roles: nu.Roles,
	}
	validate := validator.New()
	err := validate.Struct(newUser)
	if err != nil {
		// If validation fails, return an error message with the error status.
		return nil, status.Error(codes.Internal, "please provide required fields in correct format")
	}

	fmt.Println(newUser)
	return &pb.SignupResponse{UserId: newUser.Id}, nil

}
