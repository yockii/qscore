package database

type Model interface {
	TableComment() string
	AddRequired() string
	InstanceRequired() string
	CheckDuplicatedModel() Model
	UpdateConditionModel() Model
	UpdateModel() Model
	InitDefaultFields()
	UpdateRequired() string
	FuzzyQueryMap() map[string]string
	ExactMatchModel() Model
	DeleteRequired() string
	ListOmits() string
}

var Models []Model
