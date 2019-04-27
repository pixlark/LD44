package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
)

var _ = fmt.Println

func loadPath(level *Level, index int, data map[string]interface{}) {
	orbIndex  := int(data["index"].(float64))
	flagIndex := int(data["flag"].(float64))
	
	stoppersData := data["stoppers"].([]interface{})
	stoppers := make([]Stopper, len(stoppersData))
	for i := 0; i < len(stoppers); i++ {
		position := int(stoppersData[i].(float64))
		stoppers[i] = newStopper(position)
	}
	
	swappersData := data["vertSwappers"].([]interface{})
	swappers := make([]VertSwapper, len(swappersData))
	for i := 0; i < len(swappers); i++ {
		position := int(swappersData[i].(float64))
		swappers[i] = newVertSwapper(position)
	}
	
	level.paths[index] = newPath(orbIndex, flagIndex, stoppers, swappers)
}

func loadPaths(level *Level, paths []interface{}) {
	for i, path := range paths {
		path, ok := path.(map[string]interface{})
		if !ok {
			fatal("Path must be object!")
		}
		loadPath(level, i, path)
	}
}

func loadLevel(path string) Level {
	rawData, err := ioutil.ReadFile(path)
	if err != nil {
		fatal("Could not open level file!")
	}

	var data map[string]interface{}
	err = json.Unmarshal(rawData, &data)
	if err != nil {
		fatal("Problem unmarshaling level file: " + err.Error())
	}
	//fmt.Println(data)
	
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
