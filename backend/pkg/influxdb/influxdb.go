package influxdb

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// InfluxDBClient определяет общий интерфейс для работы с InfluxDB.
type InfluxDBClient interface {
	influxdb2.Client
}

// NewInfluxDB создает новый клиент InfluxDB.
func NewInfluxDB(serverURL, authToken string) InfluxDBClient {
	return influxdb2.NewClient(serverURL, authToken)
}

// NewInfluxDBWithTimeout создает новый клиент InfluxDB с увеличенным timeout.
func NewInfluxDBWithTimeout(serverURL, authToken string, timeout time.Duration) InfluxDBClient {
	options := influxdb2.DefaultOptions()
	options.SetHTTPRequestTimeout(uint(timeout.Seconds()))

	return influxdb2.NewClientWithOptions(serverURL, authToken, options)
}
