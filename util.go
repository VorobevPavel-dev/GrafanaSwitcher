package main

import (
	"bytes"
)

func removeMetaTag(temp map[string]interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	for key, value := range temp {
		if key != "meta" {
			ret[key] = value
		}
	}
	return ret
}

func repairJSON(data []byte) []byte {
	data = bytes.Replace(data, []byte("\\u003c"), []byte("<"), -1)
	data = bytes.Replace(data, []byte("\\u003e"), []byte(">"), -1)
	data = bytes.Replace(data, []byte("\\u0026"), []byte("&"), -1)
	return data
}
