package influx

import (
	"github.com/akazwz/go-crawler/global"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/http"
	"time"
)

var client influxdb2.Client

func init() {
	client = influxdb2.NewClient(global.CFG.URL, global.CFG.Token)
	// always close client at the end
	defer client.Close()
}

func Write(measurement string, tags map[string]string, fields map[string]interface{}) (err error) {
	p := influxdb2.NewPoint(measurement, tags, fields, time.Now())
	writeApi := client.WriteAPI(global.CFG.Org, global.CFG.Bucket)
	writeApi.WritePoint(p)
	writeApi.Flush()
	writeApi.SetWriteFailedCallback(func(batch string, error http.Error, retryAttempts uint) bool {
		err = &error
		return false
	})
	return
}
