package authorization

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/casbin/casbin/v2/model"
	"xorm.io/xorm"

	"github.com/yockii/qscore/pkg/logger"
	"github.com/yockii/qscore/pkg/util"
)

type adapter struct {
	engine *xorm.Engine
}

type CasbinPolicy struct {
	Id         string `json:"id" xorm:"pk varchar(50)"`
	PolicyType int    `json:"policyType" xorm:"comment('策略类型 1-p, 2-p2')"`  // p  p2
	SubjectId  string `json:"subjectId" xorm:"comment('策略主体ID')"`           // 主体
	Resource   string `json:"resource" xorm:"comment('资源内容')"`              // 资源
	Action     string `json:"action" xorm:"comment('资源使用行为')"`              // 方法
	Effect     int    `json:"effect" xorm:"comment('策略行为 1-allow，2-deny')"` // 1 allow / 2 deny
	Priority   int    `json:"priority" xorm:"comment('优先级')"`               // 优先级
	TenantId   string `json:"tenantId" xorm:"comment('租户ID')"`
	ResourceId string `json:"resourceId" xorm:"varchar(50) comment('对应的资源ID')"` // 资源ID
}

type CasbinRelationship struct {
	Id              string `json:"id,omitempty" xorm:"pk varchar(50)"`
	RelationType    int    `json:"relationType,omitempty" xorm:"comment('关系类型 1-g 2-g2')"`
	SubjectId       string `json:"subjectId,omitempty" xorm:"comment('主体ID')"`
	ParentSubjectId string `json:"parentSubjectId,omitempty" xorm:"comment('继承主体ID')"`
	TenantId        string `json:"tenantId" xorm:"comment('租户ID')"`
}

func NewAdapter(engine *xorm.Engine) (*adapter, error) {
	if err := engine.Sync2(CasbinPolicy{}, CasbinRelationship{}); err != nil {
		return nil, err
	}
	return &adapter{engine: engine}, nil
}

func (a *adapter) LoadPolicy(m model.Model) error {
	var policies []*CasbinPolicy
	if err := a.engine.Find(&policies); err != nil {
		return err
	}
	for _, policy := range policies {
		policyType := "p"
		if policy.PolicyType != 1 {
			policyType = fmt.Sprintf("p%d", policy.PolicyType)
		}
		effect := "allow"
		if policy.Effect == 2 {
			effect = "deny"
		}
		tokens := []string{
			policy.SubjectId,
			policy.Resource,
			policy.Action,
			effect,
			strconv.Itoa(policy.Priority),
			policy.TenantId,
			policy.ResourceId,
		}
		mpp := m["p"][policyType]
		mpp.Policy = append(mpp.Policy, tokens)
		mpp.PolicyMap[strings.Join(tokens, model.DefaultSep)] = len(mpp.Policy) - 1
	}

	var relations []*CasbinRelationship
	if err := a.engine.Find(&relations); err != nil {
		return err
	}
	for _, relation := range relations {
		relationType := "g"
		if relation.RelationType != 1 {
			relationType = fmt.Sprintf("g%d", relation.RelationType)
		}
		tokens := []string{
			relation.SubjectId,
			relation.ParentSubjectId,
			relation.TenantId,
		}
		mgg := m["g"][relationType]
		mgg.Policy = append(mgg.Policy, tokens)
		mgg.PolicyMap[strings.Join(tokens, model.DefaultSep)] = len(mgg.Policy) - 1
	}
	return nil
}

// SavePolicy saves all policy rules to the storage.
func (a *adapter) SavePolicy(m model.Model) error {
	sess := a.engine.NewSession()
	defer sess.Close()
	sess.Delete(&CasbinRelationship{})
	sess.Delete(&CasbinPolicy{})

	var policies []*CasbinPolicy
	var roles []*CasbinRelationship
	for policyType, ast := range m["p"] {
		for _, rule := range ast.Policy {
			if len(rule) == 7 {
				policy, err := a.parsePolicy(policyType, rule)
				if err != nil {
					logger.Error("策略规则必须是7个元素的数组", err)
					continue
				}
				policies = append(policies, policy)
			}
		}
	}
	for relationType, ast := range m["g"] {
		for _, rule := range ast.Policy {
			role, err := a.parseRelation(relationType, rule)
			if err != nil {
				logger.Error("关系规则必须是3个元素的数组", err)
				continue
			}
			roles = append(roles, role)
		}
	}

	_, err := sess.Insert(policies)
	if err != nil {
		return err
	}
	_, err = sess.Insert(roles)
	if err != nil {
		return err
	}

	return sess.Commit()
}

func (a *adapter) parseRelation(relationType string, rule []string) (*CasbinRelationship, error) {
	if len(rule) != 3 {
		return nil, errors.New("非法的父子关系规则数量，数量必须是3，实际" + strconv.Itoa(len(rule)))
	}
	gt := 1
	if s := relationType[1:]; s != "" {
		gt, _ = strconv.Atoi(relationType[1:])
	}
	role := &CasbinRelationship{
		Id:              util.GenerateDatabaseID(),
		RelationType:    gt,
		SubjectId:       rule[0],
		ParentSubjectId: rule[1],
		TenantId:        rule[2],
	}
	return role, nil
}

func (a *adapter) parsePolicy(policyType string, rule []string) (*CasbinPolicy, error) {
	if len(rule) != 7 {
		return nil, errors.New("invalid policy rule ")
	}
	pt := 1
	if s := policyType[1:]; s != "" {
		pt, _ = strconv.Atoi(policyType[1:])
	}
	effect := 1
	if rule[3] == "deny" {
		effect = 2
	}
	priority, _ := strconv.Atoi(rule[4])
	policy := &CasbinPolicy{
		Id:         util.GenerateDatabaseID(),
		PolicyType: pt,
		SubjectId:  rule[0],
		Resource:   rule[1],
		Action:     rule[2],
		Effect:     effect,
		Priority:   priority,
		TenantId:   rule[5],
		ResourceId: rule[6],
	}
	return policy, nil
}

// AddPolicy adds a policy rule to the storage.
func (a *adapter) AddPolicy(sec string, ptype string, rule []string) error {
	if sec == "p" {
		policy, err := a.parsePolicy(ptype, rule)
		if err != nil {
			return err
		}
		_, err = a.engine.InsertOne(policy)
		return err
	} else if sec == "g" {
		role, err := a.parseRelation(ptype, rule)
		if err != nil {
			return err
		}
		_, err = a.engine.InsertOne(role)
		return err
	} else {
		return errors.New("策略添加的sec非法! ")
	}
}

// RemovePolicy removes a policy rule from the storage.
func (a *adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	if sec == "p" {
		policy, err := a.parsePolicy(ptype, rule)
		if err != nil {
			return err
		}
		_, err = a.engine.Delete(&policy)
		return err
	} else if sec == "g" {
		role, err := a.parseRelation(ptype, rule)
		if err != nil {
			return err
		}
		_, err = a.engine.Delete(&role)
		return err
	} else {
		return errors.New("要删除的策略入参非法! ")
	}
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	idx := fieldIndex + len(fieldValues)
	if sec == "p" {
		pt := 1
		if s := ptype[1:]; s != "" {
			pt, _ = strconv.Atoi(ptype[1:])
		}
		policy := CasbinPolicy{
			PolicyType: pt,
		}
		if fieldIndex <= 0 && idx > 0 {
			policy.SubjectId = fieldValues[0-fieldIndex]
		}
		if fieldIndex <= 1 && idx > 1 {
			policy.Resource = fieldValues[1-fieldIndex]
		}
		if fieldIndex <= 2 && idx > 2 {
			policy.Action = fieldValues[2-fieldIndex]
		}
		if fieldIndex <= 3 && idx > 3 {
			es := fieldValues[3-fieldIndex]
			effect := 1
			if es == "deny" {
				effect = 2
			}
			policy.Effect = effect
		}
		if fieldIndex <= 3 && idx > 3 {
			policy.Priority, _ = strconv.Atoi(fieldValues[4-fieldIndex])
		}
		if fieldIndex <= 4 && idx > 4 {
			policy.TenantId = fieldValues[5-fieldIndex]
		}
		if fieldIndex <= 5 && idx > 5 {
			policy.ResourceId = fieldValues[6-fieldIndex]
		}
		_, err := a.engine.Delete(&policy)
		return err
	} else if sec == "g" {
		gt := 1
		if s := ptype[1:]; s != "" {
			gt, _ = strconv.Atoi(ptype[1:])
		}
		role := CasbinRelationship{
			RelationType: gt,
		}
		if fieldIndex <= 0 && idx > 0 {
			role.SubjectId = fieldValues[0-fieldIndex]
		}
		if fieldIndex <= 1 && idx > 1 {
			role.ParentSubjectId = fieldValues[1-fieldIndex]
		}
		if fieldIndex <= 2 && idx > 2 {
			role.TenantId = fieldValues[2-fieldIndex]
		}
		_, err := a.engine.Delete(&role)
		return err
	} else {
		return errors.New("要删除的策略入参非法! ")
	}
}
