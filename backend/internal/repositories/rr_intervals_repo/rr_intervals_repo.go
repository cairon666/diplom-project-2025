package rr_intervals_repo

import (
	"github.com/cairon666/vkr-backend/internal/config"
	"github.com/cairon666/vkr-backend/pkg/influxdb"
)

// RRIntervalsRepo реализует репозиторий для работы с R-R интервалами в InfluxDB.
type RRIntervalsRepo struct {
	influxClient influxdb.InfluxDBClient
	org          string
	bucket       string
}

// NewRRIntervalsRepo создает новый экземпляр репозитория.
func NewRRIntervalsRepo(influxClient influxdb.InfluxDBClient, conf *config.Config) *RRIntervalsRepo {
	return &RRIntervalsRepo{
		influxClient: influxClient,
		org:          conf.InfluxDB.Org,
		bucket:       conf.InfluxDB.Bucket,
	}
}
