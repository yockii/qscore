package common

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/qscore/pkg/database"
	"github.com/yockii/qscore/pkg/server"
	"github.com/yockii/qscore/pkg/util"
	"gorm.io/gorm"
	"time"
)

type Service[T database.Model] interface {
	Add(instance T, tx ...*gorm.DB) (duplicated bool, err error)
	Update(instance T, tx ...*gorm.DB) (count int64, err error)
	Delete(instance T, tx ...*gorm.DB) (count int64, err error)
	Instance(instance T) (result T, err error)
	List(condition T, paginate *server.Paginate[T], orderBy string, tcList map[string]*server.TimeCondition) (err error)

	Model() T
}

type BaseService[T database.Model] struct {
	Service[T]
}

func (*BaseService[T]) Add(instance T, tx ...*gorm.DB) (duplicated bool, err error) {
	if lackedFields := instance.AddRequired(); lackedFields != "" {
		err = errors.New(lackedFields + " is required")
		return
	}
	var c int64
	if cdm := instance.CheckDuplicatedModel(); cdm != nil {
		err = database.DB.Model(instance).Where(cdm).Count(&c).Error
		if err != nil {
			logger.Errorln(err)
			return
		}
	}
	if c > 0 {
		duplicated = true
		return
	}

	instance.InitDefaultFields()

	if len(tx) > 0 {
		if err = tx[0].Create(instance).Error; err != nil {
			logger.Errorln(err)
			return
		}
		return
	}

	if err = database.DB.Create(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	return
}

func (s *BaseService[T]) Update(instance T, tx ...*gorm.DB) (count int64, err error) {
	if lackedFields := instance.UpdateRequired(); lackedFields != "" {
		err = errors.New(lackedFields + " is required")
		return
	}
	if len(tx) > 0 {
		err = tx[0].Model(s.Model()).Where(instance.UpdateConditionModel()).Updates(instance.UpdateModel()).Error
	} else {
		err = database.DB.Model(s.Model()).Where(instance.UpdateConditionModel()).Updates(instance.UpdateModel()).Error
	}
	if err != nil {
		logger.Errorln(err)
	}
	return
}

func (s *BaseService[T]) Delete(instance T, tx ...*gorm.DB) (count int64, err error) {
	if lackedFields := instance.DeleteRequired(); lackedFields != "" {
		err = errors.New(lackedFields + " is required")
		return
	}
	var result *gorm.DB
	if len(tx) > 0 {
		result = tx[0].Model(s.Model()).Delete(instance)
	} else {
		result = database.DB.Model(s.Model()).Delete(instance)
	}
	err = result.Error
	count = result.RowsAffected
	if err != nil {
		logger.Errorln(err)
	}
	return
}

func (s *BaseService[T]) Instance(instance T) (result T, err error) {
	if lackedFields := instance.InstanceRequired(); lackedFields != "" {
		err = errors.New(lackedFields + " is required")
		return
	}
	result = s.Model()
	err = database.DB.Where(instance).First(&result).Error
	if err != nil {
		logger.Errorln(err)
	}
	return
}

func (s *BaseService[T]) List(condition T, paginate *server.Paginate[T], orderBy string, tcList map[string]*server.TimeCondition) (err error) {
	tx := database.DB.Model(s.Model())
	if paginate == nil {
		return errors.New("paginate不能为空")
	}
	if paginate.Limit > -1 {
		tx = tx.Limit(paginate.Limit)
	}
	if paginate.Offset > -1 {
		tx = tx.Offset(paginate.Offset)
	}
	if orderBy != "" {
		tx = tx.Order(util.SnakeString(orderBy))
	}
	for tc, tr := range tcList {
		if tc != "" && tr != nil {
			if !tr.Start.IsZero() && !tr.End.IsZero() {
				tx = tx.Where(tc+" between ? and ?", time.Time(tr.Start).UnixMilli(), time.Time(tr.End).UnixMilli())
			} else if tr.Start.IsZero() && !tr.End.IsZero() {
				tx = tx.Where(tc+" <= ?", time.Time(tr.End).UnixMilli())
			} else if !tr.Start.IsZero() && tr.End.IsZero() {
				tx = tx.Where(tc+" > ?", time.Time(tr.Start).UnixMilli())
			}
		}
	}
	fuzzyQueryMap := condition.FuzzyQueryMap()
	for k, v := range fuzzyQueryMap {
		tx = tx.Where(k+" like ?", v)
	}
	if condition.ListOmits() != "" {
		tx = tx.Omit(condition.ListOmits())
	}
	err = tx.Find(&paginate.Items, condition.ExactMatchModel()).Offset(-1).Limit(-1).Count(&paginate.Total).Error
	if err != nil {
		logger.Errorln(err)
	}
	return
}
