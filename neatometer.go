package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/kardianos/osext"
	config "github.com/spf13/viper"
	"net/http"
	"runtime"
	"sync"
	"time"
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
	config.SetDefault("EnableJSONServer", true)
	config.SetDefault("EnableESForwarder", true)
	config.SetDefault("PollInterval", 10)
	config.SetDefault("LogLevel", "INFO")

	setLogLevel(config.GetString("LogLevel"))
}

func main() {
	runtime.GOMAXPROCS(2)
	var wg sync.WaitGroup
	wg.Add(2)

	//Spawn http handler
	if config.GetBool("EnableJSONServer") {
		go func() {
			defer wg.Done()

			log.Info("Listening at http://" + config.GetString("Binding"))

			mx := mux.NewRouter()
			mx.HandleFunc("/", IndexHandler)

			http.ListenAndServe(config.GetString("Binding"), mx)
		}()
	} else { wg.Done() }

	//Spawn es forwarder
	if config.GetBool("EnableESForwarder") {
		go func() {
			defer wg.Done()

			for {
				duration := config.GetDuration("PollInterval") * time.Second
				time.Sleep(duration)

				newPostES(config.GetString("esURL"))
			}
		}()
	} else { wg.Done() }

	wg.Wait()
	log.Info("Terminating program")
}
