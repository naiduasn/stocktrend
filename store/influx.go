package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/joho/godotenv"
)

var (
	host       = os.Getenv("INFLUXDB_HOST")
	token      = os.Getenv("INFLUXDB_TOKEN")
	bucketName = os.Getenv("INFLUXDB_BUCKET")
	client     *influxdb3.Client
)

// init sets up the InfluxDB client and its read and write APIs.
func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host = os.Getenv("INFLUXDB_HOST")
	token = os.Getenv("INFLUXDB_TOKEN")
	bucketName = os.Getenv("INFLUXDB_BUCKET")
	client, err = influxdb3.New(influxdb3.ClientConfig{
		Host:     host,
		Token:    token,
		Database: bucketName,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(client)

}

func InsertData(data []map[string]interface{}, symbolMap map[float64]string) error {
	start_time := time.Now()
	points := []*influxdb3.Point{}
	for _, stockData := range data {
		// Create a new point
		point := influxdb3.NewPointWithMeasurement("prices").
			SetTag("SecurityID", strconv.FormatInt(int64(stockData["security_id"].(float64)), 10)).
			SetTag("Symbol", symbolMap[stockData["security_id"].(float64)]).
			SetDoubleField("Price", stockData["last_price"].(float64)).
			SetDoubleField("CZG", stockData["change_percent"].(float64)).SetTimestamp(time.Now())

		points = append(points, point)
	}
	err := client.WritePoints(context.Background(), points...)
	if err != nil {
		log.Printf("Error writing points to InfluxDB: %v", err)
		return err
	}
	fmt.Println("time taken:", time.Since(start_time))
	return nil
}
