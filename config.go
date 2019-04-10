package ghq

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var Config map[string]string

// loading config.json in memory.
func (r *Router) LoadConfig() (err error) {
	configFile, err := os.Open("config.json")
	if err != nil {
		return
	}
	defer configFile.Close()
	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		return
	}
	Config = make(map[string]string)
	err = json.Unmarshal(configBytes, &Config)
	if err != nil {
		return
	}
	//r.Config = Config
	return
}

//func (r *Router) GetConfig(configName string) (config string, ok bool) {
func GetConfig(configName string) (config string, ok bool) {
	//config,ok = r.Config[configName]
	config, ok = Config[configName]
	return
}
