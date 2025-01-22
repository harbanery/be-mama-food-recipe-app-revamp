package helper

import "github.com/gofiber/fiber/v2"

type ParamsRequest struct {
	Keywords *string `json:"keywords"`
	Page     int     `json:"page"`
	Limit    int     `json:"limit"`
	Sort     string  `json:"sort"`
	OrderBy  string  `json:"order_by"`
}

type ParamsResponse struct {
	RecordsFiltered int  `json:"records_filtered"`
	RecordsTotal    int  `json:"records_total"`
	PageTotal       int  `json:"page_total"`
	CurrentPage     int  `json:"current_page"`
	HasPrevious     bool `json:"has_previous"`
	HasNext         bool `json:"has_next"`
	Start           int  `json:"start"`
	Records         any  `json:"records"`
}

func NewParamsRequest(ctx *fiber.Ctx) *ParamsRequest {
	var page = 1
	var limit = 10
	var sort = "created_at"
	var orderBy = "desc"
	var keywords string
	params := ctx.Queries()

	if params["keywords"] != "" {
		keywords = ctx.Query("keywords")
	}

	if params["page"] != "" && params["page"] != "0" {
		page = ctx.QueryInt("page")
	}

	if params["limit"] != "" && params["limit"] != "0" {
		limit = ctx.QueryInt("limit")
	}

	if params["sort"] != "" {
		sort = ctx.Query("sort")
	}

	if params["order_by"] != "" {
		orderBy = ctx.Query("order_by")
	}

	return &ParamsRequest{
		Keywords: &keywords,
		Page:     page,
		Limit:    limit,
		Sort:     sort,
		OrderBy:  orderBy,
	}
}

func NewParamsResponse(data any, counts *RecordCount, params *ParamsRequest) *ParamsResponse {
	var pageTotal, offset int

	offset = (params.Page - 1) * params.Limit
	pageTotal = int((counts.TotalData + int64(params.Limit) - 1) / int64(params.Limit))

	if pageTotal <= 0 {
		offset = 0
		pageTotal = 1
	}

	return &ParamsResponse{
		RecordsFiltered: int(counts.FilteredData),
		RecordsTotal:    int(counts.TotalData),
		PageTotal:       pageTotal,
		CurrentPage:     params.Page,
		HasPrevious:     params.Page > 1,
		HasNext:         params.Page < pageTotal,
		Start:           offset + 1,
		Records:         data,
	}
}
