package authorization

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"xorm.io/xorm"

	"github.com/yockii/qscore/pkg/constant"
	"github.com/yockii/qscore/pkg/database"
	"github.com/yockii/qscore/pkg/logger"
)

type authorizationService struct {
	enforcer   *casbin.Enforcer
	superAdmin string
}

var defaultService *authorizationService

func SetSuperAdmin(admin string) {
	defaultService.superAdmin = admin
}

func Init() {
	defaultService = &authorizationService{
		superAdmin: constant.DefaultRoleName,
	}
	if err := defaultService.Initial(database.DB); err != nil {
		logger.Panicf("初始化默认权限系统失败，系统不应在无权限安全保护状态下运行", err)
	}
}

func (s *authorizationService) Initial(db *xorm.Engine) error {
	a, err := NewAdapter(db)
	if err != nil {
		return err
	}
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act, tenant")
	m.AddDef("p", "p", "sub, obj, act, eft, priority, tenant, resId")
	m.AddDef("g", "g", "_, _, _")
	m.AddDef("e", "e", "priority(p.eft) || deny")
	m.AddDef("m", "m", "g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && r.act == p.act && r.tenant == p.tenant || checkSuperAdmin(r.sub)")

	s.enforcer, err = casbin.NewEnforcer(m, a)
	if err != nil {
		return err
	}

	s.enforcer.AddFunction("checkSuperAdmin", func(arguments ...interface{}) (interface{}, error) {
		un := arguments[0].(string)
		return s.enforcer.HasRoleForUser(un, s.superAdmin)
	})
	return nil
}
func (s *authorizationService) AddSubjectResource(subject string, resourceTarget, action, tenantId, resourceId string) (bool, error) {
	return s.enforcer.AddPermissionForUser(subject, resourceTarget, action, "allow", "10", tenantId, resourceId)
}
func (s *authorizationService) RemoveSubjectResource(subject string, resourceTarget, action, tenantId, resourceId string) (bool, error) {
	return s.enforcer.DeletePermissionForUser(subject, resourceTarget, action, "allow", "10", tenantId, resourceId)
}
func (s *authorizationService) AddSubjectGroup(subject, group, tenantId string) (bool, error) {
	return s.enforcer.AddRoleForUser(subject, group, tenantId)
}
func (s *authorizationService) RemoveSubjectGroup(subject, group, tenantId string) (bool, error) {
	return s.enforcer.DeleteRoleForUser(subject, group, tenantId)
}
func (s *authorizationService) GetSubjectResourceIds(subject string, tenantId string) (isSuperAdmin bool, ids []string, err error) {
	isSuperAdmin, err = s.enforcer.HasRoleForUser(subject, s.superAdmin)
	if err != nil {
		return
	}
	if isSuperAdmin {
		return
	}
	var resources [][]string
	resources, err = s.enforcer.GetImplicitResourcesForUser(subject, tenantId)
	if err != nil {
		return
	}
	for _, r := range resources {
		ids = append(ids, r[6])
	}
	return
}
func (s *authorizationService) GetSubjectGroupIds(subject, tenantId string) ([]string, error) {
	roleIds, err := s.enforcer.GetRolesForUser(subject, tenantId)
	if err != nil {
		return nil, err
	}

	return roleIds, nil
}
func (s *authorizationService) CheckSubjectPermissions(subject, resource, action, tenantId string) bool {
	ok, _ := s.enforcer.Enforce(subject, resource, action, tenantId)
	return ok
}

func AddSubjectResource(subject string, resourceTarget, action, tenantId, resourceId string) (bool, error) {
	return defaultService.AddSubjectResource(subject, resourceTarget, action, tenantId, resourceId)
}
func RemoveSubjectResource(subject string, resourceTarget, action, tenantId, resourceId string) (bool, error) {
	return defaultService.RemoveSubjectResource(subject, resourceTarget, action, tenantId, resourceId)
}
func AddSubjectGroup(subject, group, tenantId string) (bool, error) {
	return defaultService.AddSubjectGroup(subject, group, tenantId)
}
func RemoveSubjectGroup(subject, group, tenantId string) (bool, error) {
	return defaultService.RemoveSubjectGroup(subject, group, tenantId)
}
func GetSubjectResourceIds(subject string, tenantId string) (isSuperAdmin bool, ids []string, err error) {
	return defaultService.GetSubjectResourceIds(subject, tenantId)
}
func GetSubjectGroupIds(subject, tenantId string) ([]string, error) {
	return defaultService.GetSubjectGroupIds(subject, tenantId)
}
func CheckSubjectPermissions(subject, resource, action, tenantId string) bool {
	return defaultService.CheckSubjectPermissions(subject, resource, action, tenantId)
}
