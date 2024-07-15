package gapi

import (
	"context"
	"errors"
	"fmt"

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

func (server *Server) Transfer(ctx context.Context, req *simplebank.TransferRequest) (*simplebank.BaseResponse, error) {

	_, err := server.authorizeUser(ctx, []string{util.BankerRole, util.DepositorRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if violations := validateTransfer(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	fromAccount, err := server.store.GetAccount(ctx, req.FromAccountId)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "%s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get from account: %s", err)
	}

	toAccount, err := server.store.GetAccount(ctx, req.ToAccountId)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "%s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get to account: %s", err)
	}

	if fromAccount.Currency != toAccount.Currency {
		err := fmt.Errorf("account [%d] - [%d] currency mismatch: %s vs %s", fromAccount.ID, toAccount.ID, fromAccount.Currency, toAccount.Currency)
		return nil, status.Errorf(codes.Canceled, "failed to get account: %s", err)
	}
	if fromAccount.Balance-req.Amount < 0 {
		err := fmt.Errorf("insuffiecent balance to transfer: Balance: [%f] - Trasfer Amount: [%f]", fromAccount.Balance, req.Amount)
		return nil, status.Errorf(codes.Canceled, "failed to get account: %s", err)
	}

	arg := db.TransferTxParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   req.ToAccountId,
		Amount:        req.Amount,
	}

	transfer, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "account not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to Update account: %s", err)
	}
	rsp := &simplebank.TransferResponse{
		Id:            transfer.Transfer.ID,
		FromAccountId: transfer.Transfer.FromAccountID,
		ToAccountId:   transfer.Transfer.ToAccountID,
		Amount:        transfer.Transfer.Amount,
		CreatedAt:     timestamppb.New(transfer.Transfer.CreatedAt),
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

func validateTransfer(req *simplebank.TransferRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateBalance(req.GetAmount()); err != nil {
		violations = append(violations, fieldViolation("amount", err))
	}

	return violations
}
