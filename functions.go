package main

import (
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
	"github.com/kidoman/embd/sensor/bmp180"
	"encoding/json"
)

type Temps struct {
	Fahrenheit float64
	Celsius    float64
}

type SensorData struct {
	Temp	 Temps
	Pressure int
	Altitude float64
}

func getData() []byte {
	bus := embd.NewI2CBus(1)
	sensor := bmp180.New(bus)

	tempc, _ := sensor.Temperature()
	tempf := tempc*(9/5)+32
	altitude, _ := sensor.Altitude()
	pressure, _ := sensor.Pressure()

	data := SensorData{
		Temp: Temps{
			Fahrenheit: tempf,
			Celsius: tempc,
		},
		Pressure: pressure,
		Altitude: altitude,
	}

	myJson,_ := json.Marshal(data)

	return myJson
}
