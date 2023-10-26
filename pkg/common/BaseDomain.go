package common

import (
	"github.com/yockii/qscore/pkg/database"
	"github.com/yockii/qscore/pkg/server"
)

type Domain[T database.Model] interface {
	GetModel() T
	GetOrderBy() string
	GetTimeConditionList() map[string]*server.TimeCondition
}

type BaseDomain[T database.Model] struct {
	OrderBy string `json:"orderBy,omitempty"`
}

func (r *BaseDomain[T]) GetOrderBy() string {
	return r.OrderBy
}

func (r *BaseDomain[T]) GetTimeConditionList() map[string]*server.TimeCondition {
	return nil
}
