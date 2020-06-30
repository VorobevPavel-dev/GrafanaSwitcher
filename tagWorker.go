package main

import (
	"fmt"
	"strconv"
	"strings"
)

//recursiveReplaceTag позволяет изменить ВСЕ значения тэга name на Value в map[string]interface{}.
//Первое возвращаемое значение отвечает за количество изменений.
//ВОЗВРАЩАЕТ ТОЛЬКО ВНУТРЕННОСТЬ ТЭГА DASHBOARD
func recursiveReplaceTag(currentTag map[string]interface{}, name string, Value interface{}, changes int) (int, map[string]interface{}) {
	count := 0
	var tempCount int
	for key, value := range currentTag {
		currentType := fmt.Sprintf("%T", value)
		switch currentType {
		case "float64":
			if key == name {
				value = Value.(float64)
				currentTag[key] = value
				count++
			}
		case "string":
			if key == name {
				value = Value.(string)
				currentTag[key] = value
				count++
			}
		case "bool":
			if key == name {
				value = Value.(bool)
				currentTag[key] = value
				count++
			}
		case "map[string]interface {}":
			tempCount, currentTag[key] = recursiveReplaceTag(value.(map[string]interface{}), name, Value, 0)
			count += tempCount
		case "nil":
			if key == name {
				value = Value.(string)
				currentTag[key] = value
				count++
			}
		case "[]interface {}":
			temp := value.([]interface{})
			for i := range temp {
				val, ok := temp[i].(map[string]interface{})
				if ok {
					tempCount, temp[i] = recursiveReplaceTag(val, name, Value, 0)
					count += tempCount
				}
			}
		}
	}
	return count, currentTag
}

//recursivePrintTag позволяет вывести все тэги из map[string]interface{} и их значения в виде дерева
func recursivePrintTag(currentTag map[string]interface{}, currentDepth int) {
	for key, value := range currentTag {
		currentType := fmt.Sprintf("%T", value)
		fmt.Println(strings.Repeat("\t", currentDepth) + key)
		switch currentType {
		case "float64":
			fmt.Println(strings.Repeat("\t", currentDepth+1) + fmt.Sprintf("%f", value.(float64)))
		case "string":
			fmt.Println(strings.Repeat(string("\t"), currentDepth+1) + value.(string))
		case "bool":
			val := strconv.FormatBool(value.(bool))
			fmt.Println(strings.Repeat("\t", currentDepth+1) + val)
		case "map[string]interface {}":
			recursivePrintTag(value.(map[string]interface{}), currentDepth+1)
		case "[]interface {}":
			temp := value.([]interface{})
			for i := range temp {
				val, ok := temp[i].(map[string]interface{})
				if ok {
					recursivePrintTag(val, currentDepth+1)
				} else {
					for i := range temp {
						fmt.Println(strings.Repeat("\t", currentDepth+1), temp[i])
					}
				}
			}
		}
	}
}
