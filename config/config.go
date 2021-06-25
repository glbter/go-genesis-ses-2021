package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	DbName          string `json:"db_name"`
	NbuLink         string `json:"nbu_link"`
	NbuParam        string `json:"nbu_param"`
	BlockchainLink  string `json:"blockchain_link"`
	BlockchainParam string `json:"blockchain_param"`
}

var config = Config{}

func InitConfig() (Config, error) {
	jsonData, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalln("error reading config")
		return config, err
	}

	var newConfig Config
	err = json.Unmarshal(jsonData, &newConfig)

	if err != nil {
		log.Fatalln("error parsing config")
		return config, err
	}
	config = newConfig

	return config, err
}

func GetConfig() Config {
	config, _ := InitConfig()
	return config
}
