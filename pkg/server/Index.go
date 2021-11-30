package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/yockii/qscore/pkg/domain"
)

type WebApp struct {
	app *fiber.App
}

func initFiberParser() {
	customDateTime := fiber.ParserType{
		Customtype: domain.DateTime{},
		Converter:  domain.DateTimeConverter,
	}
	fiber.SetParserDecoder(fiber.ParserConfig{
		IgnoreUnknownKeys: true,
		ParserType:        []fiber.ParserType{customDateTime},
		ZeroEmpty:         true,
	})
}

func InitFiber() *WebApp {
	initFiberParser()
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(e interface{}) {
		},
	}))
	app.Use(cors.New())

	return &WebApp{app}
}

func (a *WebApp) Group(prefix string, needLogin, needRouterPermission bool) fiber.Router {
	var handlers []fiber.Handler
	if needLogin {
		handlers = append(handlers, Jwtware)
	}
	if needRouterPermission {
		handlers = append(handlers, RequireRouterPermission())
	}

	return a.app.Group(prefix, handlers...)
}
