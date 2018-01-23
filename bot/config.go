package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Config ...
type Config struct {
	Token           string
	WolframAlphaKey string
}

func createBlankConfig() {
	saveConfig(Config{Token: "insert token here", WolframAlphaKey: "insert WA key here"})
}

func loadConfig() (*Config, error) {

	dataExists, _ := exists("data")
	fileExists, _ := exists("data/config.json")

	if !fileExists || !dataExists {
		os.Mkdir("data", os.ModePerm)
		createBlankConfig()
		return nil, nil
	}
	b, rerr := ioutil.ReadFile("data/config.json")

	if rerr != nil {
		return nil, rerr
	}

	var sett Config
	err := json.Unmarshal(b, &sett)

	if err != nil {
		return nil, err
	}

	return &sett, nil
}

func saveConfig(s Config) {
	if ex, _ := exists("data"); !ex {
		os.Mkdir("data", os.ModePerm)
	}

	b, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		panic(err)
	}

	if ex, _ := exists("data/config.json"); !ex {
		f, _ := os.Create("data/config.json")
		f.Close()
	}

	ioutil.WriteFile("data/config.json", b, os.ModePerm)
}
