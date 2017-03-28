package main

import (
	"bytes"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
	"github.com/kidoman/embd/sensor/bmp180"
	config "github.com/spf13/viper"
	"net/http"
	"time"
	"fmt"
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
	tempf := tempc*1.8 + 32
	altitude, _ := sensor.Altitude()
	pressure, _ := sensor.Pressure()

	return tempc, tempf, altitude, pressure
}

func getData() []byte {
	tempc, tempf, altitude, pressure := readSensor()

	data := SensorData{
		Timestamp: time.Now().Format(time.RFC3339),
		DeviceId:  config.GetString("DeviceID"),
		Temperature: Temps{
			Fahrenheit: tempf,
			Celsius:    tempc,
		},
		Pressure: pressure,
		Altitude: altitude,
	}

	myJson, err := json.Marshal(data)
	handleErr(err)

	return myJson
}

func newPostES() {
	baseUrl := config.GetString("ElasticSearchUrl")
	esIndex := config.GetString("ElasticSeachIndex")
	postUrl := fmt.Sprintf("%s/%s-%d-%d-%d/json/", baseUrl, esIndex, time.Now().Year(), time.Now().Month(), time.Now().Day())

	Block{
		Try: func() {
			resp, err := http.Post(postUrl, "application/json", bytes.NewBuffer(getData()))
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
