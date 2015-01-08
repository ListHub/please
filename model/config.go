package model

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	PersistenceAddress string
}

func Config() Configuration {
	filename := "config.json"

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(file)
	config := Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Panic("error:", err)
	}

	return config
}
