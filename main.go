package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chronark/go-queue/queue"
	f "github.com/fauna/faunadb-go/v4/faunadb"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
)

func ProduceHandler(c *fiber.Ctx) error {

	faunaToken := c.Get("Authorization")
	if faunaToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Authorization header is missing")
	}

	topic := c.Params("topic")
	if topic == "" {
		return fiber.NewError(fiber.StatusBadRequest, "topic is missing in url")
	}
	payload := c.Body()

	q, err := queue.NewQueue(f.NewFaunaClient(faunaToken))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Unable to connect to queue: %s", err.Error()))
	}

	message := queue.Message{
		Header: queue.Header{
			Id:        uuid.NewString(),
			Topic:     topic,
			CreatedAt: time.Now(),
		},
		Payload: payload,
	}

	err = q.Produce(message)
	if err != nil {
		return err
	}

	return c.JSON(message.Header)
}

func ConsumeHandler(c *fiber.Ctx) error {
	faunaToken := c.Get("Authorization")
	if faunaToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Authorization header is missing")
	}
	topic := c.Params("topic")
	if topic == "" {
		return fiber.NewError(fiber.StatusBadRequest, "topic is missing in url")
	}
	q, err := queue.NewQueue(f.NewFaunaClient(faunaToken))
	if err != nil {
		return err
	}
	message, err := q.Consume(topic)
	if err != nil {
		if strings.Contains(err.Error(), "Response error 404.") {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return err
	}

	return c.JSON(message)

}

func AcknowledgeHandler(c *fiber.Ctx) error {
	faunaToken := c.Get("Authorization")
	if faunaToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Authorization header is missing")
	}
	messageId := c.Params("messagId")
	if messageId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "messageId is missing in url")
	}
	q, err := queue.NewQueue(f.NewFaunaClient(faunaToken))
	if err != nil {
		return err
	}

	err = q.Acknowledge(messageId)
	if err != nil {
		if err != nil {
			if strings.Contains(err.Error(), "Response error 404.") {
				return c.SendStatus(fiber.StatusNotFound)
			}
			return err
		}

		return err
	}
	return c.SendStatus(fiber.StatusOK)

}

func main() {

	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())

	app.Post("/produce/:topic", ProduceHandler)
	app.Post("/acknowledge/:messagId", AcknowledgeHandler)
	app.Get("/consume/:topic", ConsumeHandler)

	err := app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		panic(err)
	}
}
