package db

import (
	"context"
	"database/sql"
	"fmt"
	"slices"

	"github.com/google/uuid"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID uuid.UUID `json:"from_account_id"`
	ToAccountID   uuid.UUID `json:"to_account_id"`
	Amount        int64     `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, and update accounts' balance
// within a single database transaction.
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Order accounts by ID to avoid deadlocks under concurrent opposite-direction transfers.
		fromID, toID := arg.FromAccountID, arg.ToAccountID
		fromAmt, toAmt := int64(-arg.Amount), arg.Amount
		if slices.Compare(fromID[:], toID[:]) > 0 {
			fromID, toID = toID, fromID
			fromAmt, toAmt = toAmt, fromAmt
		}

		if _, err = q.GetAccountForUpdate(ctx, fromID); err != nil {
			return err
		}
		if _, err = q.GetAccountForUpdate(ctx, toID); err != nil {
			return err
		}

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		first, second, err := addMoney(ctx, q, fromID, fromAmt, toID, toAmt)
		if err != nil {
			return err
		}
		if fromID == arg.FromAccountID {
			result.FromAccount, result.ToAccount = first, second
		} else {
			result.FromAccount, result.ToAccount = second, first
		}

		return nil
	})

	return result, err
}

// addMoney applies two atomic balance deltas in the given order and returns the updated accounts.
func addMoney(ctx context.Context, q *Queries, id1 uuid.UUID, amt1 int64, id2 uuid.UUID, amt2 int64) (Account, Account, error) {
	a1, err := q.AddAccountBalance(ctx, AddAccountBalanceParams{ID: id1, Amount: amt1})
	if err != nil {
		return Account{}, Account{}, err
	}
	a2, err := q.AddAccountBalance(ctx, AddAccountBalanceParams{ID: id2, Amount: amt2})
	if err != nil {
		return Account{}, Account{}, err
	}
	return a1, a2, nil
}
