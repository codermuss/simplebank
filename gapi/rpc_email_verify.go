package gapi

import (
	"context"

	db "github.com/mustafayilmazdev/simplebank/db/sqlc"
	simplebank "github.com/mustafayilmazdev/simplebank/pb"
	"github.com/mustafayilmazdev/simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *simplebank.VerifyEmailRequest) (*simplebank.VerifyEmailResponse, error) {
	if violations := validateVerifyEmailRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	txResult, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailID:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify email")
	}

	rsp := &simplebank.VerifyEmailResponse{
		IsVerified: txResult.User.IsEmailVerified,
	}
	return rsp, nil
}

func validateVerifyEmailRequest(req *simplebank.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmailID(req.GetEmailId()); err != nil {
		violations = append(violations, fieldViolation("email_id", err))
	}

	if err := val.ValidateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, fieldViolation("secret_code", err))
	}

	return violations
}
