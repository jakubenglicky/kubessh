package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	DefaultNamespace string `json:"defaultNamespace"`
}

func GetConfig() Config {

	userHome, _ := os.UserHomeDir()
	jsonFile, err := os.Open(userHome + "/.kubessh/config")

	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config Config

	json.Unmarshal(byteValue, &config)

	defer jsonFile.Close()

	return config
}
