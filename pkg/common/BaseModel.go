package common

import (
	"github.com/yockii/qscore/pkg/database"
	"github.com/yockii/qscore/pkg/util"
)

type BaseModel struct {
	ID uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
}

func (*BaseModel) TableComment() string {
	return "空表"
}
func (m *BaseModel) AddRequired() string {
	return ""
}
func (m *BaseModel) InstanceRequired() string {
	if m.ID == 0 {
		return "id"
	}
	return ""
}
func (m *BaseModel) CheckDuplicatedModel() database.Model {
	return nil
}
func (m *BaseModel) UpdateConditionModel() database.Model {
	b := new(BaseModel)
	b.ID = m.ID
	return b
}
func (m *BaseModel) UpdateModel() database.Model {
	panic("implement me")
}
func (m *BaseModel) InitDefaultFields() {
	m.ID = util.SnowflakeId()
}
func (m *BaseModel) UpdateRequired() string {
	if m.ID == 0 {
		return "id"
	}
	return ""
}
func (m *BaseModel) DeleteRequired() string {
	if m.ID == 0 {
		return "id"
	}
	return ""
}
func (m *BaseModel) FuzzyQueryMap() map[string]string {
	return nil
}
func (m *BaseModel) ExactMatchModel() database.Model {
	return m
}
func (m *BaseModel) ListOmits() string {
	return ""
}
