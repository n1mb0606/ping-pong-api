package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"time"
	"github.com/n1mb0606/ping-pong-api/influxdbclient"
)


func main() {
	app := fiber.New()
	org := os.Getenv("PING_PONG_ORG")
	bucket := os.Getenv("PING_PONG_BUCKET")
	token := os.Getenv("PING_PONG_TOKEN")
	url := os.Getenv("PING_PONG_URL")
	influxClient := influxdbclient.GetNewInfluxClient(org, bucket, url, token)

	app.Get("/api/v1/log", func(c *fiber.Ctx) error {
		querystring := "from(bucket: \"weblog\")\n" +
			"\t|> range(start: -10m)\n" +
			"\t|> sort(columns: [\"_time\"], desc: true)" +
			"\t|> limit(n:1)" +
			"\t|> filter(fn: (r) => r._measurement == \"wlog\")"

		if err := influxClient.Query(querystring); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		var resp = influxdbclient.ResponseDto{}
		for influxClient.NextResult() {
			result, err := influxClient.GetResult()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			resp = influxdbclient.ResponseDto{
				Timestamp: (result["_time"].(time.Time)).Unix(),
				Code:      int(result["_value"].(int64)),
			}
		}
		responseByte, _ := json.Marshal(resp)
		return c.SendString(string(responseByte))
	})
	app.Get("/api/v1/logs", func(c *fiber.Ctx) error {
		return c.SendString("TEST")
	})

	app.Listen(":3000")
}
