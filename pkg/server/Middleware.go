package server

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gomodule/redigo/redis"

	"github.com/yockii/qscore/pkg/authorization"
	"github.com/yockii/qscore/pkg/cache"

	"github.com/yockii/qscore/pkg/constant"
	"github.com/yockii/qscore/pkg/logger"
)

var Jwtware = jwtware.New(jwtware.Config{
	SigningKey:    []byte(constant.JWT_SECRET),
	ContextKey:    constant.JWT_CONTEXT,
	SigningMethod: "HS256",
	TokenLookup:   "header:Authorization,cookie:token",
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		if err.Error() == "Missing or malformed JWT" {
			return c.Status(fiber.StatusBadRequest).SendString("无效的token信息")
		} else {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired Authorization Token")
		}
	},
	SuccessHandler: func(c *fiber.Ctx) error {
		// 从jwt获取用户信息
		jwtToken := c.Locals(constant.JWT_CONTEXT).(*jwt.Token)
		claims := jwtToken.Claims.(jwt.MapClaims)
		uid := claims["uid"].(string)
		sid := claims["sid"].(string)
		tenantId, tenantOk := claims["tenantId"].(string)

		if cache.Enabled() {
			rConn := cache.Get()
			defer rConn.Close()
			cachedUid, err := redis.String(rConn.Do("GET", cache.Prefix+":"+constant.AppSid+":"+sid))
			if err != nil {
				if err != redis.ErrNil {
					logger.Error(err)
				}
				return c.Status(fiber.StatusUnauthorized).SendString("token信息失效")
			}
			if cachedUid != uid {
				logger.Error(err)
				return c.Status(fiber.StatusUnauthorized).SendString("token信息失效")
			}
		}

		c.Locals("userId", uid)
		if tenantOk {
			c.Locals("tenantId", tenantId)
		}
		return c.Next()
	},
})

func RequireRouterPermission() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		path := ctx.Path()
		method := ctx.Method()
		subject := ctx.Locals("userId").(string)
		if subject == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		if authorization.CheckSubjectPermissions(subject, path, method, "") {
			return ctx.Next()
		}
		return ctx.SendStatus(fiber.StatusForbidden)
	}
}
