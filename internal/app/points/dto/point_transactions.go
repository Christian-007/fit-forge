package dto

import "github.com/Christian-007/fit-forge/internal/app/points/domains"

type PointTransactionsDto struct {
	Data []domains.PointTransaction `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}