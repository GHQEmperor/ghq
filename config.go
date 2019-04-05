package ghq

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

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
	r.Config = make(map[string]string)
	return json.Unmarshal(configBytes,&r.Config)
}

func (r *Router) GetConfig(configName string) (config string,ok bool) {
	config,ok = r.Config[configName]
	return
}
