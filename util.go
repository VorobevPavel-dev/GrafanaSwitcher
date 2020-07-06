package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

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
	// for i := range routes {
	// 	fmt.Println(i)
	// }
}

//changeTag меняет тэги в карте на основе путей до этих тэгов
func changeTag(jsonMap map[string]interface{}, tagName string, newValue interface{}) (map[string]interface{}, error) {
	//Случай, если тэг без точек - просто выполняем функцию
	//В противном случае надо изменить карту и в конце поменять все иероглифы на точки
	if strings.Contains(tagName, ".") {
		toFix, _ := mapToString(jsonMap)
		toFix = strings.ReplaceAll(toFix, ".", "語")
		// fmt.Println(toFix)
		jsonMap, _ = stringToMap(toFix)
		tagName = strings.ReplaceAll(tagName, ".", "語")
	}
	workMap := func(shell map[string]interface{}) map[string]interface{} {
		return shell["dashboard"].(map[string]interface{})
	}(jsonMap)
	byteJSON, err := json.Marshal(jsonMap)
	if err != nil {
		return nil, errors.New("Cannot convert to json (util.changeTag())")
	}
	stringJSON := string(byteJSON)
	// fmt.Println(tagName)
	// fmt.Println(newValue)
	// fmt.Println(jsonMap)
	getRoutes(workMap, "dashboard", tagName)
	if len(routes) == 0 {
		return nil, errors.New("Tag " + tagName + " not found in this JSONModel")
	}
	for _, i := range routes {
		stringJSON, err = sjson.Set(stringJSON, i, newValue)
		if err != nil {
			return nil, errors.New("Cannot change tag (util.changeTag())")
		}
	}

	if strings.Contains(stringJSON, "語") {
		stringJSON = strings.ReplaceAll(stringJSON, "語", ".")
	}
	jsonMap, err = stringToMap(stringJSON)

	if err != nil {
		return nil, errors.New("Cannot convert to map (util.changeTag())")
	}
	return jsonMap, nil
}

//stringToMap превращает JSON строку в карту
func stringToMap(jsonString string) (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil, errors.New("Cannot unmarshall (util.stringToMap())")
	}
	return result, nil
}

//mapToString превращает карту в JSON
func mapToString(tag map[string]interface{}) (string, error) {
	result, err := json.MarshalIndent(tag, "", "\t")
	if err != nil {
		return "", errors.New("Cannot convert to string (mapToString())")
	}
	return string(result), nil
}

//mapToFile ползволяет записать json карту в файл
func mapToFile(filename string, jsonMap map[string]interface{}) error {
	data, err := json.MarshalIndent(jsonMap, "", "\t")
	if err != nil {
		return errors.New("Cannot caonvert to json (util.mapToFile())")
	}
	data = repairJSON(data)
	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		return errors.New("Cannot write to file (util.mapToFile())")
	}
	return nil
}
