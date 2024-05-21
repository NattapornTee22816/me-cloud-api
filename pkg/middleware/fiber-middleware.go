package middleware

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/keyauth"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	recover2 "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/utils/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"runtime/debug"
	"strings"
	"time"
)

var (
	log, _ = zap.NewProduction()
)

type FiberMiddleware struct {
	app     *fiber.App
	db      clickhouse.Conn
	session *redis.Client
	cache   *redis.Client
}

func NewFiberMiddleware(app *fiber.App, db clickhouse.Conn, cacheSession *redis.Client, cache *redis.Client) *FiberMiddleware {
	return &FiberMiddleware{
		app,
		db,
		cacheSession,
		cache,
	}
}

func (m *FiberMiddleware) UseRecovery() *FiberMiddleware {
	m.app.Use(recover2.New(recover2.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c fiber.Ctx, e any) {
			log.Error(fmt.Sprintf("%s", debug.Stack()))
		},
	}))

	return m
}

func (m *FiberMiddleware) UseCors() *FiberMiddleware {
	m.app.Use(cors.New(cors.Config{
		MaxAge: 3600,
		AllowHeaders: strings.Join([]string{
			fiber.HeaderAcceptLanguage,
			fiber.HeaderCacheControl,
			fiber.HeaderContentLength,
			fiber.HeaderContentType,
			fiber.HeaderXRequestedWith,
			fiber.HeaderAuthorization,
			fiber.HeaderXRequestID,
			fiber.HeaderOrigin,
			fiber.HeaderAccessControlAllowHeaders,
			fiber.HeaderAccessControlAllowMethods,
			fiber.HeaderAccessControlAllowOrigin,
			fiber.HeaderAccessControlAllowCredentials,
			fiber.HeaderAccessControlAllowPrivateNetwork,
		}, ","),
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodPatch,
			fiber.MethodDelete,
		}, ","),
		ExposeHeaders: strings.Join([]string{
			fiber.HeaderContentLength,
			fiber.HeaderContentType,
			fiber.HeaderXRequestedWith,
			fiber.HeaderAuthorization,
			fiber.HeaderXRequestID,
		}, ","),
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			// TODO validate origin from db
			return true
		},
	}))

	return m
}

func (m *FiberMiddleware) UseCompress() *FiberMiddleware {
	m.app.Use(compress.New(compress.Config{
		Level: compress.LevelDefault,
	}))

	return m
}

func (m *FiberMiddleware) UseHelmet() *FiberMiddleware {
	m.app.Use(helmet.New())

	return m
}

func (m *FiberMiddleware) UseKeyAuth() *FiberMiddleware {
	m.app.Use(keyauth.New(keyauth.Config{
		KeyLookup: "header:X-Api-Key",
		Next: func(ctx fiber.Ctx) bool {
			return ctx.Path() == healthcheck.DefaultLivenessEndpoint
		},
		Validator: func(ctx fiber.Ctx, key string) (bool, error) {

			// TODO validate api-key

			return true, nil
		},
		SuccessHandler: func(ctx fiber.Ctx) error {
			return ctx.Next()
		},
		ErrorHandler: func(ctx fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    "401",
				"message": "Invalid or expired API Key",
			})
		},
	}))

	return m
}

func (m *FiberMiddleware) UseLimiter(limit int) *FiberMiddleware {
	m.app.Use(limiter.New(limiter.Config{
		Max:        limit,
		Expiration: time.Minute,
		Next: func(c fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		LimiterMiddleware: limiter.SlidingWindow{},
		LimitReached: func(ctx fiber.Ctx) error {
			return ctx.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"code":    "429",
				"message": "you are too many requests, please wait and try again later.",
			})
		},
	}))

	return m
}

func (m *FiberMiddleware) UseLogger() *FiberMiddleware {
	m.app.Use(logger.New(logger.Config{
		TimeZone:      "Asia/Bangkok",
		TimeFormat:    "2006-01-02 15:04:05",
		DisableColors: true,
		Format:        "[${time}] ${ip} ${status} - ${latency} ${method} ${url}",
	}))

	return m
}

func (m *FiberMiddleware) UseRequestID() *FiberMiddleware {
	m.app.Use(requestid.New(requestid.Config{
		Generator: func() string {
			return utils.UUIDv4()
		},
	}))

	return m
}

func (m *FiberMiddleware) UseSession() *FiberMiddleware {
	m.app.Use(session.New(session.Config{}))

	return m
}

func (m *FiberMiddleware) UseHealthCheck() *FiberMiddleware {
	m.app.Get(healthcheck.DefaultLivenessEndpoint, healthcheck.NewHealthChecker(healthcheck.Config{
		Probe: func(ctx fiber.Ctx) bool {
			// check database live
			if err := m.db.Ping(context.Background()); err != nil {
				return false
			}
			if err := m.session.Ping(context.Background()).Err(); err != nil {
				return false
			}
			if err := m.cache.Ping(context.Background()).Err(); err != nil {
				return false
			}

			return true
		},
	}))

	return m
}

func (m *FiberMiddleware) UseLogging() *FiberMiddleware {
	m.app.Use(logger.New(logger.Config{
		TimeZone:      "Asia/Bangkok",
		DisableColors: true,
		LoggerFunc: func(c fiber.Ctx, data *logger.Data, cfg logger.Config) error {
			reqBody := []byte("")
			resBody := []byte("")
			if contentType := c.Get(fiber.HeaderContentType); strings.Contains(contentType, fiber.MIMEApplicationJSON) {
				reqBody = c.Body()
			}
			if resContentType := c.Response().Header.ContentType(); bytes.Contains(resContentType, []byte(fiber.MIMEApplicationJSON)) {
				resBody = c.Response().Body()
			}

			var reqHeader, resHeader []byte
			if value, err := sonic.Marshal(excludeHeader(c.GetReqHeaders())); err == nil {
				reqHeader = value
			}
			if value, err := sonic.Marshal(excludeHeader(c.GetRespHeaders())); err == nil {
				resHeader = value
			}

			log.Info("middleware-logging",
				zap.Int64("execute_time_in_ms", data.Stop.Sub(data.Start).Milliseconds()),
				zap.ByteString("url", c.Request().RequestURI()),
				zap.String("method", c.Method()),
				zap.Strings("ip", c.IPs()),
				zap.ByteString("req_header", reqHeader),
				zap.ByteString("req_body", reqBody),
				zap.ByteString("res_header", resHeader),
				zap.ByteString("res_body", resBody),
			)

			return nil
		},
	}))

	return m
}

func excludeHeader(header map[string][]string) map[string][]string {
	// TODO
	return header
}

func (m *FiberMiddleware) UseSwagger() *FiberMiddleware {
	m.app.Use(swagger.New(swagger.Config{
		BasePath: "/",
		Path:     "swagger",
		FilePath: "../../docs/swagger.json",
		Title:    "Swagger API",
		CacheAge: 3600,
	}))

	return m
}
