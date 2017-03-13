package main

type Temps struct {
	Fahrenheit float64	`json:"fahrenheit"`
	Celsius    float64	`json:"celsius"`
}

type SensorData struct {
	ID	 string		`json:"id"`
	Temp	 Temps		`json:"temp"`
	Pressure int		`json:"pressure"`
	Altitude float64	`json:"altitude"`
}
