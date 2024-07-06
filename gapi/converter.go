package gapi

import (
	db "github.com/mustafayilmazdev/simplebank/db/sqlc"
	simplebank "github.com/mustafayilmazdev/simplebank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *simplebank.User {
	return &simplebank.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
