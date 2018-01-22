package main

import "github.com/vmihailenco/msgpack"
import "os"
import "io/ioutil"

func defaultSettings() settings {
	return settings{
		catalyst: ";",
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func load() (*settings, error) {

	dataExists, _ := exists("data")
	fileExists, _ := exists("data/settings.json")

	if !dataExists || !fileExists {
		set := defaultSettings()
		save(set)
		return &set, nil
	}

	b, rerr := ioutil.ReadFile("data/settings.pack")

	if rerr != nil {
		return nil, rerr
	}

	var sett settings
	err := msgpack.Unmarshal(b, &sett)

	if err != nil {
		return nil, err
	}

	return &sett, nil
}

func save(s settings) {
	if ex, _ := exists("data"); !ex {
		os.Mkdir("data", os.ModePerm)
	}

	b, err := msgpack.Marshal(s)
	if err != nil {
		panic(err)
	}

	if ex, _ := exists("data/settings.pack"); !ex {
		f, _ := os.Create("data/settings.pack")
		f.Close()
	}

	ioutil.WriteFile("data/settings.pack", b, os.ModePerm)
}
