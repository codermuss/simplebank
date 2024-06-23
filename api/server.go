package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/mustafayilmazdev/simplebank/db/sqlc"
)

// * Note [codermuss]: Server serves HTTP requests for our banking service.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// * Note [codermuss]: NewServer creates a new HTTP server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	router.POST("/transfers", server.createTransfer)
	server.router = router
	return server
}

// * Note [codermuss]: Starts the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
