package common

import "github.com/yockii/qscore/pkg/server"

type BaseDomain[T Model] interface {
	GetModel() T
	GetOrderBy() string
	GetTimeConditionList() map[string]*server.TimeCondition
}
