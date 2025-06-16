package web

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Christian-007/fit-forge/internal/app/points/dto"
	"github.com/Christian-007/fit-forge/internal/app/points/repositories"
	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/Christian-007/fit-forge/internal/pkg/apphttp"
	"github.com/Christian-007/fit-forge/internal/pkg/applog"
	"github.com/Christian-007/fit-forge/internal/pkg/requestctx"
	"github.com/Christian-007/fit-forge/internal/utils"
)

type PointTransactionsHandler struct {
	PointTransactionsOptions
}

type PointTransactionsOptions struct {
	Logger      applog.Logger
	PointTransactionsRepository repositories.PointTransactionsRepositoryPg
}

func NewPointTransactionsHandler(options PointTransactionsOptions) PointTransactionsHandler {
	return PointTransactionsHandler{options}
}

func (p PointTransactionsHandler) GetAllWithPagination(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := requestctx.UserId(ctx)
	if !ok {
		p.Logger.Error(apperrors.ErrTypeAssertion.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	query := r.URL.Query()
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(query.Get("offset"))
	if err != nil {
		offset = 0
	}

	pointTransactions, total, err := p.PointTransactionsRepository.GetAllWithPagination(ctx, userId, limit, offset)
	if err != nil {
		p.Logger.Error("failed to get all point transactions with pagination", slog.String("error", err.Error()))
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	utils.SendResponse(w, http.StatusOK, dto.PointTransactionsDto{
		Data: pointTransactions,
		Meta: dto.PaginationMeta{
			Total: total,
			Limit: limit,
			Offset: offset,
		},
	})
}

