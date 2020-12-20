package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var Cfg Config

func Init(filePath string) error {
	cfgFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("[Error][Init] Config file not found: %s", filePath)
	}
	defer cfgFile.Close()

	decoder := json.NewDecoder(cfgFile)
	Cfg = Config{}
	err = decoder.Decode(&Cfg)
	if err != nil {
		return fmt.Errorf("[Error][Init] Unknown field in config file: %s", err.Error())
	}
	return nil
}

func Debug() string {
	res, err := json.MarshalIndent(Cfg, "", "\t")
	if err != nil {
		return "no config"
	}
	return string(res)
}
