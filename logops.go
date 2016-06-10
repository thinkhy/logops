package logops

import (
	"github.com/influxdata/influxdb/client/v2"
	"log"
	"time"
)

const (
	MyDB     = "testhub"
	username = ""
	password = ""
)

func Fire() {
	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://45.55.21.6:8086",
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  MyDB,
		Precision: "s",
	})

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	// Create a point and add to batch
	tags := map[string]string{"cpu": "cpu-total"}
	fields := map[string]interface{}{
		"idle":   10.1,
		"system": 53.3,
		"user":   46.6,
	}
	pt, err := client.NewPoint("cpu_usage", tags, fields, time.Now())

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	bp.AddPoint(pt)

	// Write the batch
	err = c.Write(bp)
	if err != nil {
		log.Fatalln("Write Error: ", err)
	}
}
