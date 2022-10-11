package authorization

import (
	"gitee.com/chunanyong/zorm"
)

// AuthorizePolicyTableName 表名常量,方便直接调用
const AuthorizePolicyTableName = "t_authorize_policy"

// AuthorizePolicy
type AuthorizePolicy struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	zorm.EntityStruct

	//Id []
	Id string `column:"id" json:"id"`

	//PolicyType 策略类型 1-p, 2-p2
	PolicyType int `column:"policy_type" json:"policyType"`

	//SubjectId 策略主体ID
	SubjectId string `column:"subject_id" json:"subjectId"`

	//ResourceCode 资源代码,辨识资源层级，比如有 aaa的权限，则认为也有 aaa:xxx的权限
	ResourceCode string `column:"resource_code" json:"resourceCode"`

	//Resource 资源内容
	Resource string `column:"resource" json:"resource"`

	//Action 资源使用行为
	Action string `column:"action" json:"action"`

	//Effect 策略行为 1-allow，2-deny
	Effect int `column:"effect" json:"effect"`

	//Priority 优先级
	Priority int `column:"priority" json:"priority"`

	//TenantId 租户ID
	TenantId string `column:"tenant_id" json:"tenantId"`

	//ResourceId 对应的资源ID
	ResourceId string `column:"resource_id" json:"resourceId"`

	//------------------数据库字段结束,自定义字段写在下面---------------//
	//如果查询的字段在column tag中没有找到,就会根据名称(不区分大小写,支持 _ 转驼峰)映射到struct的属性上

}

// GetTableName 获取表名称
// IEntityStruct 接口的方法,实体类需要实现!!!
func (entity *AuthorizePolicy) GetTableName() string {
	return AuthorizePolicyTableName
}

// GetPKColumnName 获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称
// 不支持联合主键,变通认为无主键,业务控制实现(艰难取舍)
// 如果没有主键,也需要实现这个方法, return "" 即可
// IEntityStruct 接口的方法,实体类需要实现!!!
func (entity *AuthorizePolicy) GetPKColumnName() string {
	//如果没有主键
	//return ""
	return "id"
}
