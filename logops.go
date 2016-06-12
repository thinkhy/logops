package logops

import (
	"fmt"
	client "github.com/influxdata/influxdb/client/v2"
	"log"
	"time"
)

const (
	// PrecisionDefault represents the default precision used for the InfluxDB points.
	PrecisionDefault = "ns"
	// DatabaseDefault is the default database that we will write to, if not specified otherwise in the Config for the hook.
	DatabaseDefault = "logops"
	// DefaultMeasurementValue is the default measurement that we will assign to each point, unless there is a field called "measurement".
	DefaultMeasurementValue = "logops"
	// BatchIntervalDefault represents the number of seconds that we wait for a batch to fill up.
	// After that we flush it to InfluxDB whatsoever.
	BatchIntervalDefault = 5
	// BatchSizeDefault represents the maximum size of a batch.
	BatchSizeDefault = 1000
	// MaxRetryCountDefault represents the maximum number of retrying to connect DB.
	MaxRetryCountDefault = 10
)

// Config data for InfluxDB
type Config struct {
	Precision        string
	Address          string // ip:port
	Username         string
	Password         string
	Database         string
	MeasurementValue string
	// Tags that we will extract from the log fields and set
	// them as Influx point tags.
	Tags           []string
	BatchInterval  int // seconds
	BatchSize      int
	MaxRetryCount  int
	UseUDP         bool
	retryCount     int
	chTearDown     chan bool // just for testing, don't use it in production env
	chTearDownDone chan bool // just for testing, don't use it in production env
}

// Activity represents the operation that specific module performs
type Activity struct {
	module string
	who    string
	how    string
	what   string
}

// Hook represents the hook to InfluxDB
type Hook struct {
	config       Config
	chActivities chan *Activity
}

func (config *Config) setDefaults() {
	if config.Precision == "" {
		config.Precision = PrecisionDefault
	}
	if config.Database == "" {
		config.Database = DatabaseDefault
	}
	if config.MeasurementValue == "" {
		config.MeasurementValue = DefaultMeasurementValue
	}
	if config.BatchInterval <= 0 {
		config.BatchInterval = BatchIntervalDefault
	}
	if config.BatchSize <= 0 {
		config.BatchSize = BatchSizeDefault
	}

	if config.MaxRetryCount <= 0 {
		config.MaxRetryCount = MaxRetryCountDefault
	}
}

// NewHook generate a new InfluxDB hook based on the given configuration
func NewHook(config *Config) (*Hook, error) {
	if config == nil {
		return nil, fmt.Errorf("Influxus configuration passed to InfluxDB is nil.")
	}
	config.setDefaults()
	hook := &Hook{
		config: *config,
	}

	// Make a buffered channel so that senders will not block.
	hook.chActivities = make(chan *Activity, config.BatchSize)
	hook.config.chTearDown = make(chan bool)
	hook.config.chTearDownDone = make(chan bool)
	go hook.startBatchHandler()
	return hook, nil
}

func (hook *Hook) startBatchHandler() {
	done := false
	// Make client
	for ; done || hook.config.retryCount < hook.config.MaxRetryCount; hook.config.retryCount++ {
		// wait for some seconds and retry
		time.Sleep(time.Duration(hook.config.retryCount) * time.Second)

		var err error
		var c client.Client
		if hook.config.UseUDP {
			c, err = client.NewUDPClient(client.UDPConfig{Addr: hook.config.Address})
		} else {
			c, err = client.NewHTTPClient(client.HTTPConfig{
				Addr:     hook.config.Address,
				Username: hook.config.Username,
				Password: hook.config.Password,
			})
		}
		if err != nil {
			log.Printf("[logops] Make client #%d, Error: %v", hook.config.retryCount, err)
			continue
		}
		defer c.Close()

		if hook.config.UseUDP == false {
			q := client.NewQuery(fmt.Sprintf("CREATE DATABASE %s", hook.config.Database), "", "")
			if response, err := c.Query(q); err != nil {
				fmt.Println("[logops] Failed to create db ", hook.config.Database, response.Error())
			}
		}

		for true {
			select {
			case a := <-hook.chActivities:
				// Create a new point batch
				bp, err := client.NewBatchPoints(client.BatchPointsConfig{
					Database:  hook.config.Database,
					Precision: hook.config.Precision,
				})
				if err != nil {
					log.Fatalln("[logops] NewBatchPoints Error: ", err)
					break
				}

				// Create a point and add to batch
				tags := map[string]string{"module": a.module}
				fields := map[string]interface{}{
					"value": 1,
					"who":   a.who,
					"how":   a.how,
					"what":  a.what,
				}
				pt, err := client.NewPoint(hook.config.MeasurementValue, tags, fields)
				if err != nil {
					log.Println("[logops] NewPoint Error: ", err)
					break
				}
				// fmt.Println("%s\n", pt)
				bp.AddPoint(pt)

				// Write the batch
				err = c.Write(bp)
				if err != nil {
					log.Println("[logops] Write Error: ", err)
					break
				}
			// For testing
			case <-hook.config.chTearDown:
				if hook.config.UseUDP == false {
					q := client.NewQuery(fmt.Sprintf("DROP DATABASE %s", hook.config.Database), "", "")
					if response, err := c.Query(q); err != nil {
						fmt.Println("Failed to create db ", hook.config.Database, response.Error())
					}
				}
				done = true
				hook.config.chTearDownDone <- done
				break
			} // select
		} // for true
	} // for retry
}

func (hook *Hook) Write(module, who, how, what string) {
	a := &Activity{
		module: module,
		who:    who,
		how:    how,
		what:   what,
	}
	hook.chActivities <- a
	return
}

func (hook *Hook) tearDown() {
	hook.config.chTearDown <- true
	<-hook.config.chTearDownDone
}
