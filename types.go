package main

type Temps struct {
	Fahrenheit float64	`json:"fahrenheit"`
	Celsius    float64	`json:"celsius"`
}

type SensorData struct {
	Timestamp	string	`json:"@timestamp"`
	DeviceId 	string	`json:"deviceid"`
	Temperature	Temps	`json:"temperature"`
	Pressure	int	`json:"pressure"`
	Altitude	float64	`json:"altitude"`
}
