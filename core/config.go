package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//ConfigData Configuration json structure
type ConfigData struct {
	SourceDir           string   `json:"sourceDir"`
	TargetSuffixes      []string `json:"targetExts"`
	Ccopts              string   `json:"ccopts"`
	LinkOpts            string   `json:"linkerOpts"`
	Compiler            string   `json:"compiler"`
	TatgetBinaryName    string   `json:"targetBinaryName"`
	CustomLinkerCommand string   `json:"customLinkerCmd"`
}

func handleConfigError(err *error, message string) {
	if *err != nil {
		fmt.Printf("[Error] %s\n", message)
		os.Exit(0)
	}
}

//LoadConfig Loads config from config json file
func LoadConfig(configFile string) *ConfigData {
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		handleConfigError(&err, "File "+configFile+" not found")
	}

	configOpts := ConfigData{}

	jsonData, err := ioutil.ReadFile(configFile)
	handleConfigError(&err, "Unable to open file "+configFile)

	err = json.Unmarshal([]byte(jsonData), &configOpts)
	handleConfigError(&err, "Invalid config file "+configFile)

	fmt.Printf("Loaded configuration : %v\n", configOpts)

	return &configOpts
}
