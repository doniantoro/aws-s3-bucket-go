package main

import (
	configApp "aws-s3-bucket/config"
	uploadHttp "aws-s3-bucket/domain/upload/delivery/http"
	uploadUsecase "aws-s3-bucket/domain/upload/usecase"
	"aws-s3-bucket/shared/constant"
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/google/uuid"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("REGION_NAME")),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	app := fiber.New(
		fiber.Config{
			AppName: "aws-bucket",
		},
	)
	// Middleware to set request ID and CORS headers
	app.Use(func(c *fiber.Ctx) error {
		requestId := c.Get(constant.HEADER_REQUEST_ID)
		if requestId == "" {
			requestId = uuid.New().String()
			c.Set(constant.HEADER_REQUEST_ID, requestId)
			c.Locals(constant.HEADER_REQUEST_ID, requestId)

		}
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin, Cache-Control, Pragma, X-Request-ID")
		c.Set("Access-Control-Allow-Credentials", "true")
		if c.Method() == fiber.MethodOptions {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	})

	limitThreshold, _ := strconv.Atoi(os.Getenv("LIMITER_THRESHOLD"))
	limitExpired, _ := time.ParseDuration(os.Getenv("LIMITER_EXPIRED"))

	app.Use(limiter.New(limiter.Config{
		Max:               limitThreshold,
		Expiration:        limitExpired,
		LimiterMiddleware: limiter.FixedWindow{},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		Next: func(c *fiber.Ctx) bool {
			return strings.HasPrefix(c.Path(), "/api/v1/docs")
		},
	}))

	app.Use(swagger.New(swagger.Config{
		BasePath: "/api/v1/",
		FilePath: "./docs/swagger.json",
		Path:     "docs",
		Title:    "Swagger API Docs",
		CacheAge: 1,
	}))

	v1 := app.Group("/api/v1/")
	validator := configApp.NewValidator()

	// Initialize the usecase
	multiUsecase := uploadUsecase.NewUsecase(s3Client)

	// Initialize the upload HTTP handler
	uploadHttp.NewHandler(v1, multiUsecase, validator)

	log.Fatal(app.Listen(":" + os.Getenv("APP_PORT")))
}
