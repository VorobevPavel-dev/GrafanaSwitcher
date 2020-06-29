package main

import (
	"fmt"
	"strconv"
	"strings"
)

func recursiveReplaceTag(currentTag map[string]interface{}, name string, Value interface{}) map[string]interface{} {
	for key, value := range currentTag {
		currentType := fmt.Sprintf("%T", value)
		switch currentType {
		case "float64":
			if key == name {
				value = Value.(float64)
				currentTag[key] = value
			}
		case "string":
			if key == name {
				value = Value.(string)
				currentTag[key] = value
			}
		case "bool":
			if key == name {
				fmt.Println("changed")
				value = Value.(bool)
				currentTag[key] = value
			}
		case "map[string]interface {}":
			currentTag[key] = recursiveReplaceTag(value.(map[string]interface{}), name, Value)
		case "[]interface {}":
			temp := value.([]interface{})
			for i := range temp {
				val, ok := temp[i].(map[string]interface{})
				if ok {
					currentTag[key] = recursiveReplaceTag(val, name, Value)
				}
			}
		}
	}
	return currentTag
}

func recursivePrintTag(currentTag map[string]interface{}, currentDepth int) {
	for key, value := range currentTag {
		// fmt.Println("______________________")
		currentType := fmt.Sprintf("%T", value)
		// fmt.Println(currentType)
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
			// fmt.Println("Hello")
			recursivePrintTag(value.(map[string]interface{}), currentDepth+1)
		case "[]interface {}":
			temp := value.([]interface{})
			// fmt.Println(len(temp))
			// _, ok := temp[0].(map[string]interface{})
			// if ok {
			// 	for i := range temp {
			// 		data := temp[i].(map[string]interface{})
			// 		recursivePrintTag(data, currentDepth+1)
			// 	}
			// } else {
			// 	for i := range temp {
			// 		fmt.Println(temp[i].(string))
			// 	}
			// }
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
