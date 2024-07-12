package gapi

import (
	"testing"
	"time"

	db "github.com/mustafayilmazdev/simplebank/db/sqlc"
	"github.com/mustafayilmazdev/simplebank/util"
	"github.com/mustafayilmazdev/simplebank/worker"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymetricKey:    util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)
	return server
}
