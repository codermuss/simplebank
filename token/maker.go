package token

import "time"

// * Note [codermuss]: Maker is an interface for managing tokens

type Maker interface {
	// * Note [codermuss]: CreateToken creates a new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)

	// * Note [codermuss]: VerifyToken checks if the token is valid or not

	VerifyToken(token string) (*Payload, error)
}
