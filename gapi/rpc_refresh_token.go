package gapi

import (
	"context"
	"errors"
	"time"

	db "github.com/mustafayilmazdev/simplebank/db/sqlc"
	simplebank "github.com/mustafayilmazdev/simplebank/pb"
	"github.com/mustafayilmazdev/simplebank/token"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) RefreshToken(ctx context.Context, req *simplebank.RefreshTokenRequest) (*simplebank.RefreshTokenResponse, error) {
	log.Info().Msgf("refreshToken: %s", req.GetRefreshToken())
	refreshPayload, violations := validateRefreshTokenRequest(req, server)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "%s", err)
	}
	if session.IsBlocked {
		return nil, status.Errorf(codes.Unauthenticated, "blocked session: %s", err)
	}
	if session.Username != refreshPayload.Username {
		return nil, status.Errorf(codes.Unauthenticated, "incorrect session user: %s", err)
	}

	if session.RefreshToken != req.RefreshToken {
		return nil, status.Errorf(codes.Unauthenticated, "mismatched session token: %s", err)
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, status.Errorf(codes.Unauthenticated, "expired session: %s", err)
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.Username, refreshPayload.Role, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	rsp := &simplebank.RefreshTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: timestamppb.New(accessPayload.ExpiredAt),
	}

	return rsp, nil
}

func validateRefreshTokenRequest(req *simplebank.RefreshTokenRequest, server *Server) (refreshPayload *token.Payload, violations []*errdetails.BadRequest_FieldViolation) {

	refreshPayload, err := server.tokenMaker.VerifyToken(req.GetRefreshToken())

	if err != nil {
		violations = append(violations, fieldViolation("refresh_token", err))
	}

	return refreshPayload, violations
}
