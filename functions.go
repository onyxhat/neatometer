package main

import (
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
	"github.com/kidoman/embd/sensor/bmp180"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/shirou/gopsutil/host"
	"bytes"
	"net/http"
	"time"
)

func setLogLevel(level string) {
	switch level {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	}
}

func handleErr(err error) {
	if err != nil {
		log.Error(err)
	}
}

/*
Try/Catch/Finally Implementation
Block{
	Try: func() {

	},
	Catch: func(e Exception) {

	},
	Finally: func() {

	},
}.Do()
 */

func Throw(up Exception) {
	panic(up)
}

func (tcf Block) Do() {
	if tcf.Finally != nil {

		defer tcf.Finally()
	}
	if tcf.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				tcf.Catch(r)
			}
		}()
	}
	tcf.Try()
}

func readSensor() (float64, float64, float64, int) {
	bus := embd.NewI2CBus(1)
	sensor := bmp180.New(bus)

	tempc, _ := sensor.Temperature()
	tempf := tempc*1.8+32
	altitude, _ := sensor.Altitude()
	pressure, _ := sensor.Pressure()

	return tempc, tempf, altitude, pressure
}

func getData() []byte {
	hostInfo, err := host.Info()
	handleErr(err)

	tempc, tempf, altitude, pressure := readSensor()

	data := SensorData{
		Timestamp: time.Now().Format(time.RFC3339),
		DeviceId: hostInfo.HostID,
		Temperature: Temps{
			Fahrenheit: tempf,
			Celsius: tempc,
		},
		Pressure: pressure,
		Altitude: altitude,
	}

	myJson,err := json.Marshal(data)
	handleErr(err)

	return myJson
}

func newPostES(url string) {
	Block{
		Try: func() {
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(getData()))
			handleErr(err)
			defer resp.Body.Close()
		},
		Catch: func(e Exception) {
			log.Error("POST to ES failed")
		},
		Finally: func() {
			log.Debug("POST to ES succeeded")
		},
	}.Do()
}

func isTrue(val bool) int {
	if val == true {
		return 1
	} else {
		return 0
	}
}

func countBool(myBool []bool) int {
	var count int
	for _, v := range myBool {
		count += isTrue(v)
	}
	return count
}