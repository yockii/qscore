package domain

const (
	UserIdPrefix     = "user"
	RoleIdPrefix     = "role"
	ResourceIdPrefix = "resource"
)

type User struct {
	Id         string   `json:"id,omitempty" xorm:"pk varchar(50)"`
	Username   string   `json:"username,omitempty" xorm:"index varchar(50) comment('用户名')"`
	Password   string   `json:"password,omitempty" xorm:"comment('密码')"`
	CreateTime DateTime `json:"createTime,omitempty" xorm:"created"`
}

type Role struct {
	Id         string   `json:"id,omitempty" xorm:"pk varchar(50)"`
	RoleName   string   `json:"roleName,omitempty" xorm:"varchar(50)"`
	RoleDesc   string   `json:"roleDesc,omitempty"`
	CreateTime DateTime `json:"createTime,omitempty" xorm:"created"`
}

type Resource struct {
	Id              string   `json:"id,omitempty" xorm:"pk varchar(50)"`
	ResourceName    string   `json:"resourceName,omitempty" xorm:"comment('资源名称')"`
	ResourceContent string   `json:"resourceContent,omitempty" xorm:"comment('资源内容，如url、数据分类等等')"`
	ResourceType    string   `json:"resourceType,omitempty" xorm:"comment('资源类型，定义：route、data')"`
	Action          string   `json:"action,omitempty" xorm:"comment('资源操作类型，如url有GET/POST/PUT/DELETE')"`
	CreateTime      DateTime `json:"createTime,omitempty" xorm:"created"`
}

func init() {
	SyncDomains = append(SyncDomains, User{}, Role{}, Resource{})
}

type UserRequest struct {
	User
	CreateTimeRange *TimeCondition `json:"createTimeRange,omitempty"`
}

type RoleRequest struct {
	Role
	CreateTimeRange *TimeCondition `json:"createTimeRange,omitempty"`
}

type ResourceRequest struct {
	Resource
	CreateTimeRange *TimeCondition `json:"createTimeRange,omitempty"`
}
