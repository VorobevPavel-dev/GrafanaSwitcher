package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/tidwall/sjson"
)

//todo получить путь к тэгам

var (
	routes []string
)

func GetRoutes(json map[string]interface{}, currentPath, tagToFind string) {
	//Проходим по всем тэгам
	for currentKey, currentValue := range json {
		//Тип нашего текущего тэга
		valueType := fmt.Sprintf("%T", currentValue)
		switch valueType {
		//Случаи, когда значение простое
		case "float64":
			if currentKey == tagToFind {
				routes = append(routes, currentPath+`\\`+currentKey)
			}
		case "string":
			if currentKey == tagToFind {
				routes = append(routes, currentPath+`\\`+currentKey)
			}
		case "bool":
			if currentKey == tagToFind {
				routes = append(routes, currentPath+`\\`+currentKey)
			}
		case "<nil>":
			if currentKey == tagToFind {
				routes = append(routes, currentPath+`\\`+currentKey)
			}
			//Случаи, когда значение тэга массив
		case "map[string]interface {}":
			GetRoutes(currentValue.(map[string]interface{}), currentPath+`\\`+currentKey, tagToFind)
		case "[]interface {}":
			temp := currentValue.([]interface{})
			for i := range currentValue.([]interface{}) {
				val, ok := temp[i].(map[string]interface{})
				if ok {
					GetRoutes(val, currentPath+`\\`+currentKey+strconv.Itoa(i), tagToFind)
				} else {
					if currentKey == tagToFind {
						routes = append(routes, currentPath+"\\"+currentKey)
					}
				}
			}

		}
	}
}

func stringToMap(jsonString string) map[string]interface{} {
	var result = make(map[string]interface{})
	json.Unmarshal([]byte(jsonString), &result)
	return result
}

func mapToFile(filename string, jsonMap map[string]interface{}) {
	data, _ := json.MarshalIndent(jsonMap, "", "\t")
	ioutil.WriteFile(filename, data, 0777)
}

func changeTag(jsonMap map[string]interface{}, tagName string, newValue interface{}) map[string]interface{} {
	workMap := func(shell map[string]interface{}) map[string]interface{} {
		return shell["dashboard"].(map[string]interface{})
	}(jsonMap)
	var err error
	byteJSON, _ := json.Marshal(jsonMap)
	ioutil.WriteFile("output1.json", byteJSON, 0777)
	stringJSON := string(byteJSON)
	GetRoutes(workMap, "dashboard", tagName)
	for _, i := range routes {
		fmt.Println(i)
		stringJSON, err = sjson.Set(stringJSON, i, newValue)
		if err != nil {
			fmt.Println(err)
		}
	}

	jsonMap = stringToMap(stringJSON)
	return jsonMap
}

func main() {
	var mainMap = make(map[string]interface{})
	data, _ := ioutil.ReadFile("model.json")
	json.Unmarshal(data, &mainMap)
	mainMap = changeTag(mainMap, "target", "shit")
	mapToFile("output.json", mainMap)
}
