package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/yockii/qscore/pkg/config"
	"github.com/yockii/qscore/pkg/database"
	"net"

	"github.com/gofiber/fiber/v2/middleware/recover"
	logger "github.com/sirupsen/logrus"
)

type WebApp struct {
	app        *fiber.App
	ViewEngine fiber.Views
}

var DefaultWebApp *WebApp

func init() {
	initServerDefault()
}

func InitServer() {
	if config.GetString("server.viewsDir") != "" {
		extension := config.GetString("server.viewExtension")
		if extension == "" {
			extension = ".html"
		}
		DefaultWebApp = InitWebApp(html.New(config.GetString("server.viewsDir"), extension))
	} else {
		DefaultWebApp = InitWebApp(nil)
	}
}

func InitServerWithViews(views fiber.Views) {
	DefaultWebApp = InitWebApp(views)
}

func initServerDefault() {
	config.DefaultInstance.SetDefault("server.port", 13579)
}

func initFiberParser() {
	customDateTime := fiber.ParserType{
		Customtype: database.DateTime{},
		Converter:  database.DateTimeConverter,
	}
	fiber.SetParserDecoder(fiber.ParserConfig{
		IgnoreUnknownKeys: true,
		ParserType:        []fiber.ParserType{customDateTime},
		ZeroEmpty:         true,
	})
}

func InitWebApp(views fiber.Views) *WebApp {
	initFiberParser()
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Views:                 views,
		BodyLimit:             1 * 1024 * 1024 * 1024,
	})
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(ctx *fiber.Ctx, e interface{}) {
			logger.Error(e)
		},
	}))
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
		AllowCredentials: true,
	}))

	return &WebApp{
		app:        app,
		ViewEngine: views,
	}
}

func (a *WebApp) Listener(ln net.Listener) error {
	return a.app.Listener(ln)
}
func (a *WebApp) Static(dir string) {
	a.app.Static("/", dir, fiber.Static{
		Compress: true,
	})
}

func (a *WebApp) Use(args ...interface{}) fiber.Router {
	return a.app.Use(args...)
}
func (a *WebApp) Group(path string, handlers ...fiber.Handler) fiber.Router {
	return a.app.Group(path, handlers...)
}
func (a *WebApp) All(path string, handlers ...fiber.Handler) fiber.Router {
	return a.app.All(path, handlers...)
}
func (a *WebApp) Get(path string, handlers ...fiber.Handler) fiber.Router {
	return a.app.Get(path, handlers...)
}
func (a *WebApp) Put(path string, handlers ...fiber.Handler) fiber.Router {
	return a.app.Put(path, handlers...)
}
func (a *WebApp) Post(path string, handlers ...fiber.Handler) fiber.Router {
	return a.app.Post(path, handlers...)
}
func (a *WebApp) Delete(path string, handlers ...fiber.Handler) fiber.Router {
	return a.app.Delete(path, handlers...)
}
func (a *WebApp) Start(addr string) error {
	return a.app.Listen(addr)
}
func (a *WebApp) Shutdown() error {
	return a.app.Shutdown()
}

func Listener(ln net.Listener) error {
	return DefaultWebApp.Listener(ln)
}
func Static(dir string) {
	DefaultWebApp.Static(dir)
}

func Use(args ...interface{}) fiber.Router {
	return DefaultWebApp.Use(args...)
}
func Group(path string, handlers ...fiber.Handler) fiber.Router {
	return DefaultWebApp.Group(path, handlers...)
}
func All(path string, handlers ...fiber.Handler) fiber.Router {
	return DefaultWebApp.All(path, handlers...)
}
func Get(path string, handlers ...fiber.Handler) fiber.Router {
	return DefaultWebApp.Get(path, handlers...)
}
func Put(path string, handlers ...fiber.Handler) fiber.Router {
	return DefaultWebApp.Put(path, handlers...)
}
func Post(path string, handlers ...fiber.Handler) fiber.Router {
	return DefaultWebApp.Post(path, handlers...)
}
func Delete(path string, handlers ...fiber.Handler) fiber.Router {
	return DefaultWebApp.Delete(path, handlers...)
}
func Start() error {
	return DefaultWebApp.Start(":" + config.GetString("server.port"))
}
func Shutdown() error {
	return DefaultWebApp.Shutdown()
}
