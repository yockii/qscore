package common

import (
	"github.com/gofiber/fiber/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/qscore/pkg/server"
)

type RouterController interface {
	// 初始化管理路由
	InitManage()
}

type Controller[T Model, D BaseDomain[T]] interface {
	NewModel() T
	NewDomain() D
	Add(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	Detail(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
	GetService() Service[T]
}

type BaseController[T Model, D BaseDomain[T]] struct {
	Controller[T, D]
}

func (c *BaseController[T, D]) Add(ctx *fiber.Ctx) error {
	instance := c.NewModel()
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	if lackedFields := instance.AddRequired(); lackedFields != "" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + ": " + lackedFields,
		})
	}
	duplicated, err := c.Controller.GetService().Add(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	if duplicated {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDuplicated,
			Msg:  server.ResponseMsgDuplicated,
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: instance,
	})
}

func (c *BaseController[T, D]) Update(ctx *fiber.Ctx) error {
	instance := c.NewModel()

	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if lackedFields := instance.UpdateRequired(); lackedFields != "" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + ": " + lackedFields,
		})
	}
	count, err := c.GetService().Update(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: count > 0,
	})
}

func (c *BaseController[T, D]) Delete(ctx *fiber.Ctx) error {
	instance := c.NewModel()

	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if lackedFields := instance.DeleteRequired(); lackedFields != "" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + ": " + lackedFields,
		})
	}
	count, err := c.GetService().Delete(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: count > 0,
	})
}

func (c *BaseController[T, D]) Detail(ctx *fiber.Ctx) (err error) {
	instance := c.NewModel()

	if err = ctx.QueryParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if lackedFields := instance.InstanceRequired(); lackedFields != "" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + ": " + lackedFields,
		})
	}
	if instance, err = c.GetService().Instance(instance); err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: instance,
	})
}

func (c *BaseController[T, D]) List(ctx *fiber.Ctx) error {
	domainWithModel := c.NewDomain()

	if err := ctx.QueryParser(domainWithModel); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	paginate := new(server.Paginate[T])
	if err := ctx.QueryParser(paginate); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if paginate.Limit == 0 {
		paginate.Limit = 10
	}

	if err := c.GetService().List(domainWithModel.GetModel(), paginate, domainWithModel.GetOrderBy(), domainWithModel.GetTimeConditionList()); err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: paginate,
	})
}
