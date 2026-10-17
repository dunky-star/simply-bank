package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/dunky-star/simply-bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type transferRequest struct {
	FromAccountID string          `json:"from_account_id" binding:"required"`
	ToAccountID   string          `json:"to_account_id" binding:"required"`
	Amount        decimal.Decimal `json:"amount" binding:"required"`
	Currency      string          `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromID, err := uuid.Parse(req.FromAccountID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	toID, err := uuid.Parse(req.ToAccountID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if _, valid := server.validAccount(ctx, fromID, req.Currency); !valid {
		return
	}
	if _, valid := server.validAccount(ctx, toID, req.Currency); !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: fromID,
		ToAccountID:   toID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID uuid.UUID, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("account not found")))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%s] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
