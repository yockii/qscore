package server

import "github.com/yockii/qscore/pkg/database"

type CommonResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type TimeCondition struct {
	Start database.DateTime `json:"start,omitempty" query:"start"`
	End   database.DateTime `json:"end,omitempty" query:"end"`
}

type Paginate[T any] struct {
	Total  int64 `json:"total"`
	Offset int   `json:"offset"`
	Limit  int   `json:"limit"`
	Items  []T   `json:"items"`
}
