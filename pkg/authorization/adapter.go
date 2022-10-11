package authorization

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gitee.com/chunanyong/zorm"
	"github.com/casbin/casbin/v2/model"
	logger "github.com/sirupsen/logrus"

	"github.com/yockii/qscore/pkg/util"
)

type adapter struct{}

func NewAdapter() (*adapter, error) {
	return &adapter{}, nil
}

func (a *adapter) LoadPolicy(m model.Model) error {
	ctx := context.Background()
	var policies []*AuthorizePolicy

	finder := zorm.NewSelectFinder(AuthorizePolicyTableName)

	if err := zorm.Query(ctx, finder, &policies, nil); err != nil {
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

	var relations []*AuthorizeRelationship
	finder = zorm.NewSelectFinder(AuthorizeRelationshipTableName)
	if err := zorm.Query(ctx, finder, &relations, nil); err != nil {
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
	ctx := context.Background()

	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		_, err := zorm.Delete(ctx, &AuthorizeRelationship{})
		if err != nil {
			return nil, err
		}
		_, err = zorm.Delete(ctx, &AuthorizePolicy{})
		if err != nil {
			return nil, err
		}

		var policies []zorm.IEntityStruct
		var roles []zorm.IEntityStruct
		for policyType, ast := range m["p"] {
			for _, rule := range ast.Policy {
				if len(rule) == 7 {
					policy, err := a.parsePolicy(policyType, rule)
					if err != nil {
						logger.Error("策略规则必须是7个元素的数组", err)
						continue
					}
					policy.Id = util.GenerateDatabaseID()
					policies = append(policies, policy)
				}
			}
		}
		for relationType, ast := range m["g"] {
			for _, rule := range ast.Policy {
				role, err := a.parseRelation(relationType, rule)
				if err != nil {
					logger.Error("关系规则必须是2个元素的数组", err)
					continue
				}
				role.Id = util.GenerateDatabaseID()
				roles = append(roles, role)
			}
		}

		_, err = zorm.InsertSlice(ctx, policies)
		if err != nil {
			return nil, err
		}
		_, err = zorm.InsertSlice(ctx, roles)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (a *adapter) parseRelation(relationType string, rule []string) (*AuthorizeRelationship, error) {
	if len(rule) != 3 {
		return nil, errors.New("非法的父子关系规则数量，数量必须是3，实际" + strconv.Itoa(len(rule)))
	}
	gt := 1
	if s := relationType[1:]; s != "" {
		gt, _ = strconv.Atoi(relationType[1:])
	}
	role := &AuthorizeRelationship{
		Id:              util.GenerateDatabaseID(),
		RelationType:    gt,
		SubjectId:       rule[0],
		ParentSubjectId: rule[1],
		TenantId:        rule[2],
	}
	return role, nil
}

func (a *adapter) parsePolicy(policyType string, rule []string) (*AuthorizePolicy, error) {
	if len(rule) != 8 {
		return nil, errors.New("非法的策略信息，策略字段必须为8， 实际" + strconv.Itoa(len(rule)))
	}
	pt := 1
	if s := policyType[1:]; s != "" {
		pt, _ = strconv.Atoi(policyType[1:])
	}
	effect := 1
	if rule[4] == "deny" {
		effect = 2
	}
	priority, _ := strconv.Atoi(rule[5])
	policy := &AuthorizePolicy{
		Id:           util.GenerateDatabaseID(),
		PolicyType:   pt,
		SubjectId:    rule[0],
		ResourceCode: rule[1],
		Resource:     rule[2],
		Action:       rule[3],
		Effect:       effect,
		Priority:     priority,
		TenantId:     rule[6],
		ResourceId:   rule[7],
	}
	return policy, nil
}

// AddPolicy adds a policy rule to the storage.
func (a *adapter) AddPolicy(sec string, ptype string, rule []string) error {
	ctx := context.Background()
	if sec == "p" {
		policy, err := a.parsePolicy(ptype, rule)
		if err != nil {
			return err
		}
		_, err = zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
			return zorm.Insert(ctx, policy)
		})
		return err
	} else if sec == "g" {
		role, err := a.parseRelation(ptype, rule)
		if err != nil {
			return err
		}
		_, err = zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
			return zorm.Insert(ctx, role)
		})
		return err
	} else {
		return errors.New("策略添加的sec非法! ")
	}
}

// RemovePolicy removes a policy rule from the storage.
func (a *adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	ctx := context.Background()
	if sec == "p" {
		policy, err := a.parsePolicy(ptype, rule)
		if err != nil {
			return err
		}
		_, err = zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
			return zorm.Delete(ctx, policy)
		})
		return err
	} else if sec == "g" {
		role, err := a.parseRelation(ptype, rule)
		if err != nil {
			return err
		}
		_, err = zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
			return zorm.Delete(ctx, role)
		})
		return err
	} else {
		return errors.New("要删除的策略入参非法! ")
	}
}

func (a *adapter) RemovePolicies(sec string, ptype string, rules [][]string) error {
	ctx := context.Background()
	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		for _, rule := range rules {
			if sec == "p" {
				policy, err := a.parsePolicy(ptype, rule)
				if err != nil {
					return nil, err
				}
				_, err = zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
					return zorm.Delete(ctx, policy)
				})
				return nil, err
			} else if sec == "g" {
				role, err := a.parseRelation(ptype, rule)
				if err != nil {
					return nil, err
				}
				_, err = zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
					return zorm.Delete(ctx, role)
				})
				return nil, err
			} else {
				return nil, errors.New("要删除的策略入参非法! ")
			}
		}
		return nil, nil
	})

	return err
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	ctx := context.Background()
	idx := fieldIndex + len(fieldValues)
	if sec == "p" {
		pt := 1
		if s := ptype[1:]; s != "" {
			pt, _ = strconv.Atoi(ptype[1:])
		}
		policy := AuthorizePolicy{
			PolicyType: pt,
		}
		if fieldIndex <= 0 && idx > 0 {
			policy.SubjectId = fieldValues[0-fieldIndex]
		}
		if fieldIndex <= 1 && idx > 1 {
			policy.ResourceCode = fieldValues[1-fieldIndex]
		}
		if fieldIndex <= 2 && idx > 2 {
			policy.Resource = fieldValues[2-fieldIndex]
		}
		if fieldIndex <= 3 && idx > 3 {
			policy.Action = fieldValues[3-fieldIndex]
		}
		if fieldIndex <= 4 && idx > 4 {
			es := fieldValues[4-fieldIndex]
			effect := 1
			if es == "deny" {
				effect = 2
			}
			policy.Effect = effect
		}
		if fieldIndex <= 5 && idx > 5 {
			policy.Priority, _ = strconv.Atoi(fieldValues[5-fieldIndex])
		}
		if fieldIndex <= 6 && idx > 6 {
			policy.TenantId = fieldValues[6-fieldIndex]
		}
		if fieldIndex <= 7 && idx > 7 {
			policy.ResourceId = fieldValues[7-fieldIndex]
		}
		_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
			return zorm.Delete(ctx, &policy)
		})
		return err
	} else if sec == "g" {
		gt := 1
		if s := ptype[1:]; s != "" {
			gt, _ = strconv.Atoi(ptype[1:])
		}
		role := AuthorizeRelationship{
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
		_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
			return zorm.Delete(ctx, &role)
		})
		return err
	} else {
		return errors.New("要删除的策略入参非法! ")
	}
}
