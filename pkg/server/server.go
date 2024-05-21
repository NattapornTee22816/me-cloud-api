package server

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"os"
	"time"
)

type Server struct {
	*fiber.App
}

func NewServer() *Server {
	app := fiber.New(fiber.Config{
		AppName:           "Your-Cloud",
		ServerHeader:      "X-Server",
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		JSONDecoder:       sonic.Unmarshal,
		JSONEncoder:       sonic.Marshal,
		StreamRequestBody: true,
		StructValidator: &structValidator{
			validate: validator.New(
				validator.WithRequiredStructEnabled(),
			),
		},
		ErrorHandler: func(ctx fiber.Ctx, err error) error {
			return nil
		},
	})

	return &Server{app}
}

func (s *Server) Run() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	if err := s.Listen(fmt.Sprintf(":%s", port)); err != nil {
		panic(err)
	}
}
