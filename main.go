package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Get("/api/v1/log", func(c *fiber.Ctx) error {
		querystring := "from(bucket: \"weblog\")\n" +
			"\t|> range(start: -2s)\n" +
			"\t|> filter(fn: (r) => r._measurement == \"wlog\")"+
			"\t|> limit(n:1)\n" +
			"\t|> pivot(rowKey: [\"_time\"], columnKey: [\"_field\"], valueColumn: \"_value\")"

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
				Code:      int(result["status"].(int64)),
				Up: 	   bool(result["up"].(bool)),
				Ping: 	   float64(result["ping"].(float64)),
			}
		}
		responseByte, _ := json.Marshal(resp)
		return c.SendString(string(responseByte))
	})
	app.Get("/api/v1/logs/from/:from/to/:to", func(c *fiber.Ctx) error {
		querystring := "from(bucket: \"weblog\")\n" +
			"\t|> range(start:" + c.Params("from") + ", stop:"+ c.Params("to") + ")\n" +
			"\t|> filter(fn: (r) => r._measurement == \"wlog\")"+
			"\t|> pivot(rowKey: [\"_time\"], columnKey: [\"_field\"], valueColumn: \"_value\")"

		if err := influxClient.Query(querystring); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		resps := make([]influxdbclient.ResponseDto, 0)
		for influxClient.NextResult() {
			result, err := influxClient.GetResult()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			var resp = influxdbclient.ResponseDto{
				Timestamp: (result["_time"].(time.Time)).Unix(),
				Code:      int(result["status"].(int64)),
				Up: 	   bool(result["up"].(bool)),
				Ping: 	   float64(result["ping"].(float64)),
			}
			resps = append(resps, resp)
		}
		responseByte, _ := json.Marshal(resps)
		responseByte[0] = '{'
		responseByte[len(responseByte) - 1] = '}'
		return c.SendString(string(responseByte))
	})

	app.Listen(":3000")
}
