package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/yockii/qscore/pkg/domain"
	"github.com/yockii/qscore/pkg/logger"
)

type webApp struct {
	app *fiber.App
}

var defaultApp *webApp

func init() {
	initFiberParser()
	defaultApp = InitWebApp()
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

func InitWebApp() *webApp {
	initFiberParser()
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(ctx *fiber.Ctx, e interface{}) {
			logger.Error(e)
		},
	}))
	app.Use(cors.New())

	return &webApp{app}
}

func (a *webApp) Group(prefix string, needLogin, needRouterPermission bool) fiber.Router {
	var handlers []fiber.Handler
	if needLogin {
		handlers = append(handlers, Jwtware)
	}
	if needRouterPermission {
		handlers = append(handlers, RequireRouterPermission())
	}

	return a.app.Group(prefix, handlers...)
}
func (a *webApp) All(path string, handlers ...fiber.Handler) fiber.Router {
	return a.app.All(path, handlers...)
}
func (a *webApp) Get(path string, handlers ...fiber.Handler) fiber.Router {
	return a.app.Get(path, handlers...)
}
func (a *webApp) Put(path string, handlers ...fiber.Handler) fiber.Router {
	return a.app.Put(path, handlers...)
}
func (a *webApp) Post(path string, handlers ...fiber.Handler) fiber.Router {
	return a.app.Post(path, handlers...)
}
func (a *webApp) Delete(path string, handlers ...fiber.Handler) fiber.Router {
	return a.app.Delete(path, handlers...)
}
func (a *webApp) Start(addr string) error {
	return a.app.Listen(addr)
}
func (a *webApp) Shutdown() error {
	return a.app.Shutdown()
}

//StandardRouter 标准路由，需要登录、校验权限
func StandardRouter(prefix string, add, update, delete, get, paginate fiber.Handler) fiber.Router {
	return StandardVersionRouter("v1", prefix, add, update, delete, get, paginate)
}

func StandardVersionRouter(version, prefix string, add, update, delete, get, paginate fiber.Handler) fiber.Router {
	g := defaultApp.Group(fmt.Sprintf("/api/%s%s", version, prefix), true, true)
	g.Post("/", add)
	g.Put("/", update)
	g.Delete("/", delete)
	g.Get("/instance", get)
	g.Get("/list", paginate)
	return g
}

func Group(prefix string, needLogin, needRouterPermission bool) fiber.Router {
	return defaultApp.Group(prefix, needLogin, needRouterPermission)
}
func All(path string, handlers ...fiber.Handler) fiber.Router {
	return defaultApp.All(path, handlers...)
}
func Get(path string, handlers ...fiber.Handler) fiber.Router {
	return defaultApp.Get(path, handlers...)
}
func Put(path string, handlers ...fiber.Handler) fiber.Router {
	return defaultApp.Put(path, handlers...)
}
func Post(path string, handlers ...fiber.Handler) fiber.Router {
	return defaultApp.Post(path, handlers...)
}
func Delete(path string, handlers ...fiber.Handler) fiber.Router {
	return defaultApp.Delete(path, handlers...)
}
func Start(addr string) error {
	return defaultApp.Start(addr)
}
func Shutdown() error {
	return defaultApp.Shutdown()
}
