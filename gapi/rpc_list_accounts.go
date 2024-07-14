package gapi

import (
	"context"
	"errors"

	db "github.com/mustafayilmazdev/simplebank/db/sqlc"
	simplebank "github.com/mustafayilmazdev/simplebank/pb"
	"github.com/mustafayilmazdev/simplebank/util"
	"github.com/mustafayilmazdev/simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) ListAccounts(ctx context.Context, req *simplebank.ListAccountsRequest) (*simplebank.ListAccountsResponse, error) {
	authPayload, err := server.authorizeUser(ctx, []string{util.BankerRole, util.DepositorRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if violations := validateListAccountsRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageId,
		Offset: req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "%s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed fetch accounts: %s", err)
	}

	items := []*simplebank.Accounts{}

	for _, account := range accounts {
		items = append(items, &simplebank.Accounts{
			Id:        account.ID,
			Balance:   account.Balance,
			Currency:  account.Currency,
			CreatedAt: timestamppb.New(account.CreatedAt),
		})
	}

	rsp := &simplebank.ListAccountsResponse{
		Accounts: items,
	}
	return rsp, nil
}

func validateListAccountsRequest(req *simplebank.ListAccountsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidatePageID(req.GetPageId()); err != nil {
		violations = append(violations, fieldViolation("page_id", err))
	}

	if err := val.ValidatePageSize(req.GetPageSize()); err != nil {
		violations = append(violations, fieldViolation("page_size", err))
	}

	return violations
}
