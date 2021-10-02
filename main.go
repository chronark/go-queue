package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/chronark/go-queue/queue"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
)

func ProduceHandler(c *fiber.Ctx, redisClient *redis.Client) error {
	tenant := c.Params("tenant")
	if tenant == "" {
		return fiber.NewError(fiber.StatusBadRequest,"tenant is missing in url")
	}
	topic := c.Params("topic")
	if tenant == "" {
		return fiber.NewError(fiber.StatusBadRequest,"topic is missing in url")

	}
	message := c.Body()

	q, err := queue.NewQueue("tenant", redisClient)
	if err != nil {
		return err
	}

	messageId := uuid.NewString()

	signedMessage := queue.SignedMessage{
		Message: queue.Message{
			Header: queue.Header{
				Id:        messageId,
				CreatedAt: time.Now(),
			},
			Body: message,
		},
	}

	err = q.Produce(topic, signedMessage)
	if err != nil {
		return err
	}

	return c.JSON(signedMessage)
}

func ConsumeHandler(c *fiber.Ctx, redisClient *redis.Client) error {
	tenant := c.Params("tenant")
	if tenant == "" {
		return fiber.NewError(fiber.StatusBadRequest,"tenant is missing in url")
	}
	topic := c.Params("topic")
	if tenant == "" {
		return fiber.NewError(fiber.StatusBadRequest,"topic is missing in url")
	}

	q, err := queue.NewQueue("tenant", redisClient)
	if err != nil {
		return err
	}

	message, err := q.Consume(topic)
	if err != nil {
		if err == redis.Nil {
			return c.SendStatus(http.StatusNoContent)
		}

		return err
	}
	return c.SendString(message)
}

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})

	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())

	app.Post("/:tenant/produce/:topic", func(c *fiber.Ctx) error { return ProduceHandler(c, redisClient) })
	app.Post("/:tenant/acknowledge/:messagId", func(c *fiber.Ctx) error { return ProduceHandler(c, redisClient) })
	app.Get("/:tenant/consume/:topic", func(c *fiber.Ctx) error { return ConsumeHandler(c, redisClient) })

	err := app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		panic(err)
	}
}
