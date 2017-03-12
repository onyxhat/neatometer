package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/kardianos/osext"
	config "github.com/spf13/viper"
	"net/http"
)

func init() {
	exePath, err := osext.ExecutableFolder()
	if err != nil {
		exePath = ".\\"
	}

	config.AddConfigPath(exePath)
	config.SetConfigName("config")
	config.ReadInConfig()

	config.SetDefault("Binding", "0.0.0.0:8080")
}

func main() {
	log.Info("Listening at http://" + config.GetString("Binding"))

	mx := mux.NewRouter()
	mx.HandleFunc("/", IndexHandler)

	http.ListenAndServe(config.GetString("Binding"), mx)
}
