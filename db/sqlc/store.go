package mdb

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	txOption := pgx.TxOptions{}
	tx, err := s.db.BeginTx(ctx, txOption)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("rollback err : %v, transaction err: %v", rbErr, err)
		}
		return err
	}

	return tx.Commit(ctx)
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

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

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})

			if err != nil {
				return nil
			}

			result.ToAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})

			if err != nil {
				return nil
			}
		} else {
			result.ToAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})

			if err != nil {
				return nil
			}

			result.FromAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})

			if err != nil {
				return nil
			}

		}

		return nil
	})

	return result, err
}
