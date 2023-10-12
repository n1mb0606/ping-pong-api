package influxdbclient

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type Dto struct {
	from int64
	to   int64
}

type ResponseDto struct {
	Timestamp int64
	Code      int
	Up		  bool
	Ping	  float64
}

func GetNewInfluxClient(org string, bucket string, url string, token string) *InfluxClient {
	return &InfluxClient{
		token:  token,
		url:    url,
		org:    org,
		bucket: bucket,
		client: influxdb2.NewClient(url, token),
	}
}

type InfluxClient struct {
	token        string
	url          string
	org          string
	bucket       string
	client       influxdb2.Client
	query_table  *api.QueryTableResult
	query_result map[string]interface{}
}

// Write a line of protocol to the database
func (ic *InfluxClient) Write(measurement string, tags map[string]string, fields map[string]interface{}) error {
	writeAPI := ic.client.WriteAPIBlocking(ic.org, ic.bucket)
	point := write.NewPoint(measurement, tags, fields, time.Now())
	time.Sleep(1 * time.Second)

	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
		return err
	}
	return nil
}

// Send a query and print the results
func (ic *InfluxClient) Query(queryString string) error {
	queryAPI := ic.client.QueryAPI(ic.org)
	results, err := queryAPI.Query(context.Background(), queryString)

	if err != nil {
		return err
	}
	ic.query_table = results
	if err := results.Err(); err != nil {
		return err
	}
	return nil
}

func (ic *InfluxClient) NextResult() bool {
	if ic.query_table.Next() {
		ic.query_result = ic.query_table.Record().Values()
		return true
	}
	return false
}

func (ic *InfluxClient) GetResult() (map[string]interface{}, error) {
	if err := ic.query_table.Err(); err != nil {
		return nil, err
	} else {
		return ic.query_result, nil
	}
}
