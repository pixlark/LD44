package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
)

var _ = fmt.Println

func loadPaths(level *Level, data []interface{}) {
	
}

func loadLevel(path string) Level {
	rawData, err := ioutil.ReadFile("level0.json")
	if err != nil {
		fatal("Could not open level file!")
	}

	var data map[string]interface{}
	err = json.Unmarshal(rawData, &data)
	if err != nil {
		fatal("Problem unmarshaling level file: " + err.Error())
	}
	fmt.Println(data)
	
	var level Level

	_, ok := data["width"]
	if !ok {
		fatal("No \"width\" field on level file")
	}
	level.width = int(data["width"].(float64))

	_, ok = data["height"]
	if !ok {
		fatal("No \"height\" field on level file")
	}
	level.paths = make([]Path, int(data["height"].(float64)))

	_, ok = data["paths"]
	if !ok {
		fatal("No \"paths\" field on level file")
	}

	loadPaths(&level, data["paths"].([]interface{}))
	
	return level
}
