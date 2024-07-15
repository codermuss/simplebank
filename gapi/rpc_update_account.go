package gapi

import (
	"context"
	"errors"

	db "github.com/mustafayilmazdev/simplebank/db/sqlc"
	simplebank "github.com/mustafayilmazdev/simplebank/pb"
	"github.com/mustafayilmazdev/simplebank/util"
	"github.com/mustafayilmazdev/simplebank/val"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) UpdateAccount(ctx context.Context, req *simplebank.UpdateAccountRequest) (*simplebank.BaseResponse, error) {

	authPayload, err := server.authorizeUser(ctx, []string{util.BankerRole, util.DepositorRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if violations := validateUpdateAccount(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	acc, err := server.store.GetAccount(ctx, req.AccountId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get account: %s", err)
	}

	if authPayload.Username != acc.Owner {
		return nil, status.Error(codes.PermissionDenied, "cannot update other's accounts")
	}

	arg := db.UpdateAccountParams{
		ID:      req.AccountId,
		Balance: req.Balance + acc.Balance,
	}

	account, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "account not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to Update account: %s", err)
	}
	rsp := &simplebank.UpdateAccountResponse{
		Account: &simplebank.Accounts{
			Id:        account.ID,
			Balance:   account.Balance,
			Currency:  account.Currency,
			CreatedAt: timestamppb.New(account.CreatedAt),
		},
	}
	response, err := anypb.New(rsp)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", &simplebank.BaseResponse{
			Code: int32(codes.Internal),
			Data: response,
		})
	}

	return &simplebank.BaseResponse{
		Code: 200,
		Data: response,
	}, nil
}

func validateUpdateAccount(req *simplebank.UpdateAccountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	log.Info().Msgf("params: %s", req)
	if err := val.ValidateEmailID(req.GetAccountId()); err != nil {
		violations = append(violations, fieldViolation("account_id", err))
	}

	if err := val.ValidateBalance(req.GetBalance()); err != nil {
		violations = append(violations, fieldViolation("balance", err))
	}

	return violations
}
