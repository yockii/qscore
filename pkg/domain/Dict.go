package domain

const (
	DictIdPrefix = "dict"
)

type Dict struct {
	Id         string   `json:"id,omitempty" xorm:"pk varchar(50)"`
	DictKey    string   `json:"dictKey,omitempty" xorm:"index varchar(50) comment('字典键')"`
	DictValue  string   `json:"dictValue,omitempty" xorm:"comment('字典值')"`
	DictExt    string   `json:"dictExt,omitempty" xorm:"comment('字典扩展值')"`
	ParentId   string   `json:"parentId,omitempty" xorm:"comment('父ID，若无则为字典分类')"`
	CreateTime DateTime `json:"createTime,omitempty" xorm:"created"`
}

func init() {
	SyncDomains = append(SyncDomains, Dict{})
}

type DictRequest struct {
	*Dict
	CreateTimeRange *TimeCondition `json:"createTimeRange,omitempty"`
}
