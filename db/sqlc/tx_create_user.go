package db

import "context"

type CreateUserTxParams struct {
	CreateUserParams
	// * Note [codermuss]: This function will be executed after the user is inserted,
	// * Note [codermuss]: inside the same transaction. And its output error will be used to decide
	// * Note [codermuss]: whether to commit or rollback the transaction.
	AfterCreate func(user User) error
}
type CreateUserTxResult struct {
	User User
}

// * Note [codermuss]: This method responsible with creating user.
// * Note [codermuss]: It uses execTx to handle DB Transaction error
func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}
		return arg.AfterCreate(result.User)
	})
	return result, err
}
