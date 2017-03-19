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

	//Disbale feature for non-configured setting
	if config.GetString("esURL") != "" {
		config.Set("EnableESForwarder", false)
	}

	setLogLevel(config.GetString("LogLevel"))
}

func main() {
	f := []bool{config.GetBool("EnableJSONServer"), config.GetBool("EnableESForwarder")}
	fn := countBool(f)

	runtime.GOMAXPROCS(fn)
	var wg sync.WaitGroup
	wg.Add(fn)

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
