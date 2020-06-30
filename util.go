package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/tidwall/sjson"
)

//removeMetaTag удаляет тэг meta из карты.
//При получении JSONModel в него записывается глобальная информация.
//При попытке отправить данные без парсинга возникнет ошибка.
func removeMetaTag(temp map[string]interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	for key, value := range temp {
		if key != "meta" {
			ret[key] = value
		}
	}
	return ret
}

//repairJSON позволяет изменить "сломанные" символы при разборе JSON файла
func repairJSON(data []byte) []byte {
	data = bytes.Replace(data, []byte("\\u003c"), []byte("<"), -1)
	data = bytes.Replace(data, []byte("\\u003e"), []byte(">"), -1)
	data = bytes.Replace(data, []byte("\\u0026"), []byte("&"), -1)
	return data
}

//getRoutes позволяет найти путь то изменяемых тэгов.
//Результат записывает в глобальную переменную routes[]string
func getRoutes(json map[string]interface{}, currentPath, tagToFind string) {
	//Проходим по всем тэгам
	for currentKey, currentValue := range json {
		//Тип нашего текущего тэга
		valueType := fmt.Sprintf("%T", currentValue)
		switch valueType {
		//Случаи, когда значение простое
		case "float64":
			if currentKey == tagToFind {
				routes = append(routes, currentPath+"."+currentKey)
			}
		case "string":
			if currentKey == tagToFind {
				routes = append(routes, currentPath+"."+currentKey)
			}
		case "bool":
			if currentKey == tagToFind {
				routes = append(routes, currentPath+"."+currentKey)
			}
		case "<nil>":
			if currentKey == tagToFind {
				routes = append(routes, currentPath+"."+currentKey)
			}
			//Случаи, когда значение тэга массив
		case "map[string]interface {}":
			getRoutes(currentValue.(map[string]interface{}), currentPath+"."+currentKey, tagToFind)
		case "[]interface {}":
			temp := currentValue.([]interface{})
			for i := range currentValue.([]interface{}) {
				val, ok := temp[i].(map[string]interface{})
				if ok {
					getRoutes(val, currentPath+"."+currentKey+"."+strconv.Itoa(i), tagToFind)
				} else {
					if currentKey == tagToFind {
						routes = append(routes, currentPath+"."+currentKey)
					}
				}
			}

		}
	}
}

//TODO: нормально описать ошибку
//changeTag меняет тэги в карте на основе путей до этих тэгов
func changeTag(jsonMap map[string]interface{}, tagName string, newValue interface{}) (map[string]interface{}, error) {
	workMap := func(shell map[string]interface{}) map[string]interface{} {
		return shell["dashboard"].(map[string]interface{})
	}(jsonMap)
	byteJSON, err := json.Marshal(jsonMap)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile("output1.json", byteJSON, 0777)
	stringJSON := string(byteJSON)
	if err != nil {
		return nil, err
	}
	getRoutes(workMap, "dashboard", tagName)
	for _, i := range routes {
		// fmt.Println(i)
		stringJSON, err = sjson.Set(stringJSON, i, newValue)
		if err != nil {
			panic(err)
		}
	}

	jsonMap, err = stringToMap(stringJSON)
	if err != nil {
		return nil, err
	}
	return jsonMap, nil
}

//TODO: нормально описать ошибку
//stringToMap превращает JSON строку в карту
func stringToMap(jsonString string) (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//TODO: нормально описать ошибку
//mapToFile ползволяет записать json карту в файл
func mapToFile(filename string, jsonMap map[string]interface{}) error {
	data, err := json.MarshalIndent(jsonMap, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		return err
	}
	return nil
}
