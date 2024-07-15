package gapi

import (
	"context"

	db "github.com/mustafayilmazdev/simplebank/db/sqlc"
	simplebank "github.com/mustafayilmazdev/simplebank/pb"
	"github.com/mustafayilmazdev/simplebank/util"
	"github.com/mustafayilmazdev/simplebank/val"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) AddAccount(ctx context.Context, req *simplebank.AddAccountRequest) (*simplebank.BaseResponse, error) {
	authPayload, err := server.authorizeUser(ctx, []string{util.BankerRole, util.DepositorRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if violations := validateAddAccountRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.ForeignKeyViolation || errCode == db.UniqueViolation {

			return nil, status.Errorf(codes.PermissionDenied, "%s", err)
		}
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	rsp := &simplebank.AddAccountResponse{
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

func validateAddAccountRequest(req *simplebank.AddAccountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateCurrency(req.GetCurrency()); err != nil {
		violations = append(violations, fieldViolation("page_id", err))
	}
	return violations
}
