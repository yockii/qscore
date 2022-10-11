package authorization

import (
	"gitee.com/chunanyong/zorm"
)

// AuthorizeRelationshipTableName 表名常量,方便直接调用
const AuthorizeRelationshipTableName = "t_authorize_relationship"

// AuthorizeRelationship
type AuthorizeRelationship struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	zorm.EntityStruct

	//Id []
	Id string `column:"id" json:"id"`

	//RelationType 关系类型 1-g 2-g2 3-g3
	RelationType int `column:"relation_type" json:"relationType"`

	//SubjectId 主体ID
	SubjectId string `column:"subject_id" json:"subjectId"`

	//TenantId 租户ID
	TenantId string `column:"tenant_id" json:"tenantId"`

	//ParentSubjectId 继承主体ID
	ParentSubjectId string `column:"parent_subject_id" json:"parentSubjectId"`

	//------------------数据库字段结束,自定义字段写在下面---------------//
	//如果查询的字段在column tag中没有找到,就会根据名称(不区分大小写,支持 _ 转驼峰)映射到struct的属性上

}

// GetTableName 获取表名称
// IEntityStruct 接口的方法,实体类需要实现!!!
func (entity *AuthorizeRelationship) GetTableName() string {
	return AuthorizeRelationshipTableName
}

// GetPKColumnName 获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称
// 不支持联合主键,变通认为无主键,业务控制实现(艰难取舍)
// 如果没有主键,也需要实现这个方法, return "" 即可
// IEntityStruct 接口的方法,实体类需要实现!!!
func (entity *AuthorizeRelationship) GetPKColumnName() string {
	//如果没有主键
	//return ""
	return "id"
}
