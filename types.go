package main

type Temps struct {
	Fahrenheit float64	`json:"fahrenheit"`
	Celsius    float64	`json:"celsius"`
}

type SensorData struct {
	ID	 	string	`json:"id"`
	Temperature	Temps	`json:"temperature"`
	Pressure	int	`json:"pressure"`
	Altitude	float64	`json:"altitude"`
}
