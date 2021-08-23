package influx

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/http"
	"time"
)

var client influxdb2.Client

const token = "DkLQz6aMVUaS0bOFvVq8U6Ea17IkLDpfPJK_ROHZwt6eTXU94LJigt16ATqsp0d2Qd0elQ2FakmUsSczeSWIXA=="
const bucket = "weibo"
const org = "zwz"

func init() {
	client = influxdb2.NewClient("http://hellozwz.com:8086", token)
	// always close client at the end
	defer client.Close()
}

func Write(measurement string, tags map[string]string, fields map[string]interface{}) (err error) {
	p := influxdb2.NewPoint(measurement, tags, fields, time.Now())
	writeApi := client.WriteAPI(org, bucket)
	writeApi.WritePoint(p)
	writeApi.Flush()
	writeApi.SetWriteFailedCallback(func(batch string, error http.Error, retryAttempts uint) bool {
		err = &error
		return false
	})
	return
}
