package main

import (
	"github.com/chronark/go-queue/queue"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"time"
)



func ProduceHandler(c *fiber.Ctx, redisClient *redis.Client) error {
	message := c.Body()
	topic := c.Query("topic")

	q, err := queue.NewQueue("tenant", redisClient)
	if err != nil {
		return err
	}

	err = q.Produce(topic, queue.SignedMessage{
		Message: queue.Message{
			Header: queue.Header{
				Id:        "id",
				CreatedAt: time.Now(),
			},
			Body: message,
		},
	})
	if err != nil {
		return err
	}

	return c.Send(c.Body())
}

func main() {
	redisOptions, err := redis.ParseURL("rediss://:da3d3c3719e34166a52ea4c92a29b441@eu1-funny-firefly-32822.upstash.io:32822")
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(redisOptions)

	if err != nil {
		panic(err)
	}
	// if false {

	// 	err = q.Produce("topic", queue.SignedMessage{
	// 		Message: queue.Message{
	// 			Header: queue.Header{
	// 				Id:        "id",
	// 				CreatedAt: time.Now(),
	// 			},
	// 			Body: "xxxxxx World",
	// 		},
	// 	})
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// val, err := q.Consume("topic")
	// if err != nil {
	// 	panic(err)
	// }

	// m, err := queue.Deserialize(val)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%+v\n", m)

	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())

	app.Post("/produce/:topic", func(c *fiber.Ctx) error { return ProduceHandler(c, redisClient) })

	app.Listen(":9000")

}
