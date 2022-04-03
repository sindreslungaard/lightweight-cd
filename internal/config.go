package internal

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

const path = "./data"
const fileName = "config.json"

var mu sync.Mutex

type Config struct {
	Version     string                `json:"version"`
	Deployments map[string]Deployment `json:"deployments"`
}

func FilePath() string {
	return path + "/" + fileName
}

func save(filePath string, d Config) {
	bytes, err := json.Marshal(d)

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filePath, bytes, os.ModePerm)

	if err != nil {
		panic(err)
	}
}

func ReadConfig() Config {
	mu.Lock()
	defer mu.Unlock()

	d := getOrCreateConfig()
	return d
}

func UpdateConfig(f func(Config) Config) {
	mu.Lock()
	defer mu.Unlock()

	d := getOrCreateConfig()

	data := f(d)

	save(FilePath(), data)
}

func getOrCreateConfig() Config {

	data, err := os.ReadFile(FilePath())

	if err != nil {
		println("No configuration found, setting up for first time use")

		config := Config{
			Deployments: make(map[string]Deployment),
		}

		data, err = json.MarshalIndent(config, "", "	")

		if err != nil {
			panic(err)
		}

		err = os.MkdirAll(path, os.ModeDir)

		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(FilePath(), data, os.ModePerm)

		if err != nil {
			panic(err)
		}

	}

	var d Config

	err = json.Unmarshal(data, &d)

	if err != nil {
		panic(err)
	}

	return d

}
