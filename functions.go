package main

import (
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
	"github.com/kidoman/embd/sensor/bmp180"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/shirou/gopsutil/host"
)

func handleErr(err error, errMsg string) {
	if err != nil {
		log.Error(errMsg)
		//os.Exit(-1)
	}
}

func readSensor() (float64, float64, float64, int) {
	bus := embd.NewI2CBus(1)
	sensor := bmp180.New(bus)

	tempc, err := sensor.Temperature()
	handleErr(err, "Unable to read temperature")

	tempf := tempc*1.8+32

	altitude, err := sensor.Altitude()
	handleErr(err, "Unable to read altitude")

	pressure, _ := sensor.Pressure()
	handleErr(err, "Unable to read pressure")

	return tempc, tempf, altitude, pressure
}

func getData() []byte {
	hostInfo, err := host.Info()
	handleErr(err, "Unable to detect host info")

	tempc, tempf, altitude, pressure := readSensor()

	data := SensorData{
		ID: hostInfo.HostID,
		Temperature: Temps{
			Fahrenheit: tempf,
			Celsius: tempc,
		},
		Pressure: pressure,
		Altitude: altitude,
	}

	myJson,err := json.Marshal(data)
	handleErr(err, "Unable to marshal JSON")

	return myJson
}
