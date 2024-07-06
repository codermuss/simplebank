package gapi

import (
	"fmt"

	db "github.com/mustafayilmazdev/simplebank/db/sqlc"
	simplebank "github.com/mustafayilmazdev/simplebank/pb"
	"github.com/mustafayilmazdev/simplebank/token"
	"github.com/mustafayilmazdev/simplebank/util"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	simplebank.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new gRPC server.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
